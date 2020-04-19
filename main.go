package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Response struct {
	Records []record `json:"records"`
	Status  status   `json:"status"`
}

type status struct{
	ResponseStatus string `json:"responseStatus"`
	ErrorCode int `json:"errorCode"`
}
type record struct{
	SessionKey string `json:"sessionKey"`
}

func main() {
	if os.Getenv("REDIS_ADDRESS") == "" ||
		os.Getenv("ERPLY_USERNAME") == "" ||
		os.Getenv("ERPLY_PASSWORD") == "" ||
		os.Getenv("ERPLY_CLIENT") == "" {
		panic("Missing variables!")
	}

	e := echo.New()
	e.Use(middleware.Logger())

	link := "https://" + os.Getenv("ERPLY_CLIENT") + ".erply.com/api/"
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,  // use default DB
	})

	updateSessionKey(link, *redisClient)

	e.POST("/", func(c echo.Context) error {
		sessionKey, err := redisClient.Get("sessionKey").Result()
		if err != nil {
			updateSessionKey(link, *redisClient)
			sessionKey, _ = redisClient.Get("sessionKey").Result()
		}
		params, _ := c.FormParams()
		params.Set("clientCode", os.Getenv("ERPLY_CLIENT"))
		params.Set("sessionKey", sessionKey)
		resp, err := http.PostForm(link,  params)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		if err != nil {
			panic(err)
		}

		return c.String(200,bodyString)
	})
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

func updateSessionKey(link string, redisClient redis.Client) {
	params := url.Values{}
	params.Set("clientCode", os.Getenv("ERPLY_CLIENT"))
	params.Set("username", os.Getenv("ERPLY_USERNAME"))
	params.Set("password", os.Getenv("ERPLY_PASSWORD"))
	params.Set("request", "verifyUser")

	resp, err := http.PostForm(link,  params)
	if err != nil {
		panic(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	responseObject := Response{}
	err = json.Unmarshal(bodyBytes, &responseObject)

	if responseObject.Status.ResponseStatus != "ok" || err != nil {
		panic(responseObject.Status)
	}

	redisClient.Set("sessionKey",responseObject.Records[0].SessionKey, 3599 * time.Second )
}