package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Metadata struct {
	Url string
	Data struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Image struct {
			Url string `json:"url"`
		} `json:"image"`
	} `json:"data"`
}

func main() {
	token := os.Getenv("BOOKMARKS_TELEGRAM_BOT_TOKEN")
	user := os.Getenv("BOOKMARKS_TELEGRAM_BOT_USER")
	var link string

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("[go] Error connecting to Telegram API")
	}

	bot.Debug = false

	log.Printf("[telegram] Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.From.UserName == user {
			log.Printf("[telegram] %s -- %s", update.Message.From.UserName, update.Message.Text)

			log.Print("[go] ACK back to chat")

			chat_res := fmt.Sprintf("[go] Saving the following bookmark:\n %s", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, chat_res)
			bot.Send(msg)

			link = update.Message.Text
			metadata := processLink(link)
			writeToFile(metadata)
		}
	}
}

func processLink(link string) (Metadata) {
	serviceURL := os.Getenv("BOOKMARKS_TELEGRAM_BOT_SERVICE_URL")
	resp, err := http.Get(serviceURL + "?url=" + link)

	if err != nil {
		log.Printf("[go] Error fetching metadata for provided link: %s", link)
	}

	defer resp.Body.Close()

	var metadata Metadata
	metadata.Url = link
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		log.Printf("[go] Error decoding metadata for link: %s", link)
	}

	return metadata
}

func writeToFile(data Metadata) {
  path := "/var/lib/bookmarks-telegram-bot/urls.txt"

	processed_data, err := json.Marshal(data)
	if (err != nil) {
		log.Printf("[go] Error in json Marshal-ing")
	}

	line := string(processed_data)

  f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
		log.Printf("[go] Error opening file")
	}

	defer f.Close()

  if _, err = f.WriteString(line + "\n"); err != nil {
		log.Printf("[go] Error writing to file")
	}

	log.Printf("[go] Wrote %s to file %s", line, path)
}

