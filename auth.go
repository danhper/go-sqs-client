package aws

import (
  "crypto/hmac"
  "crypto/sha256"
)

type Auth struct {
  accessKey, secretKey string
}

func sign(key []byte, msg string) []byte {
  h := hmac.New(sha256.New, key)
  h.Write([]byte(msg))
  return h.Sum(nil)
}

func getSignatureKey(secretKey string, dateStamp string,
  regionName string, serviceName string) []byte {
  kDate := sign([]byte("AWS4" + secretKey), dateStamp)
  kRegion := sign(kDate, regionName)
  kService := sign(kRegion, serviceName)
  kSigining := sign(kService, "aws4_request")
  return kSigining
}
