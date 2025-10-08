package main

import (
	"fmt"
	"github.com/lib/pq"
	"os"
	"strings"

	"cucinassistant/configs"
	"cucinassistant/database"
)

var daysNames map[string][7]string = map[string][7]string{
	"it": [7]string{"Lunedì", "Martedì", "Mercoledì", "Giovedì", "Venerdì", "Sabato", "Domenica"},
	"en": [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
}

func main() {
	// // Prints a welcome text
	fmt.Print(`CucinAssistant Migration Tool
=============================
This tool MUST BE used once and only once, as it will upgrade the database schema.
It will upgrade it from version 7 (Ciliegia) to version 8 (Banana).
Are you sure to run it? [CONFIRM] `)

	// Asks a confirm
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "CONFIRM" {
		os.Exit(1)
	} else {
		fmt.Println("Confirmed.\n")
	}

	// Initializes all the modules (it will create the missing table)
	configs.LoadAndParse()
	db := database.Connect()
	database.Bootstrap()

	// Prepares the insert statement
	stmt, _ := db.Prepare(`INSERT INTO days (mid, position, name, meals) VALUES ($1, $2, $3, $4);`)

	// Queries the menus
	rows, _ := db.Query(`SELECT m.mid, m.meals, u.email_lang FROM menus m INNER JOIN ca_users u ON m.uid = u.uid;`)

	defer rows.Close()
	for rows.Next() {
		// Parses the data
		var MID int
		var mealsStr string
		var lang string
		rows.Scan(&MID, &mealsStr, &lang)

		meals := strings.Split(mealsStr, ";")

		// Inserts the data into the days table
		for i := 0; i < 7; i++ {
			dayName := daysNames[lang][i]
			dayMeals := []string{meals[2*i], meals[2*i+1]}

			stmt.Exec(MID, i, dayName, pq.Array(&dayMeals))
		}
	}

	// Removes the meals column from the menus tables
	db.Exec(`ALTER TABLE menus DROP COLUMN meals;`)

	fmt.Println("Done.")
}
