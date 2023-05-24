package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var DiscordWebhookURL string

type DiscordWebhookMessage struct {
	Content string `json:"content"`
}

func ExecuteDiscordWebhook(message string) error {
	payload := DiscordWebhookMessage{
		Content: message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", DiscordWebhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	return nil
}
