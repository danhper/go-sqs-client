package sqs

type Message struct {
  Id        string
  MessageId string
  Body      string
  MD5       string
  Delay     int
}
