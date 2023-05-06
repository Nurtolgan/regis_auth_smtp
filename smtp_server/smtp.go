package smtp_server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"manager/config"
	"manager/debugger"
)

type mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func SendMail(hash, to string) error {
	var b bytes.Buffer
	generateHTML("/home/example.html").Execute(&b, hash)

	auth := smtp.PlainAuth("", config.Config.Smtp_from, config.Config.Smtp_password, config.Config.Smtp_host)
	request := mail{
		Sender:  config.Config.Smtp_from,
		To:      []string{to},
		Subject: "Registration to FVRK",
		Body:    b.String()}
	msg := buildMessage(request)
	return smtp.SendMail(fmt.Sprintf("%s:%d", config.Config.Smtp_host, config.Config.Smtp_port), auth, config.Config.Smtp_from, request.To, msg)
}

func buildMessage(mail mail) []byte {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", mail.To[0])
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return []byte(msg)
}

func generateHTML(temp string) *template.Template {
	userTemplate, err := template.ParseFiles(temp)
	debugger.CheckError("Parse Files", err)
	return userTemplate
}
