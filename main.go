package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vl/habit-cli/internal/app"
	"github.com/vl/habit-cli/internal/shared/db"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Seed database with test data")
	flag.Parse()

	// Open database
	dbPath := db.DefaultPath()
	database, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Seed test data if requested
	if *seedFlag {
		if err := seedTestData(database); err != nil {
			fmt.Fprintf(os.Stderr, "Error seeding data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Test data seeded successfully!")
	}

	// Create and run the application
	model := app.New(database)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

func seedTestData(database *db.DB) error {
	// Create categories
	categories := []struct {
		name  string
		color string
	}{
		{"Health", "#10B981"},
		{"Work", "#3B82F6"},
		{"Personal", "#F59E0B"},
	}

	catIDs := make(map[string]int64)
	for _, cat := range categories {
		result, err := database.DB.Exec(
			"INSERT OR IGNORE INTO categories (name, color) VALUES (?, ?)",
			cat.name, cat.color,
		)
		if err != nil {
			return err
		}
		id, _ := result.LastInsertId()
		if id == 0 {
			// Already exists, get ID
			row := database.DB.QueryRow("SELECT id FROM categories WHERE name = ?", cat.name)
			row.Scan(&id)
		}
		catIDs[cat.name] = id
	}

	// Create habits
	habits := []struct {
		name          string
		description   string
		freqType      string
		freqValue     int
		category      string
		completions   int // days back to add completions
	}{
		{"Morning Exercise", "30 min workout", "daily", 1, "Health", 12},
		{"Read", "Read for 20 minutes", "daily", 1, "Personal", 8},
		{"Meditate", "10 min meditation", "daily", 1, "Health", 15},
		{"Weekly Review", "Review goals and progress", "weekly", 1, "Work", 3},
		{"Learn Something", "Study or take a course", "times_per_week", 3, "Personal", 5},
		{"Drink Water", "8 glasses of water", "daily", 1, "Health", 10},
	}

	for _, h := range habits {
		catID := catIDs[h.category]
		result, err := database.DB.Exec(
			`INSERT OR IGNORE INTO habits (name, description, frequency_type, frequency_value, category_id, created_at)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			h.name, h.description, h.freqType, h.freqValue, catID, time.Now(),
		)
		if err != nil {
			return err
		}

		habitID, _ := result.LastInsertId()
		if habitID == 0 {
			// Already exists, get ID
			row := database.DB.QueryRow("SELECT id FROM habits WHERE name = ?", h.name)
			row.Scan(&habitID)
		}

		// Add completions for past days
		for i := 0; i < h.completions; i++ {
			date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			database.DB.Exec(
				"INSERT OR IGNORE INTO completions (habit_id, completed_at) VALUES (?, ?)",
				habitID, date,
			)
		}
	}

	return nil
}
