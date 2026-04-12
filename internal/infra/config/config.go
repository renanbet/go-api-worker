package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr           string
	MongoURI           string
	MongoDB            string
	KafkaBrokers       []string
	KafkaTopic         string
	KafkaGroupID       string
	KafkaNumPartitions int
}

func Load() (Config, error) {
	httpAddr := getenv("HTTP_ADDR", "0.0.0.0:8080")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := getenv("MONGO_DB", "orders_db")
	kafkaBrokersRaw := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := getenv("KAFKA_TOPIC", "order_events")
	kafkaGroupID := getenv("KAFKA_GROUP_ID", "orders-worker")
	kafkaNumPartitions, err := strconv.Atoi(getenv("KAFKA_NUM_PARTITIONS", "3"))
	if err != nil {
		return Config{}, fmt.Errorf("KAFKA_NUM_PARTITIONS is not a number: %w", err)
	}

	var brokers []string
	for _, b := range strings.Split(kafkaBrokersRaw, ",") {
		b = strings.TrimSpace(b)
		if b != "" {
			brokers = append(brokers, b)
		}
	}

	if mongoURI == "" {
		return Config{}, fmt.Errorf("MONGO_URI is required")
	}
	if len(brokers) == 0 {
		return Config{}, fmt.Errorf("KAFKA_BROKERS is required")
	}

	return Config{
		HTTPAddr:           httpAddr,
		MongoURI:           mongoURI,
		MongoDB:            mongoDB,
		KafkaBrokers:       brokers,
		KafkaTopic:         kafkaTopic,
		KafkaGroupID:       kafkaGroupID,
		KafkaNumPartitions: kafkaNumPartitions,
	}, nil
}

type EmailWorkerConfig struct {
	KafkaBrokers  []string
	KafkaTopic    string
	RabbitMQURL   string
	RabbitMQQueue string
}

func LoadEmailWorker() (EmailWorkerConfig, error) {
	kafkaBrokersRaw := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := getenv("KAFKA_TOPIC", "order_events")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	rabbitMQQueue := getenv("RABBITMQ_QUEUE", "email_queue")

	var brokers []string
	for _, b := range strings.Split(kafkaBrokersRaw, ",") {
		b = strings.TrimSpace(b)
		if b != "" {
			brokers = append(brokers, b)
		}
	}

	if len(brokers) == 0 {
		return EmailWorkerConfig{}, fmt.Errorf("KAFKA_BROKERS is required")
	}
	if rabbitMQURL == "" {
		return EmailWorkerConfig{}, fmt.Errorf("RABBITMQ_URL is required")
	}

	return EmailWorkerConfig{
		KafkaBrokers:  brokers,
		KafkaTopic:    kafkaTopic,
		RabbitMQURL:   rabbitMQURL,
		RabbitMQQueue: rabbitMQQueue,
	}, nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
