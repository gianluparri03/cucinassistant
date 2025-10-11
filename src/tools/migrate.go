package main

import (
	"crypto/rand"
	"fmt"
	"os"

	"cucinassistant/configs"
	"cucinassistant/database"
)

func main() {
	// Prints a welcome text
	fmt.Print(`CucinAssistant Migration Tool
=============================
This tool MUST BE used once and only once, as it will upgrade the database schema.
It will upgrade it from version 8 (Banana) to version 9 (Maracuja).
Are you sure to run it? [CONFIRM] `)

	// Asks a confirm
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "CONFIRM" {
		os.Exit(1)
	} else {
		fmt.Println("Confirmed.\n")
	}

	// Initializes all the modules
	configs.LoadAndParse()
	db := database.Connect()
	database.Bootstrap()

	// Queries the users with newsletter enabled
	rows, _ := db.Query(`SELECT uid FROM ca_users WHERE newsletter;`)
	defer rows.Close()
	var users []int
	for rows.Next() {
		var UID int
		rows.Scan(&UID)
		users = append(users, UID)
	}

	// Removes the old newsletter column
	db.Exec(`ALTER TABLE ca_users DROP COLUMN newsletter;`)

	// Creates the new newsletter column
	db.Exec(`ALTER TABLE ca_users ADD COLUMN newsletter CHAR(16);`)

	// Repopulates the newsletter column
	stmt, _ := db.Prepare(`UPDATE ca_users SET newsletter=$2 WHERE uid=$1;`)
	for _, u := range users {
		buffer := make([]byte, 8)
		rand.Read(buffer)
		t := fmt.Sprintf("%x", buffer)

		stmt.Exec(u, t)
	}

	// Sets the UNIQUE constraint
	fmt.Println(db.Exec(`ALTER TABLE ca_users ADD CONSTRAINT ca_user_newsletter_unique UNIQUE (newsletter);`))

	// Set the email_lang field nullable
	db.Exec(`ALTER TABLE ca_users ALTER email_lang DROP NOT NULL;`)
	db.Exec(`UPDATE ca_users SET email_lang=NULL WHERE email_lang='';`)

	fmt.Println("Done.")
}
