package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var kafkaProducer *kafka.Producer

// Initialize Kafka producer
func InitKafkaProducer(broker string) {
	var err error
	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	log.Println("Kafka producer initialized!")
}

// Send message to Kafka
func SendToKafka(topic string, message string) {
	if kafkaProducer == nil {
		log.Println("Kafka producer is not initialized")
		return
	}

	// Send the message
	err := kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	} else {
		log.Printf("Message sent to Kafka topic %s: %s", topic, message)
	}
}
