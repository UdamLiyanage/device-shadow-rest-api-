package main

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getShadow(c echo.Context) error {
	var (
		crud Operations = Configuration{
			Collection: shadowCollection,
		}
		shadow Shadow
	)
	result := crud.Index(bson.M{
		"device": "",
	})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.JSON(404, nil)
		}
		panic(result.Err())
	}
	err := result.Decode(&shadow)
	if err != nil {
		panic(err)
	}
	return c.JSON(200, shadow)
}

func getDeviceShadows(c echo.Context) error {
	return c.JSON(200, nil)
}

func updateShadow(c echo.Context) error {
	return c.JSON(200, nil)
}
