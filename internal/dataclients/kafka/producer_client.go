package kafka

import (
	"errors"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nutanalabs/eta-service/internal/config"
	logger "github.com/nutanalabs/rapido-logger-go"
)

type ProducerClient interface {
	GetProducer() *kafka.Producer
	CheckHealth() error
	Close()
}

type producerClientImpl struct {
	producer *kafka.Producer
}

func NewKafkaProducerClient(config *config.Config) (ProducerClient, error) {
	producer, err := newKafkaProducer(config)
	return &producerClientImpl{
		producer: producer,
	}, err
}

func (p *producerClientImpl) GetProducer() *kafka.Producer {
	return p.producer
}

func (p *producerClientImpl) CheckHealth() error {
	var topicName = ""
	metaData, err := p.producer.GetMetadata(&topicName, false, 1000)
	if err != nil {
		logger.Error(logger.Format{
			Event:   "KAFKA_CHECK_BROKERS_FAILED",
			Message: fmt.Sprintf("failed to get metadata for topic, with error as: %v", err),
		})
		return errors.New("KafkaHealthCheckError")
	}
	brokersMetadata := metaData.Brokers
	if brokersMetadata == nil || len(brokersMetadata) <= 0 {
		logger.Error(logger.Format{
			Event:   "KAFKA_CHECK_BROKERS_FAILED",
			Message: fmt.Sprintf("failed to read brokers data from Kafka, with error as: %v", brokersMetadata),
		})
		return errors.New("kafkaProducerHealthCheckError")
	}
	logger.Debug(logger.Format{
		Event: "KAFKA_HEALTH_CHECK",
		Message: "kafka health check is successful",
	})
	return nil
}

func (p *producerClientImpl) Close() {
	if p.producer != nil {
		p.producer.Close()
		logger.Info(logger.Format{
			Event:   "KAFKA_CLOSE",
			Message: "Closed Kafka Producer & Consumer Connection",
		})
	}
}

func newKafkaProducer(config *config.Config) (*kafka.Producer, error) {
	kafkaConfig := config.GetKafkaConfig()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaConfig.BootstrapServers})
	if err != nil {
		logger.Error(logger.Format{
			Event:   "KAFKA_NEW_PRODUCER_FAILED",
			Message: fmt.Sprintf("Failed to create new producer on kafka: %v", err),
			Data:    map[string]string{"config": kafkaConfig.BootstrapServers},
		})
		return nil, err
	}

	logger.Info(logger.Format{Message: "Successfully connected to kafka as a producer"})

	// Delivery report handler for produced messages ( Tells whether the the produced event reached the partition or not)
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()

	return producer, nil
}