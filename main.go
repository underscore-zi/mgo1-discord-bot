package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
	"sync"
)

var logger = logrus.StandardLogger()

var URLFormatString string
var games = make(map[int]map[string]interface{})
var gamesLock = sync.Mutex{}

func GetGameLink(gameID int) string {
	if URLFormatString == "" {
		return "'"
	}
	return fmt.Sprintf(URLFormatString, gameID)
}

func newGame(lobby int, data map[string]interface{}) {
	logger.WithFields(logrus.Fields{
		"lobby": lobby,
		"name":  data["name"],
		"id":    data["game_id"],
	}).Info("New game")

	gameID := int(data["game_id"].(float64))
	message := strings.Builder{}
	message.WriteString(data["host"].(string))
	message.WriteString(" is hosting: ")
	message.WriteString(data["name"].(string))
	if data["has_password"].(bool) {
		message.WriteString(" (Locked)")
	}

	gameLink := GetGameLink(gameID)
	if gameLink != "" {
		message.WriteString(" - ")
		message.WriteString(gameLink)
	}

	go func() { _ = ExecuteDiscordWebhook(message.String()) }()

	gamesLock.Lock()
	games[gameID] = data
	gamesLock.Unlock()
}

func removeGame(lobby int, data map[string]interface{}) {
	gameID := int(data["game_id"].(float64))
	if info, ok := games[gameID]; ok {
		logger.WithFields(logrus.Fields{
			"lobby": lobby,
			"name":  info["name"],
			"id":    info["game_id"],
		}).Info("Game closed")

		message := strings.Builder{}
		message.WriteString(info["host"].(string))
		message.WriteString(" has closed: ")
		message.WriteString(info["name"].(string))

		go func() { _ = ExecuteDiscordWebhook(message.String()) }()
	}

	gamesLock.Lock()
	delete(games, gameID)
	gamesLock.Unlock()
}

func main() {
	discord := flag.String("discord", "", "Discord webhook URL")
	socket := flag.String("socket", "", "WebSocket URL")
	urlFormat := flag.String("url", "", "URL Format String for game links (optional)")

	flag.Parse()

	if v, found := os.LookupEnv("DISCORD_WEBHOOK"); found {
		if err := flag.Set("discord", v); err != nil {
			logger.WithError(err).Fatal("Error setting Discord webhook URL from environment")
		}
	}

	if v, found := os.LookupEnv("WEBSOCKET_URL"); found {
		if err := flag.Set("socket", v); err != nil {
			logger.WithError(err).Fatal("Error setting WebSocket URL from environment")
		}
	}

	if v, found := os.LookupEnv("URL_FORMAT_STR"); found {
		if err := flag.Set("url", v); err != nil {
			logger.WithError(err).Fatal("Error setting URL Format String from environment")
		}
	}

	if *discord == "" || *socket == "" {
		flag.Usage()
		os.Exit(1)
	}

	DiscordWebhookURL = *discord
	URLFormatString = *urlFormat

	err := ConnectAndProcessWebSocket(*socket, func(e WebsocketMessage) {
		switch e.Name {
		case "game_created":
			newGame(e.Lobby, e.Data)
		case "game_deleted":
			removeGame(e.Lobby, e.Data)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
