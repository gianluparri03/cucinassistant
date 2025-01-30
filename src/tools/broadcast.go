package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"strings"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/email"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Prints a welcome text
	welcome := "CucinAssistant Broadcast Email Sender"
	fmt.Println(welcome)
	fmt.Println(strings.Repeat("=", len(welcome)))

	// Reads the configs and connects to the database
	configs.LoadAndParse()
	database.Connect()

	// Reads from stdin both the subject and the body
	fmt.Printf("\nEmail subject\n> ")
	scanner.Scan()
	subject := scanner.Text()
	fmt.Printf("\nEmail body\n> ")
	scanner.Scan()
	body := template.HTML(scanner.Text())

	// Sends a test email
	data := map[string]any{"Body": body}
	email.SendMail(subject, "empty", data, configs.EmailSender)
	fmt.Println("\nSent a test email at " + configs.EmailSender)

	// Asks if it's okay
	fmt.Printf("Type BROADCAST to send it to everyone\n> ")
	scanner.Scan()
	if scanner.Text() == "BROADCAST" {
		// Sends it to everyone
		emails := database.GetUsersEmails()
		email.SendMail(subject, "empty", data, emails...)
		fmt.Println("\nDone.")
	} else {
		fmt.Println("\nExiting.")
	}
}
