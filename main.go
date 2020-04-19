package main

import (
  "encoding/json"
  "fmt"
  "github.com/go-redis/redis/v7"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "time"
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

func main() {
  if os.Getenv("REDIS_ADDRESS") == "" ||
	  os.Getenv("ERPLY_USERNAME") == "" ||
	  os.Getenv("ERPLY_PASSWORD") == "" ||
	  os.Getenv("ADDRESS") == "" ||
	  os.Getenv("ERPLY_CLIENT") == "" {
	panic("Missing variables!")
  }

  e := echo.New()
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  link := "https://" + os.Getenv("ERPLY_CLIENT") + ".erply.com/api/"
  redisClient := redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDRESS"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       0, // use default DB
  })

  e.POST("/", func(c echo.Context) error {
	sessionKey, err := redisClient.Get("sessionKey").Result()
	if err != nil {
	  updateSessionKey(link, *redisClient, c)
	  sessionKey, _ = redisClient.Get("sessionKey").Result()
	}
	params, _ := c.FormParams()
	if params.Get("request") == "" {
	  return c.String(422, "Unprocessable Entity: required request paramter is missing!")
	}
	params.Set("clientCode", os.Getenv("ERPLY_CLIENT"))
	params.Set("sessionKey", sessionKey)
	resp, err := http.PostForm(link, params)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	if err != nil {
	  c.Logger().Error(err)
	}
	return c.String(200, bodyString)
  })
  fmt.Println("Build: " + version + " " + buildTime)
  e.Logger.Fatal(e.Start(os.Getenv("ADDRESS")))
}

func updateSessionKey(link string, redisClient redis.Client, c echo.Context) {
  params := url.Values{}
  params.Set("clientCode", os.Getenv("ERPLY_CLIENT"))
  params.Set("username", os.Getenv("ERPLY_USERNAME"))
  params.Set("password", os.Getenv("ERPLY_PASSWORD"))
  params.Set("request", "verifyUser")

  resp, err := http.PostForm(link, params)
  if err != nil {
	c.Logger().Fatal(err)
  }

  bodyBytes, _ := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()
  responseObject := Response{}
  err = json.Unmarshal(bodyBytes, &responseObject)

  if responseObject.Status.ResponseStatus != "ok" || err != nil {
	c.Logger().Fatal(responseObject.Status)
  }

  err = redisClient.Set("sessionKey",
	responseObject.Records[0].SessionKey,
	time.Duration(responseObject.Records[0].SessionLength-1)*time.Second).Err() // decreasing by one second to be cautious
  if err != nil {
	c.Logger().Fatal(err)
  }
}
