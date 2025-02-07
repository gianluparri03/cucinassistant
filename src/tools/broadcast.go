package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/langs"
)

func main() {
	// Prints a welcome text
	welcome := "CucinAssistant Broadcast Email Sender"
	fmt.Println(welcome)
	fmt.Println(strings.Repeat("=", len(welcome)))

	// Initializes all the modules
	configs.LoadAndParse()
	database.Connect()
	langs.Load()

	// Initialize the scanner and prepares
	// the bodies ararray
	scanner := bufio.NewScanner(os.Stdin)
	var bodies []email.EmailBody

	// For each (needed) language
	for lang, users := range getUsers() {
		// Reads the email from stdin
		fmt.Printf("\n[Language: %s]\n", lang)
		email := readEmail(scanner)

		// And writes the email for each user
		for _, user := range users {
			bodies = append(bodies, email.Write(user, nil))
		}
	}

	// Asks a confirm before sending the bodies
	fmt.Println("\n\nEmails ready to be sent.")
	if !confirm(scanner) {
		fmt.Println("\nAborted.")
		os.Exit(1)
	}

	// Sends them
	for _, body := range bodies {
		body.Send()
	}

	fmt.Println("Done.")
}

func getUsers() map[string][]*database.User {
	users := make(map[string][]*database.User)

	// Adds each user in a group, based on the
	// EmailLang value
	for _, user := range database.GetUsersForBroadcast() {
		lang := user.EmailLang
		if lang == "" {
			lang = langs.Default
		}

		users[lang] = append(users[lang], &user)
	}

	return users
}

func readEmail(scanner *bufio.Scanner) email.Email {
	for true {
		// Reads from stdin the subject
		fmt.Printf("Subject: ")
		scanner.Scan()
		subject := scanner.Text()

		// Reads from stdin the content
		fmt.Printf("Content: ")
		scanner.Scan()
		content := scanner.Text()

		// Repeat this process every time is needed
		if confirm(scanner) {
			return email.Email{Subject: subject, Content: content, Raw: true}
		}
	}

	return email.Email{}
}

func confirm(scanner *bufio.Scanner) bool {
	fmt.Printf("Confirm? [y/n] ")
	scanner.Scan()
	return scanner.Text() == "y"
}
