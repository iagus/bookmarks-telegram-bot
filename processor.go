package main

import(
	"log"
	"net/http"
	"encoding/xml"
	"bufio"
	s "strings"
)

type Tag struct {
	Property string `xml:"property,attr"`
	Content string `xml:"content,attr"`
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
