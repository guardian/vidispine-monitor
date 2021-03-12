package vsstoragecheck

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/common"
	"gitlab.com/codmill/customer-projects/guardian/vidispine-monitor/pagerduty"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type VSStorageCheck struct {
	VidispineHost    string
	VidispineUser    string
	VidispinePasswd  string
	PDIntegrationKey string
	VidispineHttps   bool
}

func (c VSStorageCheck) Name() string {
	return "Vidispine storages check"
}

func (c VSStorageCheck) loadStorageData() (*VSStoragesResponse, error) {
	httpClient := http.Client{}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFunc()

	proto := "https"
	if !c.VidispineHttps {
		proto = "http"
	}
	storageUrl := fmt.Sprintf("%s://%s:8080/API/storage", proto, c.VidispineHost)
	httpReq, reqErr := http.NewRequestWithContext(ctx, "GET", storageUrl, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	authStr := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.VidispineUser, c.VidispinePasswd)))
	httpReq.Header.Add("Authorization", fmt.Sprintf("Basic %s", authStr))
	httpReq.Header.Add("Accept", "application/json")
	response, httpErr := httpClient.Do(httpReq)
	if httpErr != nil {
		return nil, httpErr
	}

	defer response.Body.Close()

	contentBytes, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return nil, readErr
	}

	var storageResponse VSStoragesResponse
	unmarshalErr := json.Unmarshal(contentBytes, &storageResponse)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return &storageResponse, nil
}

func (c VSStorageCheck) CheckStorage(s *VSStorage, verboseMode bool) []*pagerduty.TriggerEvent {
	foundErrors := make([]*pagerduty.TriggerEvent, 0)

	if verboseMode {
		log.Printf("INFO (verbose) vsstoragecheck.CheckStorage %s %s state is %s", s.Type, s.Id, s.State)
	}

	if s.State != StorageStateNone && s.State != StorageStateReady {
		nowTime := time.Now()
		bodyText := fmt.Sprintf("%s storage %s entered %s state", s.Type, s.Id, s.State)

		stateErr := pagerduty.NewTriggerEvent(fmt.Sprintf("Storage %s", s.Id),
			c.PDIntegrationKey,
			pagerduty.SeverityError,
			fmt.Sprintf("vidispine-storagestate-%s", s.Id),
			bodyText,
			&nowTime)
		foundErrors = append(foundErrors, stateErr)
	}

	usedCap := s.Capacity - s.FreeCapacity
	if usedCap > s.HighWatermark {
		if verboseMode {
			log.Printf("INFO (verbose) %s %s watermark is at %d storage is over at %d", s.Type, s.Id, s.HighWatermark, usedCap)
		}
		nowTime := time.Now()
		bodyText := fmt.Sprintf("%s storage %s is at %s used, over the high watermark by %s", s.Type, s.Id, common.FormatBytes(usedCap), common.FormatBytes(usedCap-s.HighWatermark))
		watermarkErr := pagerduty.NewTriggerEvent(fmt.Sprintf("Storage %s", s.Id),
			c.PDIntegrationKey,
			pagerduty.SeverityError,
			fmt.Sprintf("vidispine-storagewatermark-%s", s.Id),
			bodyText,
			&nowTime)
		foundErrors = append(foundErrors, watermarkErr)
	} else {
		if verboseMode {
			log.Printf("INFO (verbose) %s %s watermark is at %s but storage ok at %s", s.Type, s.Id, common.FormatBytes(s.HighWatermark), common.FormatBytes(usedCap))
		}
	}

	dangerlevel := float64(s.Capacity) * 0.05

	if float64(s.FreeCapacity) < dangerlevel {
		nowTime := time.Now()
		bodyText := fmt.Sprintf("%s storage %s is over 95%% full, at %s", s.Type, s.Id,
			common.FormatBytes(usedCap))

		watermarkErr := pagerduty.NewTriggerEvent(fmt.Sprintf("Storage %s", s.Id),
			c.PDIntegrationKey,
			pagerduty.SeverityError,
			fmt.Sprintf("vidispine-storagefull-%s", s.Id),
			bodyText,
			&nowTime)
		foundErrors = append(foundErrors, watermarkErr)
	}
	return foundErrors
}

func (c VSStorageCheck) Run(verboseMode bool) ([]*pagerduty.TriggerEvent, error) {
	if verboseMode {
		log.Printf("INFO VSStorageCheck.Run Retrieving storage details")
	}

	storageInfo, err := c.loadStorageData()
	if err != nil {
		return nil, err
	}

	problems := make([]*pagerduty.TriggerEvent, 0)
	for _, storage := range storageInfo.Storage {
		problems = append(problems, c.CheckStorage(&storage, verboseMode)...)
	}

	if verboseMode {
		log.Printf("INFO VSStorageCheck.Run Out of %d storages, found %d problems (with 3 checks per storage)", len(storageInfo.Storage), len(problems))
	}
	return problems, nil
}
