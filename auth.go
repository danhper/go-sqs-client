package aws

import (
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "fmt"
  "net/http"
  "net/url"
  "sort"
  "strings"
)

type Auth struct {
  accessKey, secretKey string
}

func hash(str string) string {
  hasher := sha256.New()
  hasher.Write([]byte(str))
  return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateHashedRequest(request *HTTPRequest) string {
  return hash(GenerateCanonicalRequest(request))
}

func GenerateCanonicalRequest(request *HTTPRequest) string {
  canonicalRequest := request.Method + "\n"
  canonicalRequest += request.URL.Path + "\n"
  canonicalRequest += generateCanonicalQuery(request.URL.Query()) + "\n"
  canonicalRequest += generateCanonicalHeaders(request.Header) + "\n"
  canonicalRequest += generateSignedHeaders(request.Header) + "\n"
  canonicalRequest += hash(request.payload)

  return canonicalRequest
}

func generateSignedHeaders(headers http.Header) string {
  var signedHeaders []string
  for k, _ := range headers {
    key := strings.ToLower(k)
    if strings.HasPrefix(key, "x-") || key == "host" || key == "content-type" {
      signedHeaders = append(signedHeaders, key)
    }
  }
  sort.StringSlice(signedHeaders).Sort()
  return strings.Join(signedHeaders, ";")
}

func generateCanonicalQuery(query url.Values) string {
  var queryParams []string
  for key, value := range query {
    queryParams = append(queryParams, urlencode(key) + "=" + urlencode(value[0]))
  }
  sort.StringSlice(queryParams).Sort()
  return strings.Join(queryParams, "&")
}

func generateCanonicalHeaders(headers http.Header) string {
  var canonicalHeaders []string
  for key, value := range headers {
    canonicalKey := strings.ToLower(key)
    canonicalValue := strings.Trim(value[0], " \t\n")
    canonicalHeaders = append(canonicalHeaders, canonicalKey + ":" + canonicalValue + "\n")
  }
  sort.StringSlice(canonicalHeaders).Sort()
  return strings.Join(canonicalHeaders, "")
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

// https://gist.github.com/hnaohiro/4627658
func urlencode(s string) (result string) {
  for _, c := range(s) {
    if c <= 0x7f { // single byte
      result += fmt.Sprintf("%%%X", c)
    } else if c > 0x1fffff { // quaternary byte
      result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
        0xf0 + ((c & 0x1c0000) >> 18),
        0x80 + ((c & 0x3f000) >> 12),
        0x80 + ((c & 0xfc0) >> 6),
        0x80 + (c & 0x3f),
        )
    } else if c > 0x7ff { // triple byte
      result += fmt.Sprintf("%%%X%%%X%%%X",
        0xe0 + ((c & 0xf000) >> 12),
        0x80 + ((c & 0xfc0) >> 6),
        0x80 + (c & 0x3f),
        )
    } else { // double byte
      result += fmt.Sprintf("%%%X%%%X",
        0xc0 + ((c & 0x7c0) >> 6),
        0x80 + (c & 0x3f),
        )
    }
  }

  return result
}
