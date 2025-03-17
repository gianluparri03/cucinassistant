package email

import (
	"bytes"
	"context"
	"log/slog"
	"net/smtp"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/langs"
)

// Email is the template of an email
type Email struct {
	// Subject is both the email subject and title
	Subject string

	// Content is the text sent with the email
	Content string

	// Raw indicates if Subject and Content are IDs
	// that need to be translated (false) or are already
	// translated (true)
	Raw bool
}

var (
	// Welcome is sent after the signup
	Welcome = Email{Subject: "STR_SIGNUP_DONE", Content: "EMA_WELCOME"}

	// ResetPassword is used to send the user its token
	ResetPassword = Email{Subject: "STR_RESET_PASSWORD", Content: "EMA_RESET_PASSWORD"}

	// PasswordChanged is sent every time the password is changed
	PasswordChanged = Email{Subject: "STR_CHANGE_PASSWORD", Content: "EMA_PASSWORD_CHANGED"}

	// DeleteConfirm is used to send the user its token
	DeleteConfirm = Email{Subject: "STR_DELETE_USER", Content: "EMA_DELETE_CONFIRM"}

	// Goodbye is sent after an account has been deleted
	Goodbye = Email{Subject: "STR_GOODBYE", Content: "EMA_GOODBYE"}
)

// Write executes the email template with the given data.
// It reads from the user the username, the recipient and the language.
func (e Email) Write(user *database.User, data map[string]any) EmailBody {
	// Translates subject and content if the email is not raw
	subject := e.Subject
	content := e.Content
	if !e.Raw {
		lang := langs.Lang(user.EmailLang)
		subject = langs.Translate(lang, subject, nil)
		content = langs.Translate(lang, content, data)
	}

	// Prepares the headers of the body
	var body bytes.Buffer
	body.Write([]byte("Subject: " + subject + "\n"))
	body.Write([]byte("From: CucinAssistant <" + configs.EmailSender + ">\n"))
	body.Write([]byte("To: " + user.Email + "\n"))
	body.Write([]byte("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\n\n"))

	// Adds the subject and the content of the email
	if data == nil {
		data = make(map[string]any)
	}
	data["Username"] = user.Username
	data["Subject"] = subject
	data["Content"] = content

	// Executes the templates
	Base(subject, content, user.Username).Render(context.Background(), &body)
	return EmailBody{Body: &body, Recipient: user.Email}
}

// EmailBody is an email ready to be sent
type EmailBody struct {
	// Body contains the email
	Body *bytes.Buffer

	// Recipient is the recipient of the email
	Recipient string
}

// Send sends the body to the recipient.
// If emails are not enabled in the configs, it writes it in the terminal instead
// of sending it.
func (b EmailBody) Send() {
	// Prepares the credentials
	credentials := smtp.PlainAuth(
		"",
		configs.EmailLogin,
		configs.EmailPassword,
		configs.EmailServer,
	)

	// Prints the mesasge in the console if emails aren't enabled
	if !configs.EmailEnabled {
		slog.Warn("--- [Begin Email] ---\n" + string(b.Body.Bytes()) + "\n--- [End Email] ---")
		return
	}

	// Sends the message
	err := smtp.SendMail(
		configs.EmailServer+":"+configs.EmailPort,
		credentials,
		configs.EmailSender,
		[]string{b.Recipient},
		b.Body.Bytes(),
	)

	// Checks for errors
	if err != nil {
		slog.Error("while sending email:", "err", err)
	} else {
		slog.Debug("Sent email", "to", b.Recipient)
	}
}
