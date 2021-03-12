package vsstoragecheck

import (
	"encoding/json"
	"testing"
)

/**
We must be able to parse a real response from the server
*/
func TestParseVSStoragesResponse(t *testing.T) {
	rawResponse := `{
    "storage": [
        {
            "id": "VX-4",
            "state": "NONE",
            "priority": "MEDIUM",
            "type": "LOCAL",
            "capacity": 17361125376,
            "freeCapacity": 27721728,
            "timestamp": "2021-03-11T18:40:08.959+0000",
            "method": [
                {
                    "id": "VX-2",
                    "uri": "file:///srv/media/",
                    "read": true,
                    "write": true,
                    "browse": true,
                    "lastSuccess": "2021-03-11T18:40:08.950+0000",
                    "type": "NONE"
                },
                {
                    "id": "VX-8",
                    "uri": "https://vidispine.local/",
                    "bandwidth": 0,
                    "read": false,
                    "write": false,
                    "browse": false,
                    "type": "AUTO"
                }
            ],
            "metadata": {},
            "lowWatermark": 17361125376,
            "highWatermark": 17361125376,
            "autoDetect": true,
            "showImportables": true
        }
    ]
}`
	var result VSStoragesResponse
	unmarshalErr := json.Unmarshal([]byte(rawResponse), &result)
	if unmarshalErr != nil {
		t.Error("Got unexpected error unmarshalling server content: ", unmarshalErr)
		t.FailNow()
	}

	if len(result.Storage) != 1 {
		t.Errorf("Got unexpected storage count %d from sample data", len(result.Storage))
		t.FailNow()
	}
	s := result.Storage[0]
	if s.Id != "VX-4" {
		t.Error("Got unexpected storage ID")
	}
	if s.State != "NONE" {
		t.Error("Got unexpected storage state")
	}
}
