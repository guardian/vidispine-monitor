package vshealthcheck

import "time"

type HealthcheckEntry struct {
	Healthy   bool      `json:"healthy"`
	Duration  int       `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
}

type HealthcheckResponse struct {
	Broker        HealthcheckEntry `json:"broker"`
	Database      HealthcheckEntry `json:"database"`
	Deadlocks     HealthcheckEntry `json:"deadlocks"`
	Elasticsearch HealthcheckEntry `json:"elasticsearch"`
	Ldap          HealthcheckEntry `json:"ldap"`
}
