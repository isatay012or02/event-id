package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.homebank.kz/homebank-oauth/halykid-events/models"
	"git.homebank.kz/homebank-oauth/halykid-events/service"
	"net/http"
)

var Hook WebHook

type WebHook struct {
	log service.AppLogger
}

func (wk *WebHook) WebhookHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		return
	}

	var incomingData models.AuthEvent
	if err := json.NewDecoder(r.Body).Decode(&incomingData); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wk *WebHook) SendToWebApp(data models.AuthEvent, path string) error {

	webAppURL := path

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(webAppURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
