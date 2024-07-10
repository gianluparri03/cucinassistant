package email

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/smtp"

	"cucinassistant/config"
)

func SendMail(recipient string, subject string, templateName string, data map[string]any) {
	// Prepares the headers of the body
	var body bytes.Buffer
	body.Write([]byte("Subject: " + subject + "\n"))
	body.Write([]byte("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n\n"))

	// Writes the actual body
	formatMessage(&body, templateName, data)

	// Sends it to the recipient
	if config.Runtime.Email.Enabled {
		sendBody(recipient, &body)
		slog.Info("Sent email", "template", templateName, "recipient", recipient)
	}
}

func formatMessage(w *bytes.Buffer, templateName string, data map[string]any) {
	// Fetches the templates
	tmpl, err := template.ParseFiles(
		"email/templates/base.html",
		"email/templates/"+templateName+".html",
	)

	// Checks for errors
	if err != nil {
		slog.Error("while fetching email template:", "err", err, "template", templateName)
		return
	}

	// Adds the banner path
	data["banner"] = config.Runtime.BaseURL + "/assets/banner.png"

	// Formats the template
	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error("while executing email template:", "err", err, "template", templateName)
	}
}

func sendBody(recipient string, body *bytes.Buffer) {
	// Prepares the credentials
	credentials := smtp.PlainAuth(
		"",
		config.Runtime.Email.Login,
		config.Runtime.Email.Password,
		config.Runtime.Email.Server,
	)

	// Sends the message
	err := smtp.SendMail(
		config.Runtime.Email.Server+":"+config.Runtime.Email.Port,
		credentials,
		config.Runtime.Email.Address,
		[]string{recipient},
		body.Bytes(),
	)

	// Checks for errors
	if err != nil {
		slog.Error("while sending email:", "err", err)
	}
}
