package main

import (
	"log"

	"github.com/spf13/viper"
	"micze.io/mt/pkg/etcdstore"
	rmq "micze.io/mt/pkg/rabbitmq"
)

type Config struct {
	RabbitMqURI        string `mapstructure:"RABBITMQ_URI"`
	PrepQueueName      string `mapstructure:"PREP_QUEUE_NAME"`
	PrepExchangeName   string `mapstructure:"PREP_EXCHANGE_NAME"`
	UnprepQueueName    string `mapstructure:"UNPREP_QUEUE_NAME"`
	UnprepExchangeName string `mapstructure:"UNPREP_EXCHANGE_NAME"`
	EtcdURI            string `mapstructure:"ETCD_URI"`
}

func LoadConfig() (config Config, err error) {
	viper.SetDefault("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("PREP_QUEUE_NAME", "prepared_q")
	viper.SetDefault("PREP_EXCHANGE_NAME", "prepared_ex")
	viper.SetDefault("UNPREP_QUEUE_NAME", "unprepared_q")
	viper.SetDefault("UNPREP_EXCHANGE_NAME", "unprepared_ex")
	viper.SetDefault("ETCD_URI", "localhost:2379")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}

func main() {
	log.Println("Starting beta...")
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("error loading config: ", err.Error())
	}

	pubsubConfig := rmq.PubSubConfig{
		URI:             config.RabbitMqURI,
		PubQueueName:    config.PrepQueueName,
		PubExchangeName: config.PrepExchangeName,
		SubExchangeName: config.UnprepExchangeName,
		SubQueueName:    config.UnprepQueueName,
	}

	storeConfig := etcdstore.StoreConfig{
		URI: config.EtcdURI,
	}

	pubsub, err := rmq.NewPublisherSubscriber(pubsubConfig)
	handleError("error creating publishersubscriber:", err)

	store, err := etcdstore.NewStore(storeConfig)
	handleError("error creating publishersubscriber:", err)

	ch, err := pubsub.Subscribe()
	handleError("error subscribing for messages:", err)

	log.Println("Listening for messages...")
	for m := range ch {
		log.Printf("[M%06d] received unprepared message: %v\n", m.Id, m)
		log.Printf("[M%06d] configuring deltas for MAC: %v\n", m.Id, m.Mac)
		err := store.Put(m.Mac)
		if err != nil {
			log.Printf("[M%06d] error while configuring deltas using etcd: %s\n", m.Id, err.Error())
		}
		log.Printf("[M%06d] publishing task to prepared queue: %v\n", m.Id, m)
		if err := pubsub.Publish(m); err != nil {
			log.Printf("[M%06d] error publishing task to prepared queue: %v\n", m.Id, err.Error())
		}
		log.Printf("[M%06d] completed processing message: %v\n", m.Id, m)
	}
}

func handleError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err.Error())
	}
}
