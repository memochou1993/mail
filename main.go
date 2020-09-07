package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type Server struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Identity string `json:"identity"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mails    []Mail `json:"mails" binding:"required,dive"`
}

type Mail struct {
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		server := Server{}

		if err := c.ShouldBindJSON(&server); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		index()
	})

	r.Run()
}

func index() {
	server := Server{
		Host:     os.Getenv("MAIL_HOST"),
		Port:     os.Getenv("MAIL_PORT"),
		Identity: "",
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Mails: []Mail{
			{
				From:    "Memo Chou",
				To:      "memochou1993@gmail.com",
				Subject: "This is an example email",
				Body:    "Hello",
			},
		},
	}

	server.send()
}

func (server *Server) send() {
	go func() {
		done := make(chan bool, 10)

		for _, mail := range server.Mails {
			done <- true

			go func(mail Mail) {
				err := smtp.SendMail(
					fmt.Sprintf("%s:%s", server.Host, server.Port),
					smtp.PlainAuth(server.Identity, server.Username, server.Password, server.Host),
					server.Username,
					[]string{mail.To},
					[]byte(mail.message()),
				)

				if err != nil {
					log.Fatal(err)
				}

				<- done

				log.Println(mail.message())
			}(mail)
		}
	}()
}

func (mail *Mail) message() string {
	headers := map[string]string{
		"From":    mail.From,
		"To":      mail.To,
		"Subject": mail.Subject,
	}

	message := ""

	for header, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", header, value)
	}

	return message + "\r\n" + mail.Body
}
