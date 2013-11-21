package aws

import (
  "os"
)

type Credentials struct {
  accessKey, secretKey string
}

func MakeCredentials(accessKey string, secretKey string) *Credentials {
  return &Credentials {
    accessKey: accessKey,
    secretKey: secretKey,
  }
}

func MakeCredentialsFromEnv() *Credentials {
  return &Credentials {
    accessKey: os.Getenv("AWS_ACCESS_KEY"),
    secretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
  }
}
