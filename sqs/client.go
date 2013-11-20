package sqs

import (
  "fmt"
  "io/ioutil"
  "net/http"
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

func generateGetRequest(url string, params map[string]string) *aws.HTTPRequest {
  request, _ := aws.NewHTTPRequest("GET", url, "")
  query := request.URL.Query()
  for k, v := range params {
    query.Set(k, v)
  }
  request.URL.RawQuery = query.Encode()
  return request
}

func makeGetRequest(c *SqsClient, url string, action string) string {
  return makeGetRequestWithParams(c, url, action, make(map[string]string))
}

func makeGetRequestWithParams(c *SqsClient, url string, action string, params map[string]string) string {
  params["Action"] = action
  request := generateGetRequest(url, params)
  aws.SignRequest(c, request)
  response := makeRequest(request)
  return response
}

func (c *SqsClient) ListQueues() string {
  return makeGetRequest(c, "http://" + c.EndPoint(), "ListQueues")
}

func (sqs *SqsClient) ListQueuesWithPrefix(prefix string) {

}
