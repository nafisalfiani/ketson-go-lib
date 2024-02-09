package broker

import (
	"context"

	"github.com/nafisalfiani/ketson-go-lib/header"
	"github.com/nafisalfiani/ketson-go-lib/log"
	"github.com/nafisalfiani/ketson-go-lib/parser"
	"github.com/streadway/amqp"
)

type Config struct {
	Url string `env:"URL"`
}

type Interface interface {
	Close() error
	Channel() (*amqp.Channel, error)
	PublishMessage(topic string, obj any) error
}

type broker struct {
	server *amqp.Connection
	log    log.Interface
	json   parser.JSONInterface
}

func Init(cfg Config, log log.Interface, jsonP parser.JSONInterface) (Interface, error) {
	log.Info(context.Background(), "connecting to rabbitmq broker...")

	server, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	return &broker{
		server: server,
		log:    log,
		json:   jsonP,
	}, nil
}

func (b *broker) Close() error {
	return b.server.Close()
}

func (b *broker) Channel() (*amqp.Channel, error) {
	return b.server.Channel()
}

func (b *broker) PublishMessage(topic string, obj any) error {
	ch, err := b.server.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(topic, false, false, false, false, nil)
	if err != nil {
		return err
	}

	objJson, err := b.json.Marshal(obj)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: header.ContentTypeJSON,
		Body:        objJson,
	}

	return ch.Publish("", queue.Name, false, false, message)
}
