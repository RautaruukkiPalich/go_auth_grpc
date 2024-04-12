package kafka

import (
	"encoding/json"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rautaruukkipalich/go_auth_grpc/internal/lib/slerr"
)

type KafkaMessage struct {
	Topic   string  `json:"topic"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Email   string `json:"email"`
	Header string	`json:"header"`
	Message string `json:"message"`
}

type Broker struct {
	broker *kafka.Producer
	log    *slog.Logger
}

type Brokerer interface {
	AddToQueue(msg KafkaMessage)
	Close()
}

func New(log *slog.Logger) *Broker {

	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": "localhost:29092",
			"client.id":         "kafka-producer",
			"acks":              "all",
		},
	)
	if err != nil {
		panic(err)
	}

	log.Info("start broker")

	return &Broker{
		broker: p,
		log:    log,
	}
}

func (b *Broker) AddToQueue(msg KafkaMessage) {
	const op = "app.kafka.app.AddToQueue"
	log := b.log.With(slog.String("op", op))

	data, err := json.Marshal(msg)
	if err != nil {
		log.Error("error json data", slerr.Err(err))
	}

	err = b.broker.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &msg.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(data),
	}, nil)
	if err != nil {
		log.Error("error while sending message", slerr.Err(err))
	}
}

func (b *Broker) Close() {
	const op = "app.kafka.app.Close"
	log := b.log.With(slog.String("op", op))

	log.Info("close broker")

	b.broker.Close()
}
