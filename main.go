package main

import (
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/common"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/vshealthcheck"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/vsmetriccheck"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/vsstoragecheck"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	checkEveryStr := os.Getenv("CHECK_EVERY")                        //interval to check, parsed as a duration
	pdService := os.Getenv("PD_INTEGRATION_KEY")                     //pagerduty service ID to alert
	pdApiKey := os.Getenv("PD_API_KEY")                              //API key to communicate with PD
	vidispineHost := os.Getenv("VIDISPINE_HOST")                     //hostname to query
	vidispineMonitorHttpsStr := os.Getenv("VIDISPINE_MONITOR_HTTPS") //set to TRUE if the 9001 monitoring port is https protected
	vidispineApiHttpsStr := os.Getenv("VIDISPINE_API_HTTPS")         //set to TRUE if the 8080 API port is https protected
	vidispineApiUser := os.Getenv("VIDISPINE_API_USER")              //API user for checking storages
	vidispineApiPasswd := os.Getenv("VIDISPINE_API_PASSWD")          //API password for checking storages
	verboseStr := os.Getenv("VERBOSE")                               //whether to output verbose logging
	sendTestMessageStr := os.Getenv("TEST_MESSAGE")                  //if set, then send a test message to PD

	if vidispineHost == "" {
		log.Fatal("You must specify VIDISPINE_HOST in the environment. Note that this is the hostname not the url.")
	}

	if checkEveryStr == "" {
		log.Fatal("You must specify CHECK_EVERY in the environment, e.g. CHECK_EVERY=5minutes")
	}
	checkEvery, durParseErr := time.ParseDuration(checkEveryStr)
	if durParseErr != nil {
		log.Fatalf("CHECK_EVERY value %s is not a valid duration: %s", checkEveryStr, durParseErr)
	}

	if pdService == "" {
		log.Print("WARNING PD_INTEGRATION_KEY and/or PD_API_KEY is not set, no alerts can be raised to pagerduty")
	}

	verboseMode := false
	if verboseStr != "" {
		var boolParseErr error
		verboseMode, boolParseErr = strconv.ParseBool(verboseStr)
		if boolParseErr != nil {
			log.Fatalf("The value %s for VERBOSE is not valid, expected 'true' or 'false'", verboseStr)
		}
	}

	vidispineMonitorHttps := false
	if vidispineMonitorHttpsStr != "" {
		var boolParseErr error
		vidispineMonitorHttps, boolParseErr = strconv.ParseBool(vidispineMonitorHttpsStr)
		if boolParseErr != nil {
			log.Fatalf("The value %s for VIDISPINE_MONITOR_HTTPS was not valid, expected 'true' or 'false'", vidispineMonitorHttpsStr)
		}
	}

	vidispineApiHttps := false
	if vidispineApiHttpsStr != "" {
		var boolParseErr error
		vidispineApiHttps, boolParseErr = strconv.ParseBool(vidispineApiHttpsStr)
		if boolParseErr != nil {
			log.Fatalf("The value %s for VIDISPINE_API_HTTPS was not valid, expected 'true' or 'false'", boolParseErr)
		}
	}

	healthChecks := []common.MonitorComponent{
		vshealthcheck.VSHealthCheckMonitor{
			VidispineHost:  vidispineHost,
			VidispineHttps: vidispineMonitorHttps,
			PDServiceId:    pdService,
		},
		vsmetriccheck.VSMetricCheck{
			VidispineHost:   vidispineHost,
			VidispineHttps:  vidispineMonitorHttps,
			VidispineDbName: "vidispinedb",
			IntegrationKey:  pdService,
		},
	}

	if vidispineApiUser != "" && vidispineApiPasswd != "" {
		healthChecks = append(healthChecks, vsstoragecheck.VSStorageCheck{
			VidispineHost:    vidispineHost,
			VidispineUser:    vidispineApiUser,
			VidispinePasswd:  vidispineApiPasswd,
			PDIntegrationKey: pdService,
			VidispineHttps:   vidispineApiHttps,
		})
	} else {
		log.Print("WARNING No vidispine api user and/or password was specified, can't do storage detail checks")
	}

	if sendTestMessageStr != "" {
		nowtime := time.Now()
		testMessage := pagerduty.NewTriggerEvent("vidispine-monitor",
			pdService,
			pagerduty.SeverityInfo,
			"test-message",
			"Test message from vidispine-monitor",
			&nowtime)
		sendErr := pagerduty.SendEvent(testMessage, pdApiKey, 60*time.Second)
		if sendErr == nil {
			log.Print("INFO test message sent succesfully")
		} else {
			log.Fatal("ERROR could not send test message: ", sendErr)
		}
	}

	for {
		didFail := false
		for _, check := range healthChecks {
			alerts, runErr := check.Run(verboseMode)
			if runErr != nil {
				didFail = true
				log.Printf("ERROR running '%s' failed: %s", check.Name(), runErr)
			}
			if alerts != nil && len(alerts) > 0 {
				log.Printf("WARNING %s returned %d alerts: ", check.Name(), len(alerts))
				for _, alert := range alerts {
					log.Printf("WARNING [%s] %s", check.Name(), alert.String())
					if pdService == "" {
						log.Print("WARNING can't send to pagerduty as PD_INTEGRATION_KEY not set")
					} else {
						sendErr := pagerduty.SendEvent(alert, pdApiKey, 60*time.Second)
						if sendErr != nil {
							log.Printf("ERROR Could not sent alert %s: %s", alert, sendErr)
						}
					}
				}
			}
		}
		if didFail {
			log.Print("ERROR Some internal errors occurred while processing the warnings, terminating")
			os.Exit(1)
		}
		time.Sleep(checkEvery)
	}
}
