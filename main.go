package main

import (
  "erply-middleware/config"
  "erply-middleware/handlers"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  log "github.com/sirupsen/logrus"
)

func main() {
  log.SetFormatter(&log.TextFormatter{
	FullTimestamp: true,
  })
  configuration := config.FlagParse()
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.POST("/", func(ctx echo.Context) error {
	return handlers.MainHandler(ctx, configuration)
  })
  log.Fatal(e.Start(configuration.Address))
}

