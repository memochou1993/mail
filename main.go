package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	host     string
	port     string
	identity string
	username string
	password string
	mails    []Mail
}

type Mail struct {
	from    string
	to      string
	subject string
	body    string
}

func main() {
	server := Server{
		host:     os.Getenv("MAIL_HOST"),
		port:     os.Getenv("MAIL_PORT"),
		identity: "",
		username: os.Getenv("MAIL_USERNAME"),
		password: os.Getenv("MAIL_PASSWORD"),
		mails: []Mail{
			{
				from:    "Memo Chou",
				to:      "memochou1993@gmail.com",
				subject: "This is an example email",
				body:    "Hello",
			},
		},
	}

	server.send()
}

func (server *Server) send() {
	// TODO: should use channel
	for _, mail := range server.mails {
		err := smtp.SendMail(
			fmt.Sprintf("%s:%s", server.host, server.port),
			smtp.PlainAuth(server.identity, server.username, server.password, server.host),
			server.username,
			[]string{mail.to},
			[]byte(mail.message()),
		)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (mail *Mail) message() string {
	headers := map[string]string{
		"From":    mail.from,
		"To":      mail.to,
		"Subject": mail.subject,
	}

	message := ""

	for header, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", header, value)
	}

	return message + "\r\n" + mail.body
}
