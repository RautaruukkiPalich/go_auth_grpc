package kafka

import (
	"encoding/json"
	"log/slog"
	"context"
	"github.com/segmentio/kafka-go"
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
	broker *kafka.Writer
	log    *slog.Logger
}

type Brokerer interface {
	AddToQueue(msg KafkaMessage)
	Stop()
}

func New(log *slog.Logger) *Broker {

	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Balancer: &kafka.LeastBytes{},
	}

	log.Info("start broker")

	return &Broker{
		broker: w,
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

	err = b.broker.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: msg.Topic,
			Value: []byte(data),
		})
	if err != nil {
		log.Error("error while sending message", slerr.Err(err))
	}
}

func (b *Broker) Stop() {
	const op = "app.kafka.app.Close"
	log := b.log.With(slog.String("op", op))

	log.Info("close broker")

	defer b.broker.Close()
}
