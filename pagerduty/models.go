package pagerduty

import (
	"fmt"
	"time"
)

//https://developer.pagerduty.com/docs/events-api-v2/trigger-ev

type EventAction string

const (
	EventActionTrigger     EventAction = "trigger"
	EventActionAcknowledge EventAction = "acknowledge"
	EventActionResolve     EventAction = "resolve"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityError    Severity = "error"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

type TriggerEventPayload struct {
	Summary   string   `json:"summary"`
	Timestamp string   `json:"timestamp"`
	Source    string   `json:"source"`
	Severity  Severity `json:"severity"`
	Component string   `json:"component"`
	Group     string   `json:"group"`
	Class     string   `json:"class"`
}

type TriggerEvent struct {
	IntegrationKey string              `json:"routing_key"`
	EventAction    EventAction         `json:"event_action"` //MUST be "trigger"
	DeDupKey       string              `json:"dedup_key"`
	Payload        TriggerEventPayload `json:"payload"`
}

func NewTriggerEvent(component string, integrationKey string, severity Severity, incidentKey string, incidentBody string, timestamp *time.Time) *TriggerEvent {
	timestampStr := timestamp.Format(time.RFC3339)

	return &TriggerEvent{
		IntegrationKey: integrationKey,
		EventAction:    EventActionTrigger,
		DeDupKey:       incidentKey,
		Payload: TriggerEventPayload{
			Summary:   incidentBody,
			Timestamp: timestampStr,
			Source:    "vidispine",
			Severity:  severity,
			Component: component,
		},
	}
}

func (e *TriggerEvent) String() string {
	return e.Payload.Summary
}

//https://developer.pagerduty.com/api-reference/reference/REST/openapiv3.json/paths/~1incidents/post

type ObjectRefRequest struct {
	Id   string `json:"id"`   //REQUIRED, identity of the object
	Type string `json:"type"` //REQUIRED, type of the object
}

func PagerDutyService(serviceId string) *ObjectRefRequest {
	return &ObjectRefRequest{
		Id:   serviceId,
		Type: "service_reference",
	}
}

func PagerDutyPriority(prioId string) *ObjectRefRequest {
	return &ObjectRefRequest{
		Id:   prioId,
		Type: "priority_reference",
	}
}

func PagerDutyEscalationPolicy(policyId string) *ObjectRefRequest {
	return &ObjectRefRequest{
		Id:   policyId,
		Type: "escalation_policy_reference",
	}
}

type IncidentBody struct {
	Type    string `json:"type"` //REQUIRED, must be 'incident_body'
	Details string `json:"details"`
}

type UrgencyType string

const (
	UrgencyHigh UrgencyType = "high"
	UrgencyLow  UrgencyType = "low"
)

type CreateIncidentRequest struct {
	Type        string            `json:"type"`         //REQUIRED, must be 'incident'
	Title       string            `json:"title"`        //REQUIRED, succint title
	Service     *ObjectRefRequest `json:"service"`      //REQUIRED, service to target
	Priority    *ObjectRefRequest `json:"priority"`     //OPTIONAL, priority
	Urgency     UrgencyType       `json:"urgency"`      //OPTIONAL, low/high
	IncidentKey string            `json:"incident_key"` //OPTIONAL, for de-duplication
	Body        *IncidentBody     `json:"body"`         //REQUIRED, content of alert
}

func NewIncidentRequest(title string, serviceId string, urgency UrgencyType, incidentKey string, incidentBody string) *CreateIncidentRequest {
	return &CreateIncidentRequest{
		Type:        "incident",
		Title:       title,
		Service:     PagerDutyService(serviceId),
		Priority:    nil,
		Urgency:     urgency,
		IncidentKey: incidentKey,
		Body: &IncidentBody{
			Type:    "incident_body",
			Details: incidentBody,
		},
	}
}

func (rq *CreateIncidentRequest) String() string {
	return fmt.Sprintf("%s incident at %s priority", rq.Title, rq.Urgency)
}
