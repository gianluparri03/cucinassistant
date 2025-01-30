package email

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/smtp"

	"cucinassistant/configs"
)

// SendMail sends an email to one (or more) recipients, with a given subject, whose
// content is generated from a template with some additional data.
// TemplateName must contains only the basename of the file.
func SendMail(subject string, templateName string, data map[string]any, recipients ...string) {
	// Prepares the headers of the body
	var body bytes.Buffer
	body.Write([]byte("Subject: " + subject + "\n"))
	body.Write([]byte("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n\n"))

	// Writes the actual body
	formatMessage(&body, templateName, data)

	// Sends it to the recipient
	sendBody(&body, templateName, recipients...)
}

// formatMessage executes the template with the given data, writing all
// to the buffer
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

	// Formats the template
	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error("while executing email template:", "err", err, "template", templateName)
	}
}

// sendBody sends the message (a bytes.Buffer) to the recipients
func sendBody(body *bytes.Buffer, emailType string, recipients ...string) {
	// Prepares the credentials
	credentials := smtp.PlainAuth(
		"",
		configs.EmailLogin,
		configs.EmailPassword,
		configs.EmailServer,
	)

	for _, recipient := range recipients {
		// Sends the message (or prints it in the console)
		if configs.EmailEnabled {
			err := smtp.SendMail(
				configs.EmailServer+":"+configs.EmailPort,
				credentials,
				configs.EmailSender,
				[]string{recipient},
				body.Bytes(),
			)

			// Checks for errors
			if err != nil {
				slog.Error("while sending email:", "err", err)
			} else {
				slog.Debug("Sent email:", "emailType", emailType, "recipient", recipient)
			}
		} else {
			slog.Warn("--- [Begin Email] ---\n" + string(body.Bytes()) + "\n--- [End Email] ---")
		}
	}
}
