package main

import(
	"testing"
	"regexp"
)

func TestProcess(t *testing.T) {
	tags := []string{
		"<meta property=\"og:title\" content=\"OG Title\"/>",
		"<meta property=\"og:image\" content=\"OG Image\"/>",
		"<meta property=\"og:description\" content=\"OG Description\"/>",
	}

	metadata, _ := process("https://example.com", tags)

	wanted_link := regexp.MustCompile(`https://example.com`)
	wanted_title := regexp.MustCompile(`OG Title`)
	wanted_img := regexp.MustCompile(`OG Image`)
	wanted_desc := regexp.MustCompile(`OG Description`)

	if !wanted_link.MatchString(metadata.Link) {
		t.Errorf("Link mismatch")
	}
	if !wanted_title.MatchString(metadata.Data.Title) {
		t.Errorf("Title mismatch")
	}
	if !wanted_img.MatchString(metadata.Data.Image.URL) {
		t.Errorf("Image URL mismatch")
	}
	if !wanted_desc.MatchString(metadata.Data.Description) {
		t.Errorf("Description mismatch")
	}
}

