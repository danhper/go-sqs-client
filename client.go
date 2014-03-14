package sqs

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"

  "github.com/tuvistavie/go-aws-common"
)

const endpoint = "sqs.%s.amazonaws.com"

type Error struct {
  StatusCode int
  Code       string `xml:"Error>Code"`
  Message    string `xml:"Error>Message"`
  RequestId  string
}

type SqsClient struct {
  aws.BaseClient
}

func (s *SqsClient) ServiceName() string {
  return "sqs"
}

func (e *Error) Error() string {
  return fmt.Sprintf("[%d %s] %s: %s", e.StatusCode, e.Code, e.RequestId, e.Message)
}

func MakeClient(credentials *aws.Credentials, regionName string) *SqsClient {
  return &SqsClient{aws.MakeBaseClient(regionName, credentials)}
}

func MakeClientFromEnv() *SqsClient {
  return MakeClient(aws.MakeCredentialsFromEnv(), os.Getenv("AWS_REGION"))
}

func MakeClientWithProtocol(credentials *aws.Credentials, regionName string, protocol string) *SqsClient {
  return &SqsClient{aws.MakeBaseClientWithProtocol(regionName, credentials, protocol)}
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

func makeRequest(request *aws.HTTPRequest, resp interface{}) error {
  var r *http.Response
  var body []byte
  var err error

  client := &http.Client{}

  r, err = client.Do(request.Request)

  if err != nil {
    return err
  }

  body, err = ioutil.ReadAll(r.Body)

  if err != nil {
    return err
  }

  if r.StatusCode != 200 {
    return getError(r, body)
  }

  err = xml.Unmarshal(body, resp)
  return err
}

func getError(r *http.Response, body []byte) error {
  err := &Error{}
  e := xml.Unmarshal(body, err)
  if e != nil {
    return e
  }
  err.StatusCode = r.StatusCode
  return err
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

func (c *SqsClient) makePostRequest(action string, queue *Queue, resp interface{}) error {
  return c.makePostRequestWithParams(action, make(map[string]string), queue, resp)
}

func (c *SqsClient) makePostRequestWithParams(action string, params map[string]string, queue *Queue, resp interface{}) error {
  uri := c.getURL(queue)
  params["Action"] = action
  paramValues := url.Values{}
  for k, v := range params {
    paramValues[k] = []string{v}
  }
  request, _ := aws.NewHTTPRequest("POST", uri, paramValues.Encode())
  aws.SignRequest(c, request)
  return makeRequest(request, resp)
}

func (c *SqsClient) getURL(queue *Queue) string {
  if queue == nil {
    return aws.GenerateURL(c)
  } else {
    return queue.URL
  }
}

func (c *SqsClient) makeGetRequest(action string, queue *Queue, resp interface{}) error {
  return c.makeGetRequestWithParams(action, make(map[string]string), queue, resp)
}

func (c *SqsClient) makeGetRequestWithParams(action string, params map[string]string, queue *Queue, resp interface{}) error {
  uri := c.getURL(queue)
  params["Action"] = action
  request := generateGetRequest(uri, params)
  aws.SignRequest(c, request)
  return makeRequest(request, resp)
}
