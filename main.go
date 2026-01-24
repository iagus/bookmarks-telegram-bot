package main

import (
	"log"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("BOOKMARKS_TELEGRAM_BOT_TOKEN")
	user := os.Getenv("BOOKMARKS_TELEGRAM_BOT_USER")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.From.UserName == user {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			answer := fmt.Sprintf("{ \"url\": \"%s\" }", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)

			bot.Send(msg)
		}
	}
}

