package sqs

import (
  "strings"
)

type Queue struct {
  Name       string
  URL        string
  Attributes map[string]string
}

type QueueListResponse struct {
  Queues []*Queue
  ResponseMetadata
}

func makeQueueFromURL(url string) *Queue {
  splittedUrl := strings.Split(url, "/")
  return &Queue {
    Name       : splittedUrl[len(splittedUrl) - 1],
    URL        : url,
    Attributes : make(map[string]string),
  }
}

func makeQueuefromURLs(urlList []string) []*Queue{
  var queues []*Queue
  for _, url := range urlList {
    queues = append(queues, makeQueueFromURL(url))
  }
  return queues
}
