package sqs

import (
  "fringe81.com/tuvistavie/aws"
)

type SqsClient struct {
  aws.BaseClient
}

func (s *SqsClient) ServiceName() string {
  return "sqs"
}

func MakeClient(credentials *aws.Credentials, regionName string) *SqsClient {
  return &SqsClient {aws.MakeBaseClient(regionName, credentials)}
}

func MakeClientWithBaseCreditentials(accessKey string,
  secretKey string, regionName string) *SqsClient {
  return MakeClient(aws.MakeCredentials(accessKey, secretKey), regionName)
}
