package vshealthcheck

import (
	"encoding/json"
	"fmt"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type VSHealthCheckMonitor struct {
	VidispineHost  string
	VidispineHttps bool
	PDServiceId    string
}

/**
gets helthcheck data from the VS endpoint at the given host
*/
func (m VSHealthCheckMonitor) loadHealthcheck(vsHost string, vsHttps bool) (*HealthcheckResponse, error) {
	proto := "https:/"
	if !vsHttps {
		proto = "http:/"
	}

	urlParts := []string{
		proto,
		vsHost + ":9001",
		"healthcheck",
	}

	urlStr := strings.Join(urlParts, "/")
	_, urlParseErr := url.Parse(urlStr)
	if urlParseErr != nil {
		log.Printf("ERROR vsHealthcheck.loadHealthcheck URL %s is not valid: %s", urlStr, urlParseErr)
		return nil, urlParseErr
	}

	httpResp, httpErr := http.Get(urlStr)
	if httpErr != nil {
		return nil, httpErr
	}
	defer httpResp.Body.Close()

	content, readErr := ioutil.ReadAll(httpResp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var healthcheck HealthcheckResponse
	unmarshalErr := json.Unmarshal(content, &healthcheck)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &healthcheck, nil
}

/**
check an individual component
*/
func (m VSHealthCheckMonitor) validateHealthcheckEntry(name string, entry *HealthcheckEntry, verboseMode bool) *pagerduty.CreateIncidentRequest {
	if verboseMode {
		log.Printf("INFO (verbose) validateHealthcheckEntry for %s got %v", name, entry)
	}
	if !entry.Healthy {
		bodyText := fmt.Sprintf("The %s check failed at %s", name, entry.Timestamp.String())
		if verboseMode {
			log.Printf("INFO (verbose) %s", bodyText)
		}
		return pagerduty.NewIncidentRequest(
			fmt.Sprintf("Vidispine %s check failed", name),
			m.PDServiceId,
			pagerduty.UrgencyHigh,
			fmt.Sprintf("vidispine-%s", strings.ToLower(name)),
			bodyText,
		)
	} else {
		log.Printf("INFO (verbose) validateHealthceckEntry %s passed", name)
		return nil //we are healthy, nothing to see here
	}
}

func (m VSHealthCheckMonitor) Name() string {
	return "Vidispine basic health checks"
}

/**
runs the check on Vidispine health
*/
func (m VSHealthCheckMonitor) Run(verboseMode bool) ([]*pagerduty.CreateIncidentRequest, error) {
	if verboseMode {
		log.Printf("INFO (verbose) Checking %s on %s", m.Name(), m.VidispineHost)
	}

	healthCheckResponse, err := m.loadHealthcheck(m.VidispineHost, m.VidispineHttps)
	if err != nil {
		log.Print("ERROR vshealthcheck could not run: ", err)
		bodyText := err.Error()
		return []*pagerduty.CreateIncidentRequest{
			pagerduty.NewIncidentRequest("Vidispine healthcheck could not run", m.PDServiceId, pagerduty.UrgencyHigh, "vshealthcheck", bodyText),
		}, err
	}

	if verboseMode {
		log.Printf("INFO (verbose) Got check results, evaluating...")
	}
	checkList := []*HealthcheckEntry{
		&healthCheckResponse.Broker,
		&healthCheckResponse.Database,
		&healthCheckResponse.Deadlocks,
		&healthCheckResponse.Elasticsearch,
		&healthCheckResponse.Ldap,
	}
	checkNames := []string{
		"Broker",
		"Database",
		"Deadlocks",
		"Elasticsearch",
		"LDAP",
	}

	errors := make([]*pagerduty.CreateIncidentRequest, 0)
	for i, check := range checkList {
		problem := m.validateHealthcheckEntry(checkNames[i], check, verboseMode)
		if problem != nil {
			errors = append(errors, problem)
		}
	}

	if verboseMode {
		log.Printf("INFO (verbose) Check returned %d problems", len(errors))
	}

	return errors, nil
}
