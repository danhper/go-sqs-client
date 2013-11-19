package aws

import (
  "bytes"
  "net/http"
)

type HTTPRequest struct {
  *http.Request
  payload string
}

func NewHTTPRequest(method string, url string, body string) (*HTTPRequest, error) {
  bodyReader := bytes.NewBufferString(body)
  baseRequest, err := http.NewRequest(method, url, bodyReader)
  req := &HTTPRequest{baseRequest, body }
  return req, err
}
