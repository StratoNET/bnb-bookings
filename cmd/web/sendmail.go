package main

import (
	"time"

	"github.com/StratoNET/bnb-bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func mailListener() {
	// anonymous, asynchronous function for continuous monitoring of application's mail channel in background
	go func() {
		for {
			msg := <-app.MailChannel
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	}
}
