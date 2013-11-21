# go-sqs-sdk

This package is an SDK for the AWS SQS service.

## API support

Only the following API features are implemented in this API.

* CreateQueue
* DeleteMessage
* DeleteQueue
* GetQueueUrl
* ListQueues
* ReceiveMessage
* SendMessage

All the above API queries are supported with all their options.

## Installation

Copy the package to your `$GOPATH` with the path below

```
mkdir -p $GOPATH/fringe81.com/tuvistavie
cp -r go-aws-sdk $GOPATH/fringe81.com/tuvistavie/aws
```

You can then import the packages `aws` and `sqs` using

```go
import (
  "fringe81.com/tuvistavie/aws"
  "fringe81.com/tuvistavie/aws/sqs"
)
```

## Sample usage

Here is a sample usage with the basic API functionalities exposed.

```go
package main

import (
  "fmt"

  "fringe81.com/tuvistavie/aws"
  "fringe81.com/tuvistavie/aws/sqs"
)

func main() {
  client := sqs.MakeClientAndCreditentials("ACCESS_KEY", "SECRET_KEY", aws.Tokyo)

  queueResp, err := client.CreateQueue("foobar")
  if err != nil {
    fmt.Println(err)
    return
  }

  queue := queueResp.Queue

  client.SendMessage(queue, "foobarbaz-%d")

  client.DeleteQueue(queue)

  queueResp, _ = client.GetQueueUrl("foo")
  msgs, _ := client.ReceiveMessageWithParams(queueResp.Queue, map[string]string {
    "WaitTimeSeconds": "10",
    "MaxNumberOfMessages": "5",
    "VisibilityTimeout" : "10",
  })

  for _, msg := range msgs.Messages {
    fmt.Println(msg.Body)
  }

  if len(msgs.Messages) > 0 {
    msg := msgs.Messages[0]
    client.DeleteMessage(msg)
  }
}
```

## Public API

The `sqs` package exposes the functionalities to access and use the SQS API, while the `aws` package only exposes the credentials and region settings at the present.

All the API calls return the result of the request, and an error. One should check that the returned error is `nil` before using the result.
When the request succeeds, the results always contain a `Metadata` object containing the `RequestId` for the request.

### Client creation

The following methods are available to create a client for `sqs`.

```go
func MakeClient(credentials *aws.Credentials, regionName string) *SqsClient
func MakeClientWithProtocol(credentials *aws.Credentials, regionName string, protocol string) *SqsClient
func MakeClientAndCreditentials(accessKey string, secretKey string, regionName string) *SqsClient
func MakeClientAndCreditentialsWithProtocol(accessKey string, secretKey string, regionName string, protocol string) *SqsClient
```

The region can be given using the constants available in the `aws` package.
The protocol string should be `http` or `https`. If no protocol is specified, `http` will be used by default.

### CreateQueue

```go
func (c *SqsClient) CreateQueue(name string) (*QueueResponse, error)
func (c *SqsClient) CreateQueueWithAttributes(queue Queue) (*QueueResponse, error)
```

The `CreateQueue` method only takes the name of the queue to create. The `CreateQueueWithAttributes` method takes a `Queue` structure defined as follow

```go
type Queue struct {
  Name       string
  URL        string
  Attributes map[string]string
}
```

and can be created using the `NewQueueWithAttrs` function.

```go
func NewQueueWithAttrs(name string, attributes map[string]string) *Queue
```

the attributes should be the keys and values described in the [CreateQueue documentation](http://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/Query_QueryCreateQueue.html), with the same case. However, the integer should be converted to strings.

Example attributes:

```go
map[string]string {
  "DelaySeconds"      = "10",
  "VisibilityTimeout" = "8"
}
```

On success, the response is a `QueueResponse` containing a `Queue` pointer.

### DeleteMessage

```go
func (c *SqsClient) DeleteMessage(msg Message) (*EmptyResponse, error)
```

The message to be deleted must contain a valid `ReceiptHandle` or the request will fail.
The response does not contains any special data.

### DeleteQueue

```go
func (c *SqsClient) DeleteQueue(queue *Queue) (*EmptyResponse, error)
```

The queue should have a valid `URL` for the request to work.
The response does not contains any special data.

### GetQueueUrl

```go
func (c *SqsClient) GetQueueUrl(name string) (*QueueResponse, error)
```

The `GetQueueUrl` takes a queue name and returns a `QueueResponse` pointer containing the corresponding queue with a valid URL.
If the queue does not exist, an error will be returned.

### ListQueues

```go
func (c *SqsClient) ListQueues() (*QueueListResponse, error)
func (c *SqsClient) ListQueuesWithPrefix(prefix string) (*QueueListResponse, error)
```

The `ListQueues` and `ListQueuesWithPrefix` methods return a `QueueListResponse` containing a list of queue pointers with a valid URL. If `ListQueuesWithPrefix`, the `prefix` string will be used when getting the list of queues.

### ReceiveMessage

```go
func (c *SqsClient) ReceiveMessage(queue *Queue) (*MessageListResponse, error)
func (c *SqsClient) ReceiveMessageWithParams(queue *Queue, params map[string]string) (*MessageListResponse, error)
func (c *SqsClient) ReceiveMessageWithParamsAndAttrs(queue *Queue, params map[string]string, attrs []string) (*MessageListResponse, error)
```

The `ReceiveMessage`, `ReceiveMessageWithParams` and `ReceiveMessageWithParamsAndAttrs` methods all receive messages from the given queue. The `ReceiveMessage` will used the default parameters [described in the documentation](http://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/Query_QueryReceiveMessage.html).

To customize the parameters, a `map[string]string` can be passed to the `ReceiveMessageWithParams` method. The keys should be `MaxNumberOfMessages`, `VisibilityTimeout` or `WaitTimeSeconds` with containing values conform to the API documentation. However, integers should be transformed to string.

Finally, a slice containing the wanted attributes names can be passed as third argument of the `ReceiveMessageWithParamsAndAttrs` method.

On success, the method will return a `MessageListResponse` containing a slice of `Message` objects with valid `ReceiptHandle`.

### SendMessage

```go
func (c *SqsClient) SendMessage(queue *Queue, body string) (*MessageResponse, error)
func (c *SqsClient) SendMessageWithDelay(queue *Queue, body string, delay int) (*MessageResponse, error)
```

The `SendMessage` and `SendMessageWithDelay` methods will send a message with the given `body` to the given `Queue`. The `DelaySeconds` can be passed to `SendMessageWithDelay` to specify when the message should become avilable for processing.