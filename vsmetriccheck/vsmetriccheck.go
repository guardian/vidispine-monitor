package vsmetriccheck

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type VSMetricCheck struct {
	VidispineHost   string
	VidispineHttps  bool
	VidispineDbName string
	IntegrationKey  string
}

func (m VSMetricCheck) Name() string {
	return "Connection pool and error response rate"
}

/**
get the metrics response from the :9001 admin service
*/
func (m VSMetricCheck) loadMetrics() (*MetricsResponse, error) {
	httpClient := http.Client{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFunc()

	proto := "https:/"
	if !m.VidispineHttps {
		proto = "http:/"
	}

	urlParts := []string{
		proto,
		m.VidispineHost + ":9001",
		"metrics",
	}

	urlStr := strings.Join(urlParts, "/")
	_, urlParseErr := url.Parse(urlStr)
	if urlParseErr != nil {
		log.Printf("ERROR vsHealthcheck.loadHealthcheck URL %s is not valid: %s", urlStr, urlParseErr)
		return nil, urlParseErr
	}

	httpReq, reqErr := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	response, httpErr := httpClient.Do(httpReq)
	if httpErr != nil {
		return nil, httpErr
	}
	defer response.Body.Close()

	rawContent, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	var metrics MetricsResponse
	unmarshalErr := json.Unmarshal(rawContent, &metrics)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &metrics, nil
}

/**
returns a PD event if either active connections makes up for >90% of total pool or idle+active makes up for >80%
*/
func (m VSMetricCheck) CheckDatabasePool(metrics *MetricsResponse, verboseMode bool) *pagerduty.TriggerEvent {
	//we use MustFloat() to simplify coding, therefore we need to catch any panics that occur
	defer func() {
		if r := recover(); r != nil {
			log.Print("ERROR could not process database metrics: ", r)
		}
	}()

	sizeKey := fmt.Sprintf("io.dropwizard.db.ManagedPooledDataSource.%s.size", m.VidispineDbName)
	poolSizeTotal, havePoolSizeTotal := metrics.Gauges[sizeKey]
	idleKey := fmt.Sprintf("io.dropwizard.db.ManagedPooledDataSource.%s.idle", m.VidispineDbName)
	poolIdle, havePoolIdle := metrics.Gauges[idleKey]
	activeKey := fmt.Sprintf("io.dropwizard.db.ManagedPooledDataSource.%s.active", m.VidispineDbName)
	poolActive, havePoolActive := metrics.Gauges[activeKey]

	if !havePoolActive || !havePoolIdle || !havePoolSizeTotal {
		log.Print("WARNING metrics response was missing some of the database metrics, can't alert on database")
		return nil
	}

	if verboseMode {
		log.Printf("INFO (verbose) vsmetriccheck.CheckDatabasePool total pool size is %.1f, with %.1f currently active and %.1f idle",
			poolSizeTotal.MustFloat(), poolActive.MustFloat(), poolIdle.MustFloat())
	}

	if poolActive.MustFloat() > 0.9*poolSizeTotal.MustFloat() {
		nowTime := time.Now()
		log.Print("WARNING 90% or more of connection pool active, alerting")
		return pagerduty.NewTriggerEvent("vidispine-database",
			m.IntegrationKey,
			pagerduty.SeverityCritical,
			"vidispine-database-pool",
			"Active database connections account for over 90% of pool capacity, failure is imminent",
			&nowTime)
	}

	if (poolIdle.MustFloat() + poolActive.MustFloat()) > 0.8*poolSizeTotal.MustFloat() {
		nowTime := time.Now()
		log.Print("WARNING 80% or more of connection pool capacity is either idle or active, alerting")
		return pagerduty.NewTriggerEvent("vidispine-database",
			m.IntegrationKey,
			pagerduty.SeverityWarning,
			"vidispine-database-pool",
			"Spare database connection pool capacity (neither active nor idle) is less than 20%",
			&nowTime)
	}
	return nil
}

func (m VSMetricCheck) CheckHeapUsage(metrics *MetricsResponse, verboseMode bool) *pagerduty.TriggerEvent {
	defer func() {
		if r := recover(); r != nil {
			log.Print("ERROR could not process heap usage: ", r)
		}
	}()

	heapUsage, haveHeapUsage := metrics.Gauges["jvm.memory.heap.usage"]
	if !haveHeapUsage {
		log.Print("WARNING jvm.memory.heap.usage is not present in metric gauges, can't check heap usage")
		return nil
	}

	if verboseMode {
		log.Printf("INFO (verbose) vsmetriccheck.CheckHeapUsage JVM heap usage is %.1f%%", heapUsage.MustFloat()*100)
	}

	if heapUsage.MustFloat() > 0.9 {
		nowTime := time.Now()
		log.Print("WARNING heap usage is at 90%, alerting")
		return pagerduty.NewTriggerEvent(
			"vidispine-heap",
			m.IntegrationKey,
			pagerduty.SeverityCritical,
			"vidispine-heap",
			"Vidispine heap RAM usage is at 90%, failure is likely. Pod needs restarting and RAM allocation re-assessing",
			&nowTime,
		)
	}

	if heapUsage.MustFloat() > 0.7 {
		nowTime := time.Now()
		log.Print("WARNING heap usage is at 70%, alerting")
		return pagerduty.NewTriggerEvent(
			"vidispine-heap",
			m.IntegrationKey,
			pagerduty.SeverityWarning,
			"vidispine-heap",
			"Vidispine heap RAM usage is at 70%, monitor and update RAM allocation before failures are likely",
			&nowTime,
		)
	}
	return nil
}

func (m VSMetricCheck) Run(verboseMode bool) ([]*pagerduty.TriggerEvent, error) {
	metrics, err := m.loadMetrics()
	if err != nil {
		log.Print("ERROR could not load metrics from Vidispine admin service: ", err)
		return nil, err
	}

	alerts := make([]*pagerduty.TriggerEvent, 0)

	poolAlert := m.CheckDatabasePool(metrics, verboseMode)
	if poolAlert != nil {
		alerts = append(alerts, poolAlert)
	}

	heapAlert := m.CheckHeapUsage(metrics, verboseMode)
	if heapAlert != nil {
		alerts = append(alerts, heapAlert)
	}
	return alerts, nil
}
