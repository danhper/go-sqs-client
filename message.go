package sqs

type Message struct {
  MessageId string `xml:"MessageId"`
  Body string `xml:"Body"`
  MD5 string `xml:"MD5OfBody"`
  ReceiptHandle string `xml:"ReceiptHandle"`
  Attribute []Attribute `xml:"Attribute"`
  Delay int
  Queue *Queue
}

type Attribute struct {
  Name string `xml:"ReceiveMessageResult>Message>Attribute>Name"`
  Value string `xml:"ReceiveMessageResult>Message>Attribute>Value"`
}
