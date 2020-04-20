package handlers

import (
  "encoding/json"
  "erply-middleware/models"
  "github.com/go-redis/redis/v7"
  "github.com/labstack/echo/v4"
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "net/http"
  "net/url"
  "time"
)

func MainHandler(ctx echo.Context, config models.Config) error {
  redisClient := redis.NewClient(&redis.Options{
	Addr:     config.RedisAddress,
	Password: config.RedisPassword,
	DB:       0, // use default DB
  })

  sessionKey, err := redisClient.Get("sessionKey").Result()
  if err != nil {
	updateSessionKey(*redisClient, config)
	sessionKey, err = redisClient.Get("sessionKey").Result()
	if err != nil {
	  log.Fatal(err)
	}
  }

  params, _ := ctx.FormParams()
  if params.Get("request") == "" {
	return ctx.String(422, "Unprocessable Entity: required request paramter is missing!")
  }
  params.Set("clientCode", config.Client)
  params.Set("sessionKey", sessionKey)
  resp, err := http.PostForm(config.Link, params)
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


func updateSessionKey(redisClient redis.Client, config models.Config) {
  params := url.Values{}
  params.Set("clientCode", config.Client)
  params.Set("username", config.Username)
  params.Set("password", config.Password)
  params.Set("request", "verifyUser")

  resp, err := http.PostForm(config.Link, params)
  if err != nil {
	log.Fatal(err)
  }

  bodyBytes, _ := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()
  responseObject := models.Response{}
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