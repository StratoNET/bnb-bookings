package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
	encryptionType, _ := strconv.Atoi(os.Getenv("SMTP_ENCRYPTION"))

	server := mail.NewSMTPClient()
	server.Username = os.Getenv("SMTP_USER")
	server.Password = os.Getenv("SMTP_PASS")
	server.Host = os.Getenv("SMTP_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	server.Encryption = mail.Encryption(encryptionType)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		tmpl, err := ioutil.ReadFile(fmt.Sprintf("./templates/email/%s", m.Template))
		if err != nil {
			errorLog.Println(err)
		}
		mailTemplate := string(tmpl)
		content := strings.Replace(mailTemplate, "[%content%]", m.Content, 1)
		email.SetBody(mail.TextHTML, content)
	}

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	}
}
