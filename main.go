package main

import (
	"log"
	"fmt"
	"os"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token string = os.Getenv("BOOKMARKS_TOKEN")
var user string = os.Getenv("BOOKMARKS_USER")
var path string = os.Getenv("BOOKMARKS_PATH")
var service string = os.Getenv("BOOKMARS_SERVICE")

func main() {
	var link string

	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		log.Printf("[go] Error connecting to Telegram API")
	}

	bot.Debug = false

	log.Printf("[telegram] Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.From.UserName == user {
			log.Printf("[telegram] %s -- %s", update.Message.From.UserName, update.Message.Text)

			log.Print("[go] ACK back to chat")
			chat_res := fmt.Sprintf("[go] Saving the following bookmark:\n %s", update.Message.Text)
			msg := telegram.NewMessage(update.Message.Chat.ID, chat_res)
			bot.Send(msg)

			link = update.Message.Text
			metadata, error := handleLink(link)
			if error != nil {
				err_msg := telegram.NewMessage(update.Message.Chat.ID, "[go] Error handling bookmark")
				bot.Send(err_msg)
			} else {
				metadata.Link = link
				writeToFile(metadata)
			}
		}
	}
}

