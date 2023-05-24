package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	maxReconnectAttempts = -1
	maxBackoffDuration   = 10 * time.Minute
)

type WebsocketMessage struct {
	Name  string                 `json:"event"`
	Lobby int                    `json:"lobby"`
	Data  map[string]interface{} `json:"data"`
}

type MessageHandler func(e WebsocketMessage)

var attempt = 1

func ConnectAndProcessWebSocket(url string, handleMessage MessageHandler) error {
	backoff := time.Second

	for {
		err := connectWebSocket(url, handleMessage)
		if err == nil {
			return nil // WebSocket closed gracefully
		}

		logger.WithError(err).WithField("attempt", attempt).Error("Error connecting to WebSocket")

		if attempt == maxReconnectAttempts {
			return fmt.Errorf("maximum reconnection attempts reached")
		}

		time.Sleep(backoff)

		attempt++
		backoff *= 2
		if backoff > maxBackoffDuration {
			backoff = maxBackoffDuration
		}
	}
}

func connectWebSocket(url string, handleMessage MessageHandler) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	logger.Info("Connected to WebSocket")
	attempt = 1 // reset attempts on success
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var data WebsocketMessage
		err = json.Unmarshal(message, &data)
		if err != nil {
			logger.WithError(err).WithField("message", string(message)).Error("Failed to parse JSON message")
			continue
		}

		go handleMessage(data)
	}
}
