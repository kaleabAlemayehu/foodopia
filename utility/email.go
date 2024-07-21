package utility

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"sync"
)

type singleton struct {
	smtpHost    string
	smtpAddress string
	from        string
	password    string
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = newClient()
	})
	return instance
}

func newClient() *singleton {
	host := os.Getenv("SMTP_HOST")
	from := os.Getenv("GMAIL_USERNAME")
	password := os.Getenv("GMAIL_PASSWORD")

	return &singleton{
		smtpHost:    host,
		smtpAddress: fmt.Sprintf("%s:587", host),
		from:        from,
		password:    password,
	}
}

func (s *singleton) Send(to string, htmlTemplate string, subject string, data any) error {
	body, err := createBody(htmlTemplate, subject, data)
	if err != nil {
		return err
	}

	return smtp.SendMail(s.smtpAddress,
		smtp.PlainAuth("", s.from, s.password, s.smtpHost),
		s.from, []string{to}, body)
}

func createBody(htmlTemplate string, subject string, data any) ([]byte, error) {

	var body bytes.Buffer

	const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, mime)))
	t, err := template.ParseFiles(htmlTemplate)
	if err != nil {
		log.Panicf("failed to parse '%s' with '%v'", htmlTemplate, err)
		return nil, err
	}
	err = t.Execute(&body, data)
	if err != nil {
		log.Panicf("failed to execute template '%s' with '%v'", htmlTemplate, err)
		return nil, err
	}

	return body.Bytes(), nil
}
