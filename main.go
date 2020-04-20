package main

import (
  "erply-middleware/handlers"
  "erply-middleware/models"
  "flag"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  log "github.com/sirupsen/logrus"
  "os"
)


var (
  showVersion = flag.Bool("version", false, "Print version")
  erplyClient  = flag.String("client", os.Getenv("ERPLY_CLIENT"), "Set erply client code")
  redisAddress = flag.String("redisAddress", os.Getenv("REDIS_ADDRESS"), "Redis server address")
  redisPassword = flag.String("redisPassword", os.Getenv("REDIS_PASSWORD"), "Redis server password")
  username     = flag.String("username", os.Getenv("ERPLY_USERNAME"), "Erply username")
  password     = flag.String("password", os.Getenv("ERPLY_PASSWORD"), "Erply password")
  address      = flag.String("address", os.Getenv("ADDRESS"), "Server address and port written like this 127.0.0.1:1232")
)

var (
  buildTime string
  version   string
)

func flagParse() models.Config {
  flag.Parse()
  if *showVersion {
	log.Info("Build:", version, buildTime)
	os.Exit(0)
  }

  if *redisAddress == "" ||
	  *username == "" ||
	  *password == "" ||
	  *address == "" ||
	  *erplyClient == "" {
    	log.Fatal("Some parameters are missing!")
  }
  return models.Config{
    Client: *erplyClient,
    Username: *username,
    Password: *password,
    Link: "https://" + *erplyClient + ".erply.com/api/",
    RedisAddress: *redisAddress,
    RedisPassword: *redisPassword,
  }
}

func main() {
  log.SetFormatter(&log.TextFormatter{
	FullTimestamp: true,
  })
  config := flagParse()
  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.POST("/", func(ctx echo.Context) error {
	return handlers.MainHandler(ctx, config)
  })
  log.Info("Build: " + version + " " + buildTime)
  log.Fatal(e.Start(*address))
}

