package aws

import (
  "fmt"
  "strings"
)

const endpoint = "sqs.%s.amazonaws.com"

const (
  Virginia   = "us-east-1"
  California = "us-west-1"
  Oregon     = "us-west-2"
  Ireland    = "eu-west-1"
  Singapore  = "ap-southeast-1"
  Sydney     = "ap-southeast-2"
  Tokyo      = "ap-northeast-1"
  SaoPaulo   = "sa-east-1"
)

type Client interface {
  RegionName() string
  ServiceName() string
  Credentials() *Credentials
}

type BaseClient struct {
  regionName string
  credentials *Credentials
}

func MakeBaseClient(regionName string, credentials *Credentials) BaseClient {
  return BaseClient {
    regionName  : regionName,
    credentials : credentials,
  }
}

func (b *BaseClient) RegionName() string {
  return b.regionName
}

func (b *BaseClient) Credentials() *Credentials {
  return b.credentials
}

func (b *BaseClient) GetEndPoint() string {
  return fmt.Sprintf(endpoint, b.RegionName())
}

func SignRequest(c Client, request *HTTPRequest) *HTTPRequest {
  if strings.ToLower(request.Method) == "get" {
    return signGetRequest(c, request)
  } else {
    return  signPostRequest(c, request)
  }
}

func signGetRequest(c Client, request *HTTPRequest) *HTTPRequest {
  setSignedHeaders(c, request)
  fmt.Println(request.Header.Get("X-Amz-Credential"))
  return request
}

func signPostRequest(c Client, request *HTTPRequest) *HTTPRequest {
  return request
}
