package sqs

import (
  "strings"
)

type Queue struct {
  Name       string
  URL        string
  Attributes map[string]string
}

func NewQueueWithAttrs(name string, attributes map[string]string) *Queue {
  return &Queue {
    Name: name,
    Attributes: attributes,
  }
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
