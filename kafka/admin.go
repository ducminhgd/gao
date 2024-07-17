package kafka

import (
	"fmt"

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

func (kc *KafkaClient) GetLagMessages(group string) (map[string]map[int32]int64, error) {
	topics, err := kc.Topics()
	if err != nil {
		return nil, err
	}

	lagMap := make(map[string]map[int32]int64)

	for _, topic := range topics {
		partitions, err := kc.Partitions(topic)
		if err != nil {
			continue
		}

		partitionMap := make(map[int32]int64)

		for _, partition := range partitions {
			offsetManager, err := sarama.NewOffsetManagerFromClient(group, kc)
			if err != nil {
				continue
			}

			partitionOffsetManager, err := offsetManager.ManagePartition(topic, partition)
			if err != nil {
				continue
			}

			newestOffset, err := kc.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				continue
			}

			currentOffset, _ := partitionOffsetManager.NextOffset()
			lag := newestOffset - currentOffset

			if lag > 0 {
				partitionMap[partition] = lag
			}
		}

		if len(partitionMap) > 0 {
			lagMap[topic] = partitionMap
		}
	}

	return lagMap, nil
}

func (kc *KafkaClient) GetConsumerGroupState(name string) (string, error) {
	admin, err := sarama.NewClusterAdminFromClient(kc.Client)
	if err != nil {
		fmt.Printf("Failed to create cluster admin: %v", err)
		return "Unknown", err
	}
	defer admin.Close()

	group, err := admin.DescribeConsumerGroups([]string{name})
	if err != nil {
		fmt.Printf("Failed to describe consumer group: %v", err)
		return "Unknown", err
	}
	if len(group) == 0 {
		fmt.Printf("Consumer group %s not found", name)
		return "Unknown", err
	}
	return group[0].State, nil
}
