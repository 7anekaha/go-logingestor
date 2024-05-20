package main

import (
	"context"
	"fmt"
	"time"

	server "github.com/7anekaha/go-logingestor/api/server"
)

func main() {
	// Create log entries with different timestamps
	log1 := server.Log{
		Level:      "INFO",
		Message:    "This is a log message",
		ResourceId: "1234",
		Timestamp:  parseTime("2023-09-15T08:00:00Z"),
		TraceId:    "trace-1234",
		SpanId:     "span-1234",
		Commit:     "abc123",
		Metadata: server.Metadata{
			ParentResourceId: "5678",
		},
	}
	log2 := server.Log{
		Level:      "ERROR",
		Message:    "This is an error message",
		ResourceId: "5678",
		Timestamp:  parseTime("2023-09-16T09:30:00Z"),
		TraceId:    "trace-5678",
		SpanId:     "span-5678",
		Commit:     "def456",
		Metadata: server.Metadata{
			ParentResourceId: "9012",
		},
	}
	log3 := server.Log{
		Level:      "INFO",
		Message:    "This is another log message",
		ResourceId: "9012",
		Timestamp:  parseTime("2023-09-17T11:45:00Z"),
		TraceId:    "trace-9012",
		SpanId:     "span-9012",
		Commit:     "ghi789",
		Metadata: server.Metadata{
			ParentResourceId: "3456",
		},
	}
	log4 := server.Log{
		Level:      "WARN",
		Message:    "This is a warning message",
		ResourceId: "3456",
		Timestamp:  parseTime("2023-09-18T14:20:00Z"),
		TraceId:    "trace-3456",
		SpanId:     "span-3456",
		Commit:     "jkl012",
		Metadata: server.Metadata{
			ParentResourceId: "7890",
		},
	}
	log5 := server.Log{
		Level:      "DEBUG",
		Message:    "This is a debug message",
		ResourceId: "7890",
		Timestamp:  parseTime("2023-09-19T16:55:00Z"),
		TraceId:    "trace-7890",
		SpanId:     "span-7890",
		Commit:     "mno345",
		Metadata: server.Metadata{
			ParentResourceId: "0123",
		},
	}

	// Connect to the database
	mongoConfig := server.NewMongoConfig(
		"mongodb://nico:secret@localhost:27017",
		10,
		50,
		30,
	)

	client, cancel, err := server.NewMongoClient(mongoConfig)
	if err != nil {
		panic("error connecting to mongodb")
	}

	defer func() {
		cancel()
		if err := client.Disconnect(context.Background()); err != nil {
			panic("mongodb disconnect error")
		}
	}()

	repo := server.NewMongoRepository(client)
	repo.Add(context.Background(), log1)
	repo.Add(context.Background(), log2)
	repo.Add(context.Background(), log3)
	repo.Add(context.Background(), log4)
	repo.Add(context.Background(), log5)
}

// parseTime parses a time string in RFC3339 format and returns a time.Time object
func parseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(fmt.Sprintf("error parsing time: %v", err))
	}
	return t
}
