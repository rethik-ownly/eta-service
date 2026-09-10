package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nutanalabs/eta-service/internal/config"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type Repository interface {
	SendMessage(topic string, message []byte) error
	SendKeyedMessage(topic string, key string, message []byte) error
}

type repositoryImpl struct {
	config *config.Config
	producer *kafka.Producer
}

func NewKafkaRepository(config *config.Config, kafkaClient ProducerClient) Repository {
	return &repositoryImpl{
		config: config,
		producer: kafkaClient.GetProducer(),
	}
}

func (r *repositoryImpl) SendMessage(topic string, message []byte) error {
	err := r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, nil)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "KAFKA_SEND_MESSAGE_FAILED",
			Message: fmt.Sprintf("Failed to send message to Kafka server: %v", err),
			Data:    map[string]string{"topic": topic},
		})
		return err
	}

	return nil
}

func (r *repositoryImpl) SendKeyedMessage(topic string, key string, message []byte) error {
	err := r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
		Key: []byte(key),
	}, nil)

	if err != nil {
		logger.Error(logger.Format{
			Event:   "KAFKA_SEND_KEYED_MESSAGE_FAILED",
			Message: fmt.Sprintf("Failed to send keyed message to Kafka server: %v", err),
			Data:    map[string]string{"topic": topic, "key": key},
		})
		return err
	}
	return nil
}