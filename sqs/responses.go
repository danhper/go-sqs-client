package sqs

type ListQueuesResponse struct {
  QueueUrl []string `xml:"ListQueuesResult>QueueUrl"`
  ResponseMetadata
}

type ResponseMetadata struct {
  RequestId string
  BoxUsage float64
}
