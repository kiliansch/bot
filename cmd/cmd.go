package cmd

import (
	"log"
	"net/http"

	"github.com/kmptnz/bot/internal/config"
	_ "github.com/kmptnz/bot/internal/db"
	"github.com/kmptnz/bot/internal/logger"
	"github.com/kmptnz/bot/internal/matcher"
	"github.com/kmptnz/bot/internal/telegram"
)

// Create an HTTP server listening for webhook requests from Telegram on port 3000
func RunServer() {
	telegram.SetWebhookURL()

	log.Fatal(http.ListenAndServe(config.Cfg.Server.ListenAddr, http.HandlerFunc(HandleWebhookRequest)))
}

// Handle a webhook request sent by Telegram
func HandleWebhookRequest(res http.ResponseWriter, req *http.Request) {
	// Parse the request
	requestMessage, err := telegram.ParseRequest(res, req)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	// Send request message to matchers
	matcher.ExecuteMatchers(*requestMessage)
}
