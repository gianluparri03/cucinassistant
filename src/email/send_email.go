package email

import (
	"bytes"
	"html/template"
	"log"
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
		log.Print("Sent email (" + templateName + ") to <" + recipient + ">")
	}
}

func formatMessage(w *bytes.Buffer, templateName string, data map[string]any) {
	// Fetches the templates
	tmpl := template.Must(template.ParseFiles(
		"email/templates/base.html",
		"email/templates/"+templateName+".html",
	))

	// Adds the banner path
	data["banner"] = config.Runtime.BaseURL + "/assets/banner.png"

	// Formats the template
	tmpl.Execute(w, data)
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
	smtp.SendMail(
		config.Runtime.Email.Server+":"+config.Runtime.Email.Port,
		credentials,
		config.Runtime.Email.Address,
		[]string{recipient},
		body.Bytes(),
	)
}
