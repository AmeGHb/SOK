package mail

import (
	"mail-sender/config"

	"gopkg.in/gomail.v2"
)

func SendMail(email, body string, conf *config.Config) error {

	message := gomail.NewMessage()
	message.SetHeader("From", conf.MailFrom)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Transaction inforamtion.")
	message.SetBody("text/plain", body)

	d := gomail.Dialer{Host: conf.MailHost, Port: conf.MailPort}
	if err := d.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
