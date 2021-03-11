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
		if result.Service.Id != "someservice" {
			t.Errorf("alert had incorrect service id '%s'", result.Service.Id)
		}
		if result.Type != "incident" {
			t.Errorf("alert had incorrect type '%s'", result.Type)
		}
		if result.Title != "Vidispine test check failed" {
			t.Errorf("alert had incorrect title '%s'", result.Title)
		}
		if result.Urgency != pagerduty.UrgencyHigh {
			t.Errorf("alert had incorrect urgency '%s'", result.Urgency)
		}
		if result.IncidentKey != "vidispine-test" {
			t.Errorf("alert had incorrect incident key '%s'", result.IncidentKey)
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
