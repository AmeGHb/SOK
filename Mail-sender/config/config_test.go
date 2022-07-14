package config_test

import (
	"strconv"
	"testing"

	"mail-sender/config"
)

func TestConfigNew(t *testing.T) {

	newConf := config.New("config", ".\\config\\")
	if newConf == nil {
		t.Fatal("/config.New()/ Expected config but got none")
	}
}

func TestGet(t *testing.T) {

	newConf := config.New("config", ".\\config\\")

	if strconv.Itoa(newConf.HttpPort) != newConf.GetHttpPort() {
		t.Errorf("/config.GetHttpPort()/ Exprcted %d but got %s", newConf.HttpPort, newConf.GetHttpPort())
	}

	if strconv.Itoa(newConf.KafkaPort) != newConf.GetKafkaPort() {
		t.Errorf("/config.GetKafkaPort()/ Exprcted %d but got %s", newConf.KafkaPort, newConf.GetKafkaPort())
	}

	if newConf.MailPort != newConf.GetMailPort() {
		t.Errorf("/config.GetMailPort()/ Exprcted %d but got %d", newConf.MailPort, newConf.GetMailPort())
	}

	if newConf.MailHost != newConf.GetMailHost() {
		t.Errorf("/config.GetMailHost()/ Exprcted %s but got %s", newConf.MailHost, newConf.GetMailHost())
	}

	if newConf.MailFrom != newConf.GetMailFrom() {
		t.Errorf("/config.GetMailFrom()/ Exprcted %s but got %s", newConf.MailFrom, newConf.GetMailFrom())
	}
}
