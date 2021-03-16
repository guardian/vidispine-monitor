package vsmetriccheck

import (
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"
	"testing"
)

func TestVSMetricCheck_CheckDatabasePool_normal(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.size":   {Value: 100.0},
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.idle":   {Value: 20.0},
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.active": {Value: 8.0},
		},
		Counters: nil,
		Meters:   nil,
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckDatabasePool(fakeMetrics, false)
	if result != nil {
		t.Error("CheckDatabasePool returned an alert when everything is in order")
	}
}

func TestVSMetricCheck_CheckDatabasePool_70pc(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.size":   {Value: 100.0},
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.idle":   {Value: 30.0},
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.active": {Value: 55.0},
		},
		Counters: nil,
		Meters:   nil,
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckDatabasePool(fakeMetrics, false)
	if result == nil {
		t.Error("CheckDatabasePool returned no alert at 70% utilization")
	} else {
		if result.Payload.Severity != pagerduty.SeverityWarning {
			t.Errorf("CheckDatabasePool returned severity %s instead of warning for >70%%", result.Payload.Severity)
		}
	}
}

func TestVSMetricCheck_CheckDatabasePool_90pc(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.size":   {Value: 100.0},
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.idle":   {Value: 1.0},
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.active": {Value: 95.0},
		},
		Counters: nil,
		Meters:   nil,
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckDatabasePool(fakeMetrics, false)
	if result == nil {
		t.Error("CheckDatabasePool returned no alert at 90% utilization")
	} else {
		if result.Payload.Severity != pagerduty.SeverityCritical {
			t.Errorf("CheckDatabasePool returned severity %s instead of critical for >90%%", result.Payload.Severity)
		}
	}
}

func TestVSMetricCheck_CheckDatabasePool_notfound(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.db.ManagedPooledDataSource.vsdb.size": {Value: 100.0},
		},
		Counters: nil,
		Meters:   nil,
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckDatabasePool(fakeMetrics, false)

	if result != nil {
		t.Errorf("CheckDatabasePool returned an alert %v when there was no data", result)
	}
}

func TestVSMetricCheck_CheckHeapUsage_normal(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"jvm.memory.heap.usage": {Value: 0.31},
		},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckHeapUsage(fakeMetrics, false)
	if result != nil {
		t.Errorf("CheckHeapUsage returned an alert %v when value was in-range", result)
	}
}

/**
should return a critical error if heap usage >90%
*/
func TestVSMetricCheck_CheckHeapUsage_90pc(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"jvm.memory.heap.usage": {Value: 0.95},
		},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckHeapUsage(fakeMetrics, false)
	if result == nil {
		t.Error("CheckHeapUsage returned no alert when heap was at 95%")
	} else {
		if result.Payload.Severity != pagerduty.SeverityCritical {
			t.Errorf("CheckHeapUsage returned a %s error for 95%% heap when it should have been 'critical'.", result.Payload.Severity)
		}
	}
}

/**
should return a warning if heap usage >70%
*/
func TestVSMetricCheck_CheckHeapUsage_70pc(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"jvm.memory.heap.usage": {Value: 0.8},
		},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckHeapUsage(fakeMetrics, false)
	if result == nil {
		t.Error("CheckHeapUsage returned no alert when heap was at 80%")
	} else {
		if result.Payload.Severity != pagerduty.SeverityWarning {
			t.Errorf("CheckHeapUsage returned a %s for 80%% heap when it should have been 'warning'.", result.Payload.Severity)
		}
	}
}

/**
should return no error if no data
*/
func TestVSMetricCheck_CheckHeapUsage_nodata(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges:  map[string]MetricGauge{},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckHeapUsage(fakeMetrics, false)
	if result != nil {
		t.Errorf("CheckHeapUsage returned an alert %v when there was no data", result)
	}
}

func TestVSMetricCheck_CheckExcessive500s_normal(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-15m": {Value: 0.05},
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-5m":  {Value: 0.01},
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-1m":  {Value: 0.0},
		},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckExcessive500s(fakeMetrics, false)
	if result != nil {
		t.Errorf("CheckExcessive500s returned alert %v meters were in-range", result)
	}
}

func TestVSMetricCheck_CheckExcessive500s_1min(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-15m": {Value: 0.45},
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-5m":  {Value: 0.68},
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-1m":  {Value: 0.97},
		},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckExcessive500s(fakeMetrics, false)
	if result == nil {
		t.Error("CheckExcessive500s returned no alert when 97% of responses in 1min were 5xx", result)
	} else {
		if result.Payload.Severity != pagerduty.SeverityError {
			t.Errorf("CheckExcessive500s returned a '%s' severity when it should have been 'error'", result.Payload.Severity)
		}
	}
}

func TestVSMetricCheck_CheckExcessive500s_5min(t *testing.T) {
	fakeMetrics := &MetricsResponse{
		Version: "4.0.0",
		Gauges: map[string]MetricGauge{
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-15m": {Value: 0.45},
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-5m":  {Value: 0.68},
			"io.dropwizard.jetty.MutableServletContextHandler.percent-5xx-1m":  {Value: 0.39},
		},
	}

	c := VSMetricCheck{VidispineDbName: "vsdb"}
	result := c.CheckExcessive500s(fakeMetrics, false)
	if result == nil {
		t.Error("CheckExcessive500s returned no alert when 68% of responses in 5min were 5xx", result)
	} else {
		if result.Payload.Severity != pagerduty.SeverityWarning {
			t.Errorf("CheckExcessive500s returned a '%s' severity when it should have been 'warning'", result.Payload.Severity)
		}
	}
}
