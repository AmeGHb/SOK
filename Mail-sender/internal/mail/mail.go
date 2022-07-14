package mail

import (
	"crypto/tls"

	"mail-sender/config"

	"gopkg.in/gomail.v2"
)

func SendMail(email, body string, conf *config.Config) error {

	message := gomail.NewMessage()
	message.SetHeader("From", conf.MailFrom)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Transaction inforamtion.")
	message.SetBody("text/plain", body)

	d := gomail.NewDialer(
		conf.GetMailHost(),
		conf.GetMailPort(),
		conf.GetMailFrom(),
		conf.GetMailPassword(),
	)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
