package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func getShadow(c echo.Context) error {
	var (
		crud Operations = Configuration{
			Collection: shadowCollection,
		}
		shadow Shadow
	)
	if c.QueryParam("shadowName") != "" {
		result := crud.Index(bson.M{
			"device": c.Param("urn"),
			"name":   c.QueryParam("shadowName"),
		})
		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				return c.JSON(404, nil)
			}
			panic(result.Err())
		}
		err := result.Decode(&shadow)
		if err != nil {
			return c.JSON(500, nil)
		}
		return c.JSON(200, shadow)
	} else {
		response, err := crud.Read(bson.M{
			"device": c.Param("urn"),
		})
		if err != nil {
			return c.JSON(500, nil)
		}
		return c.JSON(200, response)
	}
}

func getDeviceShadows(c echo.Context) error {
	return c.JSON(200, nil)
}

func updateShadow(c echo.Context) error {
	var (
		body  map[string]interface{}
		token mqtt.Token
	)
	for !mqttClient.IsConnected() {
		connectToBroker()
	}
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return c.JSON(500, err)
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return c.JSON(500, err)
	}
	if c.QueryParam("shadowName") == "" {
		token = mqttClient.Publish(
			// /shadows/:urn/shadow/update
			os.Getenv("MQTT_PUBLISH_BASE_TOPIC")+c.Param("urn")+"/shadow/update",
			0,
			false, payload)
		token.WaitTimeout(250)
	} else {
		token = mqttClient.Publish(
			// /shadows/:urn/shadowName/shadow/update
			os.Getenv("MQTT_PUBLISH_BASE_TOPIC")+c.Param("urn")+"/"+c.QueryParam("shadowName")+"/shadow/update",
			0,
			false, payload)
		token.WaitTimeout(250)
	}
	if token.Error() != nil {
		return c.JSON(500, token.Error().Error())
	}
	return c.JSON(200, nil)
}
