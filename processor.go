package main

import(
	"log"
	"net/http"
	"encoding/xml"
  "bufio"
	s "strings"
)

const	Title = "og:title"
const	Img = "og:image"
const	Desc = "og:description"

type Tag struct {
	Property string `xml:"property,attr"`
	Content string `xml:"content,attr"`
}

func handleLink(link string) (Metadata, error) {
	var m Metadata
	tags, error := fetchMetadata(link)
	if error != nil {
		log.Printf("[go:handleLink:1] Error handling link")
	} else {
		m, error = process(tags)
		if error != nil {
			log.Printf("[go:handleLink:2] Error processing link")
		}
	}

	return m, error
}

func fetchMetadata(link string) ([]string, error) {
	res, error := http.Get(link)
	var tags []string
	if error != nil {
		log.Printf("[go:fetchMetadata:1] Error fetching link %s", link)
	} else {
		defer res.Body.Close()
		log.Printf("[go:fetchMetadata:2] Response status for link %s %d", link, res.StatusCode)
		b := bufio.NewReader(res.Body)
		for {
			html, error := b.ReadString('>')
			if (error != nil) {
				break
			} else {
				if s.Contains(html, "/head") {
					break
				} else if s.Contains(html, Title) || s.Contains(html, Img) || s.Contains(html, Desc) {
					html = s.Replace(html, ">", "/>", 1) // formatting for valid xml
					tags = append(tags, html)
				}
			}
		}
	}

	return tags, error
}

func process(tags []string) (Metadata, error) {
	var metadata Metadata
	var tag Tag
	var error error

	for _, data := range tags {
		error = xml.Unmarshal([]byte(data), &tag)
		if error != nil {
			log.Printf("[go:process:1] Error unmarshalling into Tag")
			break
		}

		switch tag.Property {
		case Title:
			log.Printf("[go:process:2] Found Title meta tag: %s", tag.Content)
			metadata.Data.Title = tag.Content
		case Img:
			log.Printf("[go:process:2] Found Image meta tag: %s", tag.Content)
			metadata.Data.Image.URL = tag.Content
		case Desc:
			log.Printf("[go:process:2] Found Description meta tag: %s", tag.Content)
			metadata.Data.Description = tag.Content
		}
	}


	return metadata, error
}

