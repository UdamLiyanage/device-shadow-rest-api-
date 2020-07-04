package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var (
	shadow     *mongo.Collection
	mqttClient mqtt.Client
)

func init() {
	client := connect()
	shadow = databaseCollection(client, "Shadow", "shadow")
}

func onMessageReceived(_ mqtt.Client, msg mqtt.Message) {
	println(string(msg.Payload()))
}

func setClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("BROKER_URL"))
	opts.SetUsername(os.Getenv("BROKER_USERNAME"))
	opts.SetPassword(os.Getenv("BROKER_PASSWORD"))
	opts.SetDefaultPublishHandler(onMessageReceived)
	return opts
}

func main() {
	opts := setClientOptions()

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Opts Init")
	}
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Println("Connected to MQTT Broker")
	}
}
