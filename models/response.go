package models
// Erply response struct
type Response struct {
  Records []record `json:"records"`
  Status  status   `json:"status"`
}
// Status part in the response
type status struct {
  ResponseStatus string `json:"responseStatus"`
  ErrorCode      int    `json:"errorCode"`
}
// record struct in response since response can have multiple records
type record struct {
  SessionKey    string `json:"sessionKey"`
  SessionLength int    `json:"sessionLength"`
}