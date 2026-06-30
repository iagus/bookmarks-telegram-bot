package main

import(
	"os"
	"log"
)

func writeToFile(data string) {
	line := string(data)

  f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
		log.Printf("[go:writeToFile:1] Error opening file")
	}

	defer f.Close()

  if _, err = f.WriteString(line + "\n"); err != nil {
		log.Printf("[go:writeToFile:2] Error writing to file")
	}

	log.Printf("[go:writeToFile:3] Wrote %s to file %s", line, path)
}

