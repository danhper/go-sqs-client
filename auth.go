package aws

import (
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "fmt"
  "net/http"
  "net/url"
  "regexp"
  "sort"
  "strings"
  "time"
)

const algorithm = "AWS4-HMAC-SHA256"

func hash(str string) string {
  hasher := sha256.New()
  hasher.Write([]byte(str))
  return hex.EncodeToString(hasher.Sum(nil))
}

func getLongTime(t time.Time) string {
  return t.Format("20060102T150405Z")
}

func getShortTime(t time.Time) string {
  return t.Format("20060102")
}

func setSignedHeaders(c Client, request *HTTPRequest) {
  credentials := getCredentials(c.Credentials().accessKey, request.createdTime,
    c.RegionName(), c.ServiceName())
  headers := request.Header
  headers.Set("X-Amz-Algorithm", algorithm)
  headers.Set("X-Amz-Credential", credentials)
  headers.Set("X-Amz-Date", getLongTime(request.createdTime))
  headers.Set("X-Amz-SignedHeaders",  generateSignedHeaders(headers))
}

func getCredentials(accessKey string, t time.Time,
  regionName string, serviceName string) string {
  scope := getScope(t, regionName, serviceName)
  return accessKey + "/" + scope
}

func getScope(t time.Time, regionName string, serviceName string) string {
  return getShortTime(t) + "/" + regionName + "/" + serviceName + "/aws4_request"
}

func GenerateStringToSign(regionName string, serviceName string,
  t time.Time, hashedRequest string) string {
  scope := getScope(t, regionName, serviceName)
  return algorithm + "\n" + getLongTime(t) + "\n" + scope + "\n" + hashedRequest
}

func GenerateHashedRequest(request *HTTPRequest) string {
  return hash(generateCanonicalRequest(request))
}

func generateCanonicalRequest(request *HTTPRequest) string {
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
    if key == "x-amz-date" || key == "host" || key == "content-type" {
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

func GenerageSignature(signingKey []byte, stringToSign string) string {
  return hex.EncodeToString(sign(signingKey, stringToSign))
}

func GetSignatureKey(secretKey string, t time.Time,
  regionName string, serviceName string) []byte {
  kDate := sign([]byte("AWS4" + secretKey), getShortTime(t))
  kRegion := sign(kDate, regionName)
  kService := sign(kRegion, serviceName)
  kSigining := sign(kService, "aws4_request")
  return kSigining
}

// https://gist.github.com/hnaohiro/4627658
func urlencode(s string) (result string) {
  authorizedChars := regexp.MustCompile("[A-Za-z0-9_.~-]")
  for _, c := range(s) {
    char := fmt.Sprintf("%c", c)
    if authorizedChars.FindStringIndex(char) != nil {
      result += char
    } else if c <= 0x7f { // single byte
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
