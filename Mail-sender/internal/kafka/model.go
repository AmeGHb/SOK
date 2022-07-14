package kafka

type ToKafkaStartFunc struct {
	KafkaAdress []string
	Topic       string
	GroupID     string
}
