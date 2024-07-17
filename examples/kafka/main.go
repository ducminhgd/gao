package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/ducminhgd/gao/kafka"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{"35.240.198.124:9094", "35.240.198.124:9095", "35.240.198.124:9096"} // Update with your Kafka broker addresses
	// topic := "gbadmin.topic_subtopic_tables"                                                 // Update with the topic you want to monitor
	group := "GROUP_SINK_KTABLE_MCS_USER_W_GROUP"

	kafkaClient := kafka.NewKafkaClient(brokers, config)

	lagMsg, err := kafkaClient.GetLagMessages(group)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%v\n", lagMsg)
	}

	s, err := kafkaClient.GetConsumerGroupState(group)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(s)
	}

}
