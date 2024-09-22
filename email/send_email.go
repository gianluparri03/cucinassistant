package email

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/smtp"

	"cucinassistant/config"
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
	if !config.Runtime.Testing {
		sendBody(&body, templateName, recipients...)
	}
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

	// Adds the banner path
	data["banner"] = config.Runtime.BaseURL + "/assets/banner.png"

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
		config.Runtime.Email.Address,
		config.Runtime.Email.Password,
		config.Runtime.Email.Server,
	)

	for _, recipient := range recipients {
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
		} else {
			slog.Debug("Sent email:", "emailType", emailType, "recipient", recipient)
		}
	}
}
