package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	Host     string `json:"host" validate:"required"`
	Port     string `json:"port" validate:"required"`
	Identity string `json:"identity"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Mails    []Mail `json:"mails" validate:"required,dive"`
}

type Mail struct {
	From    string `json:"from" validate:"required"`
	To      string `json:"to" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func main() {
	http.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	var server Server

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		response(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error": err.Error(),
		})

		return
	}

	if err := server.validate(); err != nil {
		response(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error": err.Error(),
		})

		return
	}

	server.send()
}

func (server *Server) validate() error {
	var validate *validator.Validate

	validate = validator.New()

	return validate.Struct(server)
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

				log.Println(mail.message())

				<-done
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

func response(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
