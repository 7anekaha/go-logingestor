package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	server "github.com/7anekaha/go-logingestor/api/server"
)

func main() {

	levels := []string{"INFO", "ERROR", "WARN", "DEBUG"}
	messages := []string{"This is a message 1", "This is a message 2", "This is a message 3", "This is a debug message"}

	// iterate 100 times to populate the database
	for i := 0; i < 100; i++ {

		// 1000 log entries
		for j := 0; j < 1000; j++ {
			// Create log entries with different timestamps
			logReq := server.Log{
				Level:      levels[rand.Intn(4)],
				Message:    messages[rand.Intn(4)],
				ResourceId: generateRandomStringNumber(4),
				Timestamp:  time.Now(),
				TraceId:    "trace-" + generateRandomStringNumber(4),
				SpanId:     "span-" + generateRandomStringNumber(4),
				Commit:     generateRandomStringNumber(6),
				Metadata: server.Metadata{
					ParentResourceId: generateRandomStringNumber(4),
				},
			}
			jsonLogReq, err := json.Marshal(logReq)
			if err != nil {
				panic("error marshalling log request")
			}

			res, err := http.Post("http://localhost:3000/add", "application/json", bytes.NewBuffer(jsonLogReq))
			if err != nil {
				panic("error posting log request")
			}
			log.Println("Status code: ", res.StatusCode)

			timeToSleep := time.Duration(rand.Intn(200)) * time.Millisecond

			time.Sleep(timeToSleep)
		}
	}

}

func generateRandomStringNumber(n int) string {
	const letters = "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
