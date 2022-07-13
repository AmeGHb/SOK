package mail

import (
	"bytes"
	"log"
	"net/smtp"
)

func SendMail(email, body string) error {

	client, err := smtp.Dial("mail.example.com:25")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Set the sender and recipient.
	client.Mail("sender@example.org")
	client.Rcpt(email)

	// Send the email body.
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()

	buf := bytes.NewBufferString(body)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}

	return nil
}
