package main

import (
	"log"

	"github.com/spf13/viper"
	"micze.io/mt/pkg/ad"
	"micze.io/mt/pkg/api"
	"micze.io/mt/pkg/rabbitmq"
)

func LoadConfigs() (pubsubConfig rabbitmq.PubSubConfig, err error) {
	viper.SetDefault("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("PUB_QUEUE_NAME", "unprepared_q")
	viper.SetDefault("PUB_EXCHANGE_NAME", "unprepared_ex")
	viper.SetDefault("BIND_ADDRESS", ":8080")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&pubsubConfig)
	return pubsubConfig, err
}

func main() {
	log.Println("Starting beta...")
	pubsubConfig, err := LoadConfigs()
	if err != nil {
		log.Fatal("error loading configs: ", err.Error())
	}

	publisher, err := rabbitmq.NewPublisherSubscriber(pubsubConfig)
	if err != nil {
		log.Fatal("error creating store ", err.Error())
	}

	// TODO
	ad := ad.NewAD("localhost", "DC=corp,DC=local", 389)

	server := api.NewServer(publisher, ad)
	address := viper.GetString("BIND_ADDRESS")

	log.Println("Listening for messages...")
	if err = server.Start(address); err != nil {
		log.Fatal("error starting server ", err.Error())
	}

}
