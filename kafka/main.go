package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{"35.240.198.124:9094", "35.240.198.124:9095", "35.240.198.124:9096"} // Update with your Kafka broker addresses
	topic := "gbadmin.topic_subtopic_tables"                                                 // Update with the topic you want to monitor
	group := "terry.gbadmin_cdc.topic_tables-earliest"                                       // Update with your consumer group

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	totalLag := int64(0)

	for _, partition := range partitions {
		offsetManager, err := sarama.NewOffsetManagerFromClient(group, client)
		if err != nil {
			log.Fatal(err)
		}
		defer offsetManager.Close()

		partitionOffsetManager, err := offsetManager.ManagePartition(topic, partition)
		if err != nil {
			log.Fatal(err)
		}
		defer partitionOffsetManager.Close()

		// oldestOffset, err := client.GetOffset(topic, partition, sarama.OffsetOldest)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		newestOffset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal(err)
		}

		consumerOffset, _ := partitionOffsetManager.NextOffset()

		lag := newestOffset - consumerOffset
		totalLag += lag

		fmt.Printf("Partition %d Lag: %d\n", partition, lag)
	}

	fmt.Printf("Total Lag: %d\n", totalLag)
}
