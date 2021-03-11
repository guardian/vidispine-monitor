package common

import "gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"

type MonitorComponent interface {
	Run(verboseMode bool) ([]*pagerduty.CreateIncidentRequest, error) //perform the monitor checks. Return a list of CreateIncidentRequest, for each problem identified.
	Name() string                                                     //return a descriptive name for this check
}
