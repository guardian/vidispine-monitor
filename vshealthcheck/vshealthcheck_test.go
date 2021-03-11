package vshealthcheck

import (
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"
	"testing"
	"time"
)

func TestVSHealthCheckMonitor_validateHealthcheckEntry_unhealthy(t *testing.T) {
	m := VSHealthCheckMonitor{
		VidispineHost:  "somehost",
		VidispineHttps: false,
		PDServiceId:    "someservice",
	}

	faketime, _ := time.Parse(time.RFC3339, "2010-01-02T03:04:05.678Z")
	toValidate := &HealthcheckEntry{
		Healthy:   false,
		Duration:  123,
		Timestamp: faketime,
	}

	result := m.validateHealthcheckEntry("test", toValidate, false)
	if result == nil {
		t.Error("validateHealthcheckEntry returned no problem on an unhealthy entry")
	} else {
		if result.IntegrationKey != "someservice" {
			t.Errorf("alert had incorrect IntegrationKey '%s'", result.IntegrationKey)
		}
		if result.EventAction != pagerduty.EventActionTrigger {
			t.Errorf("alert had incorrect EventAction '%s'", result.EventAction)
		}
		if result.Payload.Summary != "The test check failed at 2010-01-02 03:04:05.678 +0000 UTC" {
			t.Errorf("alert had incorrect summary '%s'", result.Payload.Summary)
		}
		if result.Payload.Severity != pagerduty.SeverityError {
			t.Errorf("alert had incorrect severity '%s'", result.Payload.Severity)
		}
		if result.DeDupKey != "vidispine-test" {
			t.Errorf("alert had incorrect incident key '%s'", result.DeDupKey)
		}
	}
}

func TestVSHealthCheckMonitor_validateHealthcheckEntry_healthy(t *testing.T) {
	m := VSHealthCheckMonitor{
		VidispineHost:  "somehost",
		VidispineHttps: false,
		PDServiceId:    "someservice",
	}

	faketime, _ := time.Parse(time.RFC3339, "2010-01-02T03:04:05.678Z")
	toValidate := &HealthcheckEntry{
		Healthy:   true,
		Duration:  123,
		Timestamp: faketime,
	}

	result := m.validateHealthcheckEntry("test", toValidate, false)
	if result != nil {
		t.Error("validateHealthcheckEntry returned an unexpected problem: ", result.String())
	}

}
