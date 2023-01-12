package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SwitchEndpointParams struct {
	Mac       string `json:"mac,omitempty"`
	IpAddr    string `json:"ipaddr,omitempty"`
	SysId     string `json:"sysid,omitempty"`
	Username  string `json:"username,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Id        int64  `json:"id,omitempty"`
}

type PubSubConfig struct {
	PubExchangeName string `mapstructure:"PUB_EXCHANGE_NAME"`
	PubQueueName    string `mapstructure:"PUB_QUEUE_NAME"`
	SubExchangeName string `mapstructure:"SUB_EXCHANGE_NAME"` // not required TO REMOVE after verification
	SubQueueName    string `mapstructure:"SUB_QUEUE_NAME"`
	URI             string `mapstructure:"RABBITMQ_URI"`
}

type PublisherSubscriber struct {
	channel *amqp.Channel
	config  PubSubConfig
}

func (ps *PublisherSubscriber) Publish(params SwitchEndpointParams) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = ps.channel.PublishWithContext(ctx,
		ps.config.PubExchangeName, // exchange
		"",                        // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}

func (ps *PublisherSubscriber) Subscribe() (chan SwitchEndpointParams, error) {
	ch := make(chan SwitchEndpointParams)

	_, err := ps.channel.QueueDeclare(
		ps.config.SubQueueName, // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ps.channel.Consume(
		ps.config.SubQueueName, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return nil, err
	}

	go func() {
		for data := range msgs {
			var msg SwitchEndpointParams
			json.Unmarshal(data.Body, &msg)
			ch <- msg
		}
	}()

	return ch, nil
}

func (ps *PublisherSubscriber) Close() {
	ps.channel.Close()
}

func NewPublisherSubscriber(config PubSubConfig) (*PublisherSubscriber, error) {
	pubsub := &PublisherSubscriber{}

	conn, err := amqp.Dial(config.URI)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		config.PubExchangeName, // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(
		config.PubQueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		q.Name,
		"",
		config.PubExchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	pubsub.config = config
	pubsub.channel = channel

	return pubsub, nil
}
