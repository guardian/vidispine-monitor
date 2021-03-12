package vsmetriccheck

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseMetrics(t *testing.T) {
	file, openErr := os.Open("sample_metrics_formatted.json")
	if openErr != nil {
		t.Error(openErr)
		t.FailNow()
	}
	defer file.Close()

	rawBytes, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		t.Error(readErr)
		t.FailNow()
	}

	var parsed MetricsResponse
	unmarshalErr := json.Unmarshal(rawBytes, &parsed)
	if unmarshalErr != nil {
		t.Error(unmarshalErr)
		t.FailNow()
	}

	servletActiveRequests, haveServletActiveRequests := parsed.Counters["io.dropwizard.jetty.MutableServletContextHandler.active-requests"]
	if !haveServletActiveRequests {
		t.Error("Expected Counters to have io.dropwizard.jetty.MutableServletContextHandler.active-requests")
	} else {
		if servletActiveRequests.Count != 0 {
			t.Error("Expected io.dropwizard.jetty.MutableServletContextHandler.active-requests to be a counter of value 0")
		}
	}

	clusterSize, haveClusterSize := parsed.Gauges["cluster.size"]
	if !haveClusterSize {
		t.Error("Expected gauges to have cluster.size")
	} else {
		floatValue, floatErr := clusterSize.FloatValue()
		if floatErr != nil {
			t.Error(floatErr)
		} else {
			if floatValue != 1.0 {
				t.Error("Expected cluster.size gauge to be 1.0, not ", floatValue)
			}
			mustValue := clusterSize.MustFloat()
			if mustValue != 1.0 {
				t.Error("Expected cluster.size gauge to be 1.0, not ", floatValue)
			}
		}
	}
}
