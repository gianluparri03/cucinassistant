package email

import (
	"bytes"
	"log/slog"
	"net/smtp"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/langs"
)

// Email is the template of an email
type Email struct {
	// Subject is both the email subject and title
	Subject langs.String

	// Content is the text sent with the email
	Content langs.String
}

var (
	// Welcome is sent after the signup
	Welcome = Email{Subject: langs.STR_SIGNUP_DONE, Content: langs.STR_WELCOME_EMAIL}

	// ResetPassword is used to send the user its token
	ResetPassword = Email{Subject: langs.STR_RESET_PASSWORD, Content: langs.STR_RESET_PASSWORD_EMAIL}

	// PasswordChanged is sent every time the password is changed
	PasswordChanged = Email{Subject: langs.STR_CHANGE_PASSWORD, Content: langs.STR_PASSWORD_CHANGED_EMAIL}

	// DeleteConfirm is used to send the user its token
	DeleteConfirm = Email{Subject: langs.STR_DELETE_USER, Content: langs.STR_DELETE_CONFIRM_EMAIL}

	// Goodbye is sent after an account has been deleted
	Goodbye = Email{Subject: langs.STR_GOODBYE, Content: langs.STR_GOODBYE_EMAIL}
)

// Write creates an EmailBody.
// It reads from the user the username, the recipient and the language.
func (e Email) Write(user *database.User, link string) EmailBody {
	ctx := langs.Get(user.EmailLang).Ctx()

	return RawEmail{
		Subject: langs.Translate(ctx, e.Subject),
		Content: langs.Translate(ctx, e.Content),
	}.Write(user, link, false)
}

// RawEmail is an email that has already been translated
type RawEmail struct {
	// Subject is both the email subject and title
	Subject string

	// Content is the text sent with the email
	Content string

	// Newsletter indicates whether to put the unsubscribe link at the bottom
	Newsletter bool
}

// Write executes the email template with the given data.
// It reads from the user the username, the recipient and the language.
func (e RawEmail) Write(user *database.User, link string, newsletter bool) EmailBody {
	var disableUrl string

	// Prepares the headers of the body
	var body bytes.Buffer
	body.Write([]byte("Subject: " + e.Subject + "\n"))
	body.Write([]byte("From: CucinAssistant <" + configs.EmailSender + ">\n"))
	body.Write([]byte("To: " + user.Email + "\n"))
	if newsletter && user.Newsletter != nil {
		disableUrl = configs.BaseURL + "/disable_newsletter?token=" + *user.Newsletter
		body.Write([]byte("List-Unsubscribe: <" + disableUrl + ">\n"))
		body.Write([]byte("List-Unsubscribe-Post: List-Unsubscribe=One-Click\n"))
	}
	body.Write([]byte("MIME-version: 1.0\nContent-Type: text/html; charset=\"UTF-8\"\n\n\n\n"))

	// Executes the templates
	Base(e.Subject, e.Content, user.Username, link, disableUrl).
		Render(langs.Get(user.EmailLang).Ctx(), &body)
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
