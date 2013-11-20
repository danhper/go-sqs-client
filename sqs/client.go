package sqs

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"

  "fringe81.com/tuvistavie/aws"
)

const endpoint = "sqs.%s.amazonaws.com"

type xmlErrors struct {
  RequestId string
  Errors []Error `xml:"Errors>Error"`
}

type Error struct {
  StatusCode int
  Code string
  Message string
  RequestId string
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

func makeRequest(request *aws.HTTPRequest, resp interface{}) error {
  var r *http.Response
  var body []byte
  var err error

  client := &http.Client{}

  r, err = client.Do(request.Request)

  if err != nil {
    return nil
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
  errors := xmlErrors{}
  xml.Unmarshal(body, errors)
  var err Error
  if len(errors.Errors) > 0 {
    err = errors.Errors[0]
  }
  err.RequestId = errors.RequestId
  err.StatusCode = r.StatusCode
  if err.Message == "" {
    err.Message = r.Status
  }
  return &err
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
  uri :=  c.getURL(queue)
  params["Action"] = action
  paramValues := url.Values{}
  for k, v := range params {
    paramValues[k] = []string{v}
  }
  request, _ := aws.NewHTTPRequest("POST", uri, paramValues.Encode())
  aws.SignRequest(c, request)
  return makeRequest(request, resp)
}

func (c * SqsClient) getURL(queue *Queue) string {
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

func (c *SqsClient) ListQueues() (*QueueListResponse, error) {
  return c.ListQueuesWithPrefix("")
}

func (c *SqsClient) ListQueuesWithPrefix(prefix string) (*QueueListResponse, error) {
  resp := &ListQueuesResponse{}
  params := make(map[string]string)
  if prefix != "" {
    params["QueueNamePrefix"] = prefix
  }
  err := c.makeGetRequestWithParams("ListQueues", params, nil, resp)
  if err == nil {
    return &QueueListResponse{makeQueuefromURLs(resp.QueueUrl), resp.ResponseMetadata}, nil
  } else {
    return nil, err
  }
}
