package main

import(
	"log"
	"net/http"
	x "encoding/xml"
	j "encoding/json"
  "bufio"
	s "strings"
)

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

func handleLink(link string) (string, error) {
	var m Metadata
	var l string
	tags, error := fetchMetadata(link)

	if error != nil {
		log.Printf("[go:handleLink:1] Error handling link")
	} else {
		m, error = process(link, tags)
		l, error = serialize(m)
		if error != nil {
			log.Printf("[go:handleLink:2] Error processing link")
		}
	}

	return l, error
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
			break
		}
		counter++
		if s.Contains(html, "<meta") {
			tags = append(tags, html)
		} else if s.Contains(html, "/head") {
			break
		}
	}

	log.Printf("[go:fetchMetadata:3] Found %d tags", len(tags))
	return tags, error
}


func process(link string, tags []string) (Metadata, error) {
	var metadata Metadata
	var tag Tag
	var error error

	for _, data := range tags {
		data = normalize(data)
		error = x.Unmarshal([]byte(data), &tag)
		if error != nil {
			log.Printf("[go:process:1] Error unmarshalling into Tag")
			log.Printf("%v", error)
			break
		}
		if s.Contains(tag.Property, "title") {
			metadata.Data.Title = tag.Content
		} else if s.Contains(tag.Property, "image") {
			metadata.Data.Image.URL = tag.Content
		} else if s.Contains(tag.Property, "description") {
			metadata.Data.Description = tag.Content
		} else {
			log.Printf("[go:process:2] Meta tag %s not targeted. Dropped", tag.Content)
		}
	}

	metadata.Link = link
	return metadata, error
}

func normalize(data string) string {
	data = s.Replace(data, ">", "/>", 1)
	data = s.Replace(data, "//>", "/>", 1)

	return data
}

func serialize(metadata Metadata) (string, error) {
	data, error := j.Marshal(metadata)
	if error != nil {
		log.Printf("[go:serialize:1] Error marshalling metadata into JSON")
	}
	var l string = string(data)

	return l, error
}

