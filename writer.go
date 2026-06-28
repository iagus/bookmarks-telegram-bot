package main

import(
	"encoding/json"
	"os"
	"log"
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

