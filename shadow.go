package main

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

func getShadow(c echo.Context) error {
	var (
		crud Operations = Configuration{
			Collection: shadowCollection,
		}
		shadow map[string]interface{}
	)
	objID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		return c.JSON(500, err)
	}
	err = crud.Index(bson.M{
		"_id": objID,
	}).Decode(&shadow)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(500, nil)
	}
	return c.JSON(200, shadow)
}

func getDeviceShadows(c echo.Context) error {
	var (
		crud Operations = Configuration{
			Collection: shadowCollection,
		}
		shadows []map[string]interface{}
	)
	response, err := crud.Read(bson.M{
		"device": c.Param("urn"),
	})
	if err != nil {
		if response.Err() == mongo.ErrNoDocuments {
			return c.JSON(404, nil)
		}
		if err == mongo.ErrNoDocuments {
			return c.JSON(404, nil)
		}
	}
	err = response.All(context.Background(), &shadows)
	if err != nil {
		log.Error(err.Error())
	}
	return c.JSON(200, shadows)
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
