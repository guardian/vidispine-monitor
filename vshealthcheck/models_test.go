package vshealthcheck

import (
	"encoding/json"
	"testing"
	"time"
)

/**
Ensure that a valid json from Vidispine can be parsed correctly
*/
func TestParseHealthcheckResponse(t *testing.T) {
	testContent := `{
  "broker": {
    "healthy": true,
    "duration": 0,
    "timestamp": "2021-03-11T11:49:20.441Z"
  },
  "database": {
    "healthy": true,
    "duration": 0,
    "timestamp": "2021-03-11T11:49:20.435Z"
  },
  "deadlocks": {
    "healthy": true,
    "duration": 0,
    "timestamp": "2021-03-11T11:49:20.442Z"
  },
  "elasticsearch": {
    "healthy": true,
    "duration": 4,
    "timestamp": "2021-03-11T11:49:20.440Z"
  },
  "ldap": {
    "healthy": true,
    "duration": 1,
    "timestamp": "2021-03-11T11:49:20.441Z"
  }
}
`
	var parsedContent HealthcheckResponse
	unmarshalErr := json.Unmarshal([]byte(testContent), &parsedContent)
	if unmarshalErr != nil {
		t.Error("Got unexpected unmarshalling error: ", unmarshalErr)
		t.FailNow()
	}

	if parsedContent.Broker.Healthy != true {
		t.Error("Broker healthy should be true")
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2021-03-11T11:49:20.441Z")
	if parsedContent.Broker.Timestamp != expectedTime {
		t.Errorf("Broker check was wrong, got %s expected %s", parsedContent.Broker.Timestamp, expectedTime)
	}
	if parsedContent.Broker.Duration != 0 {
		t.Error("Broker duration should be 0")
	}

	if parsedContent.Database.Healthy != true {
		t.Error("Database healthy should be true")
	}
	expectedTimeDb, _ := time.Parse(time.RFC3339, "2021-03-11T11:49:20.435Z")
	if parsedContent.Database.Timestamp != expectedTimeDb {
		t.Errorf("Database check was wrong, got %s expected %s", parsedContent.Database.Timestamp, expectedTimeDb)
	}
	if parsedContent.Database.Duration != 0 {
		t.Error("Database duration should be 0")
	}

	if parsedContent.Deadlocks.Healthy != true {
		t.Error("Deadlocks healthy should be true")
	}
	expectedTimeDl, _ := time.Parse(time.RFC3339, "2021-03-11T11:49:20.442Z")
	if parsedContent.Deadlocks.Timestamp != expectedTimeDl {
		t.Errorf("Deadlocks check was wrong, got %s expected %s", parsedContent.Deadlocks.Timestamp, expectedTimeDl)
	}
	if parsedContent.Deadlocks.Duration != 0 {
		t.Error("Deadlocks duration should be 0")
	}

	if parsedContent.Elasticsearch.Healthy != true {
		t.Error("Elasticsearch healthy should be true")
	}
	expectedTimeEs, _ := time.Parse(time.RFC3339, "2021-03-11T11:49:20.440Z")
	if parsedContent.Elasticsearch.Timestamp != expectedTimeEs {
		t.Errorf("Elasticsearch check was wrong, got %s expected %s", parsedContent.Elasticsearch.Timestamp, expectedTimeEs)
	}
	if parsedContent.Elasticsearch.Duration != 4 {
		t.Error("Elasticsearch duration should be 4")
	}

	if parsedContent.Ldap.Healthy != true {
		t.Error("LDAP healthy should be true")
	}
	expectedTimeLd, _ := time.Parse(time.RFC3339, "2021-03-11T11:49:20.441Z")
	if parsedContent.Ldap.Timestamp != expectedTimeLd {
		t.Errorf("LDAP check was wrong, got %s expected %s", parsedContent.Ldap.Timestamp, expectedTimeLd)
	}
	if parsedContent.Ldap.Duration != 1 {
		t.Error("LDAP duration should be 1")
	}
}
