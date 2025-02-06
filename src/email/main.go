package email

import (
	"bytes"
	"log/slog"
	"net/smtp"

	"cucinassistant/configs"
	"cucinassistant/langs"
)

// Email is the template of an email
type Email struct {
	// Subject is both the email subject and title
	Subject string

	// Content is the text sent with the email
	Content string
}

var (
	// Welcome is sent after the signup
	Welcome = Email{Subject: "STR_SIGNUP_DONE"}

	// ResetPassword is used to send the user its token
	ResetPassword = Email{Subject: "STR_RESET_PASSWORD", Content: "EMA_RESET_PASSWORD"}

	// PasswordChanged is sent every time the password is changed
	PasswordChanged = Email{Subject: "STR_CHANGE_PASSWORD", Content: "EMA_PASSWORD_CHANGED"}

	// DeleteConfirm is used to send the user its token
	DeleteConfirm = Email{Subject: "STR_DELETE_USER", Content: "EMA_DELETE_CONFIRM"}

	// Goodbye is sent after an account has been deleted
	Goodbye = Email{Subject: "STR_GOODBYE", Content: "EMA_GOODBYE"}
)

// Write executes the email template with the given language and data
func (e Email) Write(recipient string, lang string, data map[string]any) bytes.Buffer {
	// Prepares the headers of the body
	var body bytes.Buffer
	body.Write([]byte("Subject: " + langs.Translate(lang, e.Subject, nil) + "\n"))
	body.Write([]byte("From: CucinAssistant <" + configs.EmailSender + ">\n"))
	body.Write([]byte("To: " + recipient + "\n"))
	body.Write([]byte("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n\n"))

	// Adds the subject and the content of the email
	if data == nil {
		data = make(map[string]any)
	}
	data["Subject"] = e.Subject
	data["Content"] = e.Content

	// Executes the templates
	langs.ExecuteTemplates(&body, lang, []string{"email/template.html"}, data)
	return body
}

// Send sends an email to a recipient.
// It uses e.Write() and sendBody().
func (e Email) Send(recipient string, lang string, data map[string]any) {
	body := e.Write(recipient, lang, data)
	SendBody(&body, recipient)
}

// SendBody sends the body (a bytes.Buffer) to the recipients.
// If emails are not enabled in the configs, it writes it in the terminal instead
// of sending it.
func SendBody(body *bytes.Buffer, recipients ...string) {
	// Prepares the credentials
	credentials := smtp.PlainAuth(
		"",
		configs.EmailLogin,
		configs.EmailPassword,
		configs.EmailServer,
	)

	for _, recipient := range recipients {
		// Prints the mesasge in the console if emails aren't enabled
		if !configs.EmailEnabled {
			slog.Warn("--- [Begin Email] ---\n" + string(body.Bytes()) + "\n--- [End Email] ---")
			continue
		}

		// Sends the message
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
			slog.Debug("Sent email", "to", recipient)
		}
	}
}
