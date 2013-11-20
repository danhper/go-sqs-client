package sqs

import (
  "strings"
)

type Queue struct {
  name       string
  url        string
  attributes map[string]string
}

func makeQueueFromURL(url string) *Queue {
  splittedUrl := strings.Split(url, "/")
  return &Queue {
    name       : splittedUrl[len(splittedUrl) - 1],
    url        : url,
    attributes : make(map[string]string),
  }
}

func makeQueuefromURLs(urlList []string) []*Queue{
  var queues []*Queue
  for _, url := range urlList {
    queues = append(queues, makeQueueFromURL(url))
  }
  return queues
}
