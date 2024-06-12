package main

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Email struct {
	To      string
	Subject string
	Body    string
}

type Sender struct {
	Email *email.Email
	Auth  smtp.Auth
	From  string
	Addr  string
}

func NewEmail(to, subject, body string) *Email {
	return &Email{To: to, Subject: subject, Body: body}
}

func NewSender(from, password, host, addr string) *Sender {
	email := email.NewEmail()
	email.From = from
	auth := smtp.PlainAuth("", from, password, host)

	return &Sender{Email: email, Auth: auth, From: from, Addr: addr}
}

func (s *Sender) Send(e *Email) error {
	s.Email.To = []string{e.To}
	s.Email.Subject = e.Subject
	s.Email.HTML = []byte(e.Body)

	done := make(chan error)

	go func() {
		done <- s.Email.Send(s.Addr, s.Auth)
	}()

	err := <-done

	return err
}
