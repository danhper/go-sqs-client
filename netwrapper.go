package aws

import (
  "bytes"
  "net/http"
  "time"
)

type HTTPRequest struct {
  *http.Request
  payload string
  createdTime time.Time
}

func NewHTTPRequestWithTime(method string, url string,
  body string, createdTime time.Time) (*HTTPRequest, error) {
  bodyReader := bytes.NewBufferString(body)
  baseRequest, err := http.NewRequest(method, url, bodyReader)
  req := &HTTPRequest{baseRequest, body, createdTime}
  return req, err
}

func NewHTTPRequest(method string, url string, body string) (*HTTPRequest, error) {
  return NewHTTPRequestWithTime(method, url, body, time.Now().UTC())
}
