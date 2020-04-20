package config

import (
  "erply-middleware/models"
  "flag"
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

func FlagParse(buildTime string, version string) models.Config {

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
  log.Info("Build: " + version + " " + buildTime)
  return models.Config{
	Client: *erplyClient,
	Username: *username,
	Password: *password,
	Link: "https://" + *erplyClient + ".erply.com/api/",
	RedisAddress: *redisAddress,
	RedisPassword: *redisPassword,
	Address: *address,
  }
}