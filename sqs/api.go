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

type MessageResponse struct {
  Message
  ResponseMetadata
}

type ResponseMetadata struct {
  RequestId string
  BoxUsage float64
}

type listQueuesResponse struct {
  QueueUrl []string `xml:"ListQueuesResult>QueueUrl"`
  ResponseMetadata
}

type createQueueResponse struct {
  QueueUrl string `xml:"CreateQueueResult>QueueUrl"`
  ResponseMetadata
}

type getQueueUrlResponse struct {
  QueueUrl string `xml:"GetQueueUrlResult>QueueUrl"`
  ResponseMetadata
}

type sendMessageResponse struct {
  MD5 string `xml:"SendMessageResult>MD5OfMessageBody"`
  Id string `xml:"SendMessageResult>MessageId"`
  ResponseMetadata
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
  if err != nil {
    return nil, err
  }
  return &QueueListResponse{makeQueuefromURLs(resp.QueueUrl), resp.ResponseMetadata}, nil
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
  if err != nil {
    return nil, err
  }
  return &QueueResponse{*makeQueueFromURL(resp.QueueUrl), resp.ResponseMetadata}, nil
}

func (c *SqsClient) SendMessage(queue *Queue, body string) (*MessageResponse, error) {
  return c.SendMessageWithDelay(queue, body, -1)
}

func (c *SqsClient) SendMessageWithDelay(queue *Queue, body string, delay int) (*MessageResponse, error) {
  resp := &sendMessageResponse{}
  params := make(map[string]string)
  params["MessageBody"] = body
  if delay >= 0 {
    params["Delay"] = fmt.Sprintf("%d", delay)
  }
  err := c.makePostRequestWithParams("SendMessage", params, queue, resp)
  if err != nil {
    return nil, err
  }
  message := Message{MessageId: resp.Id, Body: body, MD5: resp.MD5}
  if delay >= 0 {
    message.Delay = delay
  }
  return &MessageResponse{message, resp.ResponseMetadata}, nil
}

func (c *SqsClient) GetQueueUrl(name string) (*QueueResponse, error) {
  resp := &getQueueUrlResponse{}
  params := map[string]string {
    "QueueName": name,
  }
  err := c.makeGetRequestWithParams("GetQueueUrl", params, nil, resp)
  if err != nil {
    return nil, err
  }
  return &QueueResponse{*makeQueueFromURL(resp.QueueUrl), resp.ResponseMetadata}, nil
}
