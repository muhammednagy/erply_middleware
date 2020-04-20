package handlers

import (
  "erply-middleware/config"
  "github.com/labstack/echo/v4"
  "github.com/stretchr/testify/assert"
  "net/http"
  "net/http/httptest"
  "net/url"
  "strings"
  "testing"
)


func TestMainHandlerWithRequest(t *testing.T) {
  configuration := config.FlagParse()
  e := echo.New()
  form := url.Values{}
  form.Set("request", `getProducts`)
  req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))

  req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

  rec := httptest.NewRecorder()
  c := e.NewContext(req, rec)

  if assert.NoError(t, MainHandler(c, configuration)) {
	assert.Equal(t, http.StatusOK, rec.Code)
  }
}

func TestMainHandlerWithoutRequest(t *testing.T) {
  configuration := config.FlagParse()
  e := echo.New()
  req, _ := http.NewRequest(http.MethodPost, "/", nil)
  rec := httptest.NewRecorder()
  c := e.NewContext(req, rec)

  if assert.NoError(t, MainHandler(c, configuration)) {
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
  }
}
