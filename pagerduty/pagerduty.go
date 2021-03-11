package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

/**
marshals a CreateIncidentRequest into a JSON request body and returns a ByteReader to it
*/
func generateEventBody(req *TriggerEvent) (io.Reader, error) {
	bodyContent, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return bytes.NewReader(bodyContent), nil
}

func SendEvent(req *TriggerEvent, apiKey string, timeout time.Duration) error {
	httpClient := http.Client{}

	bodyReader, marshalErr := generateEventBody(req)
	if marshalErr != nil {
		return marshalErr
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc() //need to call at termination to clean up the context

	httpRq, rqErr := http.NewRequestWithContext(ctx, "POST", "https://events.pagerduty.com/v2/enqueue", bodyReader)
	if rqErr != nil {
		return rqErr
	}

	if apiKey != "" {
		httpRq.Header.Add("Authorization", fmt.Sprintf("Token token=%s", apiKey))
	}
	httpRq.Header.Add("Content-Type", "application/json")
	httpRq.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")

	response, err := httpClient.Do(httpRq)
	if err != nil {
		return err
	} else {
		defer response.Body.Close()
		contentBytes, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Printf("ERROR pagerduty.SendEvent could not read server response: %s", readErr)
			contentBytes = []byte("")
		}

		if response.StatusCode == 200 {
			log.Printf("INFO pagerduty.SendEvent Submitted event to PagerDuty")
		} else {
			log.Printf("ERROR pagerduty.SendEvent Pagerduty returned a %d error: %s", response.StatusCode, string(contentBytes))
			return errors.New("server rejected message")
		}
	}
	return nil
}
