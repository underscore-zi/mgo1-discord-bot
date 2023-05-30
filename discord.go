package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var DiscordWebhookURL string

type DiscordMentionKey string

const (
	MentionParse DiscordMentionKey = "parse"
	//MentionRoles DiscordMentionKey = "roles"
	//MentionUsers DiscordMentionKey = "users"
)

type DiscordWebhookMessage struct {
	Content         string                         `json:"content"`
	AllowedMentions map[DiscordMentionKey][]string `json:"allowed_mentions,omitempty"`
}

func ExecuteDiscordWebhook(message string) error {
	payload := DiscordWebhookMessage{
		Content: message,
		AllowedMentions: map[DiscordMentionKey][]string{
			MentionParse: {},
		},
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
