package main

import (
  "encoding/json"
  "flag"
  "github.com/go-redis/redis/v7"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "time"
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

type Response struct {
  Records []record `json:"records"`
  Status  status   `json:"status"`
}

type status struct {
  ResponseStatus string `json:"responseStatus"`
  ErrorCode      int    `json:"errorCode"`
}
type record struct {
  SessionKey    string `json:"sessionKey"`
  SessionLength int    `json:"sessionLength"`
}

func flagParse() {
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
}

func MainHandler(ctx echo.Context) error {
  link := "https://" + *erplyClient + ".erply.com/api/"
  redisClient := redis.NewClient(&redis.Options{
	Addr:     *redisAddress,
	Password: *redisPassword,
	DB:       0, // use default DB
  })

  sessionKey, err := redisClient.Get("sessionKey").Result()
  if err != nil {
	updateSessionKey(link, *redisClient)
	sessionKey, err = redisClient.Get("sessionKey").Result()
	if err != nil {
	  log.Fatal(err)
	}
  }

  params, _ := ctx.FormParams()
  if params.Get("request") == "" {
	return ctx.String(422, "Unprocessable Entity: required request paramter is missing!")
  }
  params.Set("clientCode", *erplyClient)
  params.Set("sessionKey", sessionKey)
  resp, err := http.PostForm(link, params)
  if err != nil {
	log.Error(err)
  }

  defer resp.Body.Close()
  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
	log.Error(err)
  }
  bodyString := string(bodyBytes)

  return ctx.String(200, bodyString)
}

func main() {
  log.SetFormatter(&log.TextFormatter{
	FullTimestamp: true,
  })
  flagParse()

  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.POST("/", MainHandler)
  log.Info("Build: " + version + " " + buildTime)
  log.Fatal(e.Start(os.Getenv("ADDRESS")))
}

func updateSessionKey(link string, redisClient redis.Client) {
  params := url.Values{}
  params.Set("clientCode", os.Getenv("ERPLY_CLIENT"))
  params.Set("username", os.Getenv("ERPLY_USERNAME"))
  params.Set("password", os.Getenv("ERPLY_PASSWORD"))
  params.Set("request", "verifyUser")

  resp, err := http.PostForm(link, params)
  if err != nil {
	log.Fatal(err)
  }

  bodyBytes, _ := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()
  responseObject := Response{}
  err = json.Unmarshal(bodyBytes, &responseObject)

  if responseObject.Status.ResponseStatus != "ok" || err != nil {
	log.Fatal(responseObject.Status)
  }

  err = redisClient.Set("sessionKey",
	responseObject.Records[0].SessionKey,
	time.Duration(responseObject.Records[0].SessionLength-1)*time.Second).Err() // decreasing by one second to be cautious
  if err != nil {
	log.Fatal(err)
  }
}
