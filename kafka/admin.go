package kafka

import (
	"github.com/IBM/sarama"
)

type KafkaClient struct {
	sarama.Client
}

func NewKafkaClient(brokers []string, config *sarama.Config) *KafkaClient {
	client, _ := sarama.NewClient(brokers, config)
	return &KafkaClient{
		client,
	}
}
