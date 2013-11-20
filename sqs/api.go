package sqs

import (
  "fmt"
)

type QueueResponse struct {
  Queue
  ResponseMetadata
}

type QueueListResponse struct {
  Queues []*Queue
  ResponseMetadata
}

type listQueuesResponse struct {
  QueueUrl []string `xml:"ListQueuesResult>QueueUrl"`
  ResponseMetadata
}

type createQueueResponse struct {
  QueueUrl string `xml:"CreateQueueResult>QueueUrl"`
  ResponseMetadata
}

type ResponseMetadata struct {
  RequestId string
  BoxUsage float64
}

func (c *SqsClient) ListQueues() (*QueueListResponse, error) {
  return c.ListQueuesWithPrefix("")
}

func (c *SqsClient) ListQueuesWithPrefix(prefix string) (*QueueListResponse, error) {
  resp := &listQueuesResponse{}
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

func (c *SqsClient) CreateQueue(name string) (*QueueResponse, error) {
  queue := &Queue{
    Name: name,
    Attributes: make(map[string]string),
  }
  return c.CreateQueueWithAttributes(queue)
}

func (c *SqsClient) CreateQueueWithAttributes(queue *Queue) (*QueueResponse, error) {
  resp := &createQueueResponse{}
  params := map[string]string {
    "QueueName": queue.Name,
  }
  i := 1
  for key, value := range queue.Attributes {
    params[fmt.Sprintf("Attributes.%d.Name", i)]  = key
    params[fmt.Sprintf("Attributes.%d.Value", i)] = value
    i++
  }
  err := c.makePostRequestWithParams("CreateQueue", params, nil, resp)
  if err == nil {
    return &QueueResponse{*makeQueueFromURL(resp.QueueUrl), resp.ResponseMetadata}, nil
  } else {
    return nil, err
  }
}
