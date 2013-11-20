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

func MakeClientWithBaseCreditentials(accessKey string,
  secretKey string, regionName string) *SqsClient {
  return MakeClient(aws.MakeCredentials(accessKey, secretKey), regionName)
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

func (c *SqsClient) makePostRequest(uri string, action string) string {
  return c.makePostRequestWithParams(uri, action, make(map[string]string))
}

func (c *SqsClient) makePostRequestWithParams(uri string, action string, params map[string]string) string {
  params["Action"] = action
  paramValues := url.Values{}
  for k, v := range params {
    paramValues[k] = []string{v}
  }
  request, _ := aws.NewHTTPRequest("POST", uri, paramValues.Encode())
  aws.SignRequest(c, request)
  return makeRequest(request)
}

func (c *SqsClient) makeGetRequest(uri string, action string) string {
  return c.makeGetRequestWithParams(uri, action, make(map[string]string))
}

func (c *SqsClient) makeGetRequestWithParams(uri string, action string, params map[string]string) string {
  params["Action"] = action
  request := generateGetRequest(uri, params)
  aws.SignRequest(c, request)
  response := makeRequest(request)
  return response
}

func (c *SqsClient) ListQueues() string {
  // return c.makeGetRequest("http://" + c.EndPoint(), "ListQueues")
  return c.makePostRequest("http://" + c.EndPoint(), "ListQueues")
}

func (sqs *SqsClient) ListQueuesWithPrefix(prefix string) {

}
