package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/lib/pq"
)

type Event struct {
	Id        string `json:"id"`
	EventType string `json:"event_type"`
	PageURL   string `json:"page_url"`
	Timestamp string `json:"timestamp"`
	UserID    string `json:"user_id"`
}

func main() {
	// Connect to PostgreSQL
	connStr := ""
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connect to Elasticsearch
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "elastic",
		Password:  "",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Query events from PostgreSQL
	rows, err := db.Query("SELECT id, event_type, page_url, timestamp, user_id FROM events")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	start := time.Now()
	fmt.Println("Starting indexing")
	// Index events into Elasticsearch
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.EventType, &event.PageURL, &event.Timestamp, &event.UserID)
		if err != nil {
			log.Fatal(err)
		}

		// Convert event to JSON
		eventJSON, err := json.Marshal(event)
		if err != nil {
			log.Fatal(err)
		}

		// Index the event
		req := esapi.IndexRequest{
			Index:      "events",
			DocumentID: event.Id,
			Body:       strings.NewReader(string(eventJSON)),
			Refresh:    "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error indexing document: %s", err)
		}
		defer res.Body.Close()

		fmt.Printf("Indexed event: %s\n", event.Id)
	}

	fmt.Printf("Indexed events in %s\n", time.Since(start))

	fmt.Println("All events indexed into Elasticsearch.")
}
