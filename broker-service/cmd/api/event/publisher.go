package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	connection *amqp.Connection
}

func (p *Publisher) setup() error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (p *Publisher) Push(event string, severity string) error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	// publish method
	err = channel.Publish(
		"topics",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventPublisher(conn *amqp.Connection) (Publisher, error) {
	publisher := Publisher{
		connection: conn,
	}

	err := publisher.setup()
	if err != nil {
		return Publisher{}, err
	}

	return publisher, nil
}
