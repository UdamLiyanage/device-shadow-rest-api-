package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func connectToBroker() {
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

func setupRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} method=${method}, uri=${uri}, status=${status} path=${path} latency=${latency_human}\n",
	}))

	shadowGroup := e.Group("/api/v1/shadows")

	shadowGroup.GET("/:urn/shadow", getShadow)
	shadowGroup.GET("/device/:urn/shadows", getDeviceShadows)

	shadowGroup.POST("/:urn/shadow/update", updateShadow)

	return e
}

func main() {
	connectToBroker()
	r := setupRouter()
	r.Logger.Fatal(r.Start(":8080"))
}
