package mail_test

import (
	// "log"
	// "os"
	"testing"

	"mail-sender/config"
	"mail-sender/internal/mail"
	// httpServer "mail-sender/internal/server"
)

func TestMail(t *testing.T) {

	testConf := config.New("config_test", ".")
	// testLogger := log.New(os.Stderr, "", log.Lshortfile)

	// go httpServer.New(testConf, testLogger)

	testEmail := "NeveR5MorE@yandex.ru"
	testBody := "Your transaction has been successfully finished."

	err := mail.SendMail(testEmail, testBody, testConf)
	if err != nil {
		t.Fatalf("/mail.SendMail()/ Mail send failed: \n%v", err)
	}
}
