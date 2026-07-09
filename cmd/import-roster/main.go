package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AlumniRosterEntry struct {
	FullName  string
	ClassYear int
	House     string
	Phone     string
	Email     string
}

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: import-roster <csv-file> <database-url>")
		fmt.Println("Example: import-roster roster.csv \"postgres://user:pass@localhost:5432/dbname?sslmode=disable\"")
		os.Exit(1)
	}

	csvFile := os.Args[1]
	dbURL := os.Args[2]

	// Open CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	if len(records) < 2 {
		log.Fatal("CSV file is empty or has no data rows")
	}

	// Parse header
	header := records[0]
	// Expected columns: full_name, class_year, house, phone, email
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	// Validate required columns
	requiredCols := []string{"full_name", "class_year"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			log.Fatalf("Missing required column: %s", col)
		}
	}

	// Parse data rows
	var entries []AlumniRosterEntry
	for i, record := range records[1:] {
		if len(record) != len(header) {
			log.Printf("Skipping row %d: column count mismatch", i+2)
			continue
		}

		fullName := record[colIndex["full_name"]]
		classYearStr := record[colIndex["class_year"]]
		
		classYear, err := strconv.Atoi(classYearStr)
		if err != nil {
			log.Printf("Skipping row %d: invalid class year: %s", i+2, classYearStr)
			continue
		}

		entry := AlumniRosterEntry{
			FullName:  fullName,
			ClassYear: classYear,
		}

		if houseIdx, ok := colIndex["house"]; ok {
			entry.House = record[houseIdx]
		}
		if phoneIdx, ok := colIndex["phone"]; ok {
			entry.Phone = record[phoneIdx]
		}
		if emailIdx, ok := colIndex["email"]; ok {
			entry.Email = record[emailIdx]
		}

		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		log.Fatal("No valid entries found in CSV")
	}

	fmt.Printf("Parsed %d entries from CSV\n", len(entries))

	// Connect to database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database")

	// Insert entries
	inserted := 0
	for i, entry := range entries {
		query := `
			INSERT INTO alumni_roster (full_name, class_year, house, phone, email)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING
		`

		var house, phone, email *string
		if entry.House != "" {
			house = &entry.House
		}
		if entry.Phone != "" {
			phone = &entry.Phone
		}
		if entry.Email != "" {
			email = &entry.Email
		}

		result, err := db.Exec(query, entry.FullName, entry.ClassYear, house, phone, email)
		if err != nil {
			log.Printf("Failed to insert entry %d: %v", i+1, err)
			continue
		}

		rows, _ := result.RowsAffected()
		if rows > 0 {
			inserted++
		}
	}

	fmt.Printf("Successfully imported %d entries\n", inserted)
	fmt.Printf("Skipped %d duplicate entries\n", len(entries)-inserted)
}
