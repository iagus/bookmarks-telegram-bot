package main

import(
	"testing"
	"regexp"
)

func TestProcessLink(t *testing.T) {
	url := "https://www.youtube.com/watch?v=LVUGbW8BnRw"
	metadata, _ := processLink(url)
	wanted_title := regexp.MustCompile(`Our Old Man Is Getting Old`)
	wanted_img := regexp.MustCompile(`https://i.ytimg.com/vi/LVUGbW8BnRw/maxresdefault.jpg`)

	if !wanted_title.MatchString(metadata.Data.Title) {
		t.Errorf("Title mismatch")
	}
	if !wanted_img.MatchString(metadata.Data.Image.URL) {
		t.Errorf("Image URL mismatch")
	}
}

