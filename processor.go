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
	var tags []string

	client := http.Client{}
	req, error := http.NewRequest("GET", link, nil)
	if error != nil {
		log.Printf("Error creating request")
		return tags, error
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:153.0) Gecko/20100101 Firefox/153.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")

	res, error := client.Do(req)
	if error != nil {
		log.Printf("[go:fetchMetadata:1] Error fetching link %s", link)
		return tags, error
	}
	defer res.Body.Close()

	log.Printf("[go:fetchMetadata:2] Response status for link %s %d", link, res.StatusCode)

	counter := 0
	b := bufio.NewReader(res.Body)
	for {
		html, error := b.ReadString('>')
		if error != nil {
			log.Printf("[go:fetchMetadata:3] Error reading buffer on line %d", counter)
			log.Printf("%v", error)
		}

		if s.Contains(html, "/head") {
			break
		}

		if s.Contains(html, "<meta") {
			tags = append(tags, html)
		}

		counter++
	}

	log.Printf("[go:fetchMetadata:3] Found %d tags", len(tags))

	return tags, error
}


func process(tags []string) (Metadata, error) {
	var metadata Metadata
	var tag Tag
	var error error

	for _, data := range tags {
		data = normalize(data)

		error = xml.Unmarshal([]byte(data), &tag)
		if error != nil {
			log.Printf("[go:process:1] Error unmarshalling into Tag")
			log.Printf("%v", error)
			break
		}

		if s.Contains(tag.Property, "title") {
			log.Printf("[go:process:2] Found Title meta tag: %s", tag.Content)
			metadata.Data.Title = tag.Content
		} else if s.Contains(tag.Property, "image") {
			log.Printf("[go:process:2] Found Image meta tag: %s", tag.Content)
			metadata.Data.Image.URL = tag.Content
		} else if s.Contains(tag.Property, "description") {
			log.Printf("[go:process:2] Found Description meta tag: %s", tag.Content)
			metadata.Data.Description = tag.Content
		} else {
			log.Printf("[go:process:2] Meta tag %s not targeted. Dropped", tag.Content)
		}
	}

	return metadata, error
}

func normalize(data string) string {
	data = s.Replace(data, ">", "/>", 1)
	data = s.Replace(data, "//>", "/>", 1)

	return data
}

