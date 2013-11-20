package aws

import (
  "strings"
)

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
  RegionName()  string
  ServiceName() string
  Credentials() *Credentials
  EndPoint()    string
  Protocol()    string
}

type BaseClient struct {
  regionName string
  credentials *Credentials
  protocol   string
}

func MakeBaseClient(regionName string, credentials *Credentials) BaseClient {
  return BaseClient {
    regionName  : regionName,
    credentials : credentials,
    protocol    : "http",
  }
}

func GenerateURL(c Client) string {
  return c.Protocol() + "://" + c.EndPoint()
}

func MakeBaseClientWithProtocol(regionName string, credentials *Credentials, protocol string) BaseClient {
  return BaseClient {
    regionName  : regionName,
    credentials : credentials,
    protocol    : protocol,
  }
}

func (b *BaseClient) Protocol() string {
  return b.protocol
}

func (b *BaseClient) RegionName() string {
  return b.regionName
}

func (b *BaseClient) Credentials() *Credentials {
  return b.credentials
}

func SignRequest(c Client, request *HTTPRequest) {
  request.Header.Set("host", c.EndPoint())
  if strings.ToLower(request.Method) == "get" {
    signGetRequest(c, request)
  } else {
    signPostRequest(c, request)
  }
}

func signGetRequest(c Client, request *HTTPRequest) {
  updateQueryWithAuth(c, request)
  query := request.URL.Query()
  query.Set("X-Amz-Signature", getSignature(c, request))
  request.URL.RawQuery = query.Encode()
}

func signPostRequest(c Client, request *HTTPRequest) {
  updateHeadersForAuth(c, request)
}
