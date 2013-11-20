package aws

import (
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
  "net/http"
  "net/url"
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

func updateHeadersForAuth(c Client, request *HTTPRequest) {
  request.Header.Set("x-amz-date", getLongTime(request.createdTime))
  if request.Header.Get("Content-Type") == "" {
    request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
  }
  credential := "Credential=" + getCredentials(c.Credentials().accessKey, request.createdTime, c.RegionName(), c.ServiceName())
  signedHeaders := "SignedHeaders=" + generateSignedHeaders(request.Header)
  signature := "Signature=" + getSignature(c, request)
  authString := algorithm + " " + credential + ", " + signedHeaders + ", " + signature
  request.Header.Set("Authorization", authString)
}

func updateQueryWithAuth(c Client, request *HTTPRequest) {
  credentials := getCredentials(c.Credentials().accessKey, request.createdTime,
    c.RegionName(), c.ServiceName())
  query := request.URL.Query()
  query.Set("X-Amz-Algorithm", algorithm)
  query.Set("X-Amz-Credential", credentials)
  query.Set("X-Amz-Date", getLongTime(request.createdTime))
  query.Set("X-Amz-SignedHeaders",  generateSignedHeaders(request.Header))
  request.URL.RawQuery = query.Encode()
}

func getCredentials(accessKey string, t time.Time,
  regionName string, serviceName string) string {
  scope := getScope(t, regionName, serviceName)
  return accessKey + "/" + scope
}

func getScope(t time.Time, regionName string, serviceName string) string {
  return getShortTime(t) + "/" + regionName + "/" + serviceName + "/aws4_request"
}

func generateStringToSign(regionName string, serviceName string,
  t time.Time, hashedRequest string) string {
  scope := getScope(t, regionName, serviceName)
  return algorithm + "\n" + getLongTime(t) + "\n" + scope + "\n" + hashedRequest
}

func generateHashedRequest(request *HTTPRequest) string {
  return hash(generateCanonicalRequest(request))
}

func prependSlash(url string) string {
  if strings.HasPrefix(url, "/") {
    return url
  } else {
    return "/" + url
  }
}

func generateCanonicalRequest(request *HTTPRequest) string {
  canonicalRequest := request.Method + "\n"
  canonicalRequest += prependSlash(request.URL.Path) + "\n"
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
    q := url.Values{key: {value[0]}}
    queryParams = append(queryParams, q.Encode())
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

func getSignature(c Client, r *HTTPRequest) string {
  hashedRequest := generateHashedRequest(r)
  stringToSign := generateStringToSign(c.RegionName(), c.ServiceName(), r.createdTime, hashedRequest)
  signingKey := getSignatureKey(c.Credentials().secretKey, r.createdTime, c.RegionName(), c.ServiceName())
  return generateSignature(signingKey, stringToSign)
}

func generateSignature(signingKey []byte, stringToSign string) string {
  return hex.EncodeToString(sign(signingKey, stringToSign))
}

func getSignatureKey(secretKey string, t time.Time,
  regionName string, serviceName string) []byte {
  kDate := sign([]byte("AWS4" + secretKey), getShortTime(t))
  kRegion := sign(kDate, regionName)
  kService := sign(kRegion, serviceName)
  kSigining := sign(kService, "aws4_request")
  return kSigining
}
