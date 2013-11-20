package sqs

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "fringe81.com/tuvistavie/aws"
)

const endpoint = "sqs.%s.amazonaws.com"

type SqsClient struct {
  aws.BaseClient
}

func (s *SqsClient) ServiceName() string {
  return "sqs"
}

func MakeClient(credentials *aws.Credentials, regionName string) *SqsClient {
  return &SqsClient {aws.MakeBaseClient(regionName, credentials)}
}

func MakeClientWithProtocol(credentials *aws.Credentials, regionName string, protocol string) *SqsClient {
  return &SqsClient {aws.MakeBaseClientWithProtocol(regionName, credentials, protocol)}
}

func MakeClientAndCreditentials(accessKey string,
  secretKey string, regionName string) *SqsClient {
  return MakeClient(aws.MakeCredentials(accessKey, secretKey), regionName)
}

func MakeClientAndCreditentialsWithProtocol(accessKey string,
  secretKey string, regionName string, protocol string) *SqsClient {
  return MakeClientWithProtocol(aws.MakeCredentials(accessKey, secretKey), regionName, protocol)
}

func (c *SqsClient) EndPoint() string {
  return fmt.Sprintf(endpoint, c.RegionName())
}

func makeRequest(request *aws.HTTPRequest) string {
  client := &http.Client{}
  resp, _ := client.Do(request.Request)
  body, _ := ioutil.ReadAll(resp.Body)
  return fmt.Sprintf("%s", body)
}

func generateGetRequest(uri string, params map[string]string) *aws.HTTPRequest {
  request, _ := aws.NewHTTPRequest("GET", uri, "")
  query := request.URL.Query()
  for k, v := range params {
    query.Set(k, v)
  }
  request.URL.RawQuery = query.Encode()
  return request
}

func (c *SqsClient) makePostRequest(action string) string {
  return c.makePostRequestWithParams(action, make(map[string]string))
}

func (c *SqsClient) makePostRequestWithParams(action string, params map[string]string) string {
  uri := aws.GenerateURL(c)
  params["Action"] = action
  paramValues := url.Values{}
  for k, v := range params {
    paramValues[k] = []string{v}
  }
  request, _ := aws.NewHTTPRequest("POST", uri, paramValues.Encode())
  aws.SignRequest(c, request)
  return makeRequest(request)
}

func (c *SqsClient) makeGetRequest(action string) string {
  return c.makeGetRequestWithParams(action, make(map[string]string))
}

func (c *SqsClient) makeGetRequestWithParams(action string, params map[string]string) string {
  uri := aws.GenerateURL(c)
  params["Action"] = action
  request := generateGetRequest(uri, params)
  aws.SignRequest(c, request)
  response := makeRequest(request)
  return response
}

func (c *SqsClient) ListQueues() string {
  return c.makeGetRequest("ListQueues")
}

func (c *SqsClient) ListQueuesWithPrefix(prefix string) string {
  return c.makeGetRequestWithParams("ListQueues", map[string]string{
    "QueueNamePrefix": prefix,
  })
}
