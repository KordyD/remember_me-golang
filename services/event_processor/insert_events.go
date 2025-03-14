package main

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	connStr := GetDBCredentials()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	number := 100000
	start := time.Now()
	for i := 1; i <= number; i++ {
		eventType := gofakeit.Sentence(5)
		pageURL := gofakeit.URL()
		timestamp := gofakeit.Date()
		userID := gofakeit.UUID()

		_, err := db.Exec(`
			INSERT INTO events (event_type, page_url, timestamp, user_id)
			VALUES ($1, $2, $3, $4)`,
			eventType, pageURL, timestamp, userID)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Execution took %s\n", time.Since(start))
	fmt.Printf("Inserted %d events into PostgreSQL.", number)
}
