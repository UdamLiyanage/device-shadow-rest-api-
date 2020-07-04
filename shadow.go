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
	if c.QueryParam("shadowName") == "" {
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
	return c.JSON(200, nil)
}
