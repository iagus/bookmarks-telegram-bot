package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"encoding/xml"
	"bufio"
	s "strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token string = os.Getenv("BOOKMARKS_TOKEN")
var user string = os.Getenv("BOOKMARKS_USER")
var path string = os.Getenv("BOOKMARKS_PATH")
var service string = os.Getenv("BOOKMARS_SERVICE")

type Metadata struct {
	Link string
	Data struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Image struct {
			URL string `json:"url"`
		} `json:"image"`
	} `json:"data"`
}

type Tag struct {
	Property string `xml:"property,attr"`
	Content string `xml:"content,attr"`
}

func main() {
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

			chat_res := fmt.Sprintf("Saving the following bookmark:\n %s", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, chat_res)
			bot.Send(msg)

			link = update.Message.Text
			metadata, _ := processLink(link)
			writeToFile(metadata)
		}
	}
}

func processLink(link string) (Metadata, error) {
	res, err := http.Get(link)

	var metadata Metadata

	metadata.Link = link

	if err != nil {
		log.Printf("[go] Error fetching link %s", link)
	} else {
		defer res.Body.Close()

		log.Printf("[go] Response status for link %s %d", link, res.StatusCode)

		b := bufio.NewReader(res.Body)
		total := 0
		for {
			html, err := b.ReadString('>')
			if (total == 3 || err != nil) {
				break
			} else {
				var tag Tag
				if s.Contains(html, "og:title") || s.Contains(html, "og:image\"") || s.Contains(html, "og:description") {
					html = s.Replace(html, ">", "/>", 1)
					err = xml.Unmarshal([]byte(html), &tag)

					if (err != nil) {
						break
					}

					switch tag.Property {
					case "og:title":
						metadata.Data.Title = tag.Content
					case "og:image":
						metadata.Data.Image.URL = tag.Content
					case "og:description":
						metadata.Data.Description = tag.Content
					}

					total = total + 1
				}
			}
		}
	}

	return metadata, err
}

func writeToFile(data Metadata) {
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

