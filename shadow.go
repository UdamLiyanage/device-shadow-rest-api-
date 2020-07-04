package main

import "github.com/labstack/echo/v4"

func getShadow(c echo.Context) error {
	return c.JSON(200, nil)
}

func getDeviceShadows(c echo.Context) error {
	return c.JSON(200, nil)
}

func updateShadow(c echo.Context) error {
	return c.JSON(200, nil)
}
