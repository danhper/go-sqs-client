package aws

type Credentials struct {
  accessKey, secretKey string
}

func MakeCredentials(accessKey string, secretKey string) *Credentials {
  return &Credentials {
    accessKey: accessKey,
    secretKey: secretKey,
  }
}
