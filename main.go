// gcloud functions deploy sign --entry-point=SignedUrl --runtime=go111 --trigger-http --quiet

package sign

import (
        "os"
        "fmt"
        "log"
        "time"
        "strings"
        "strconv"
        "net/http"
        "io/ioutil"
        "golang.org/x/oauth2/google"
        "cloud.google.com/go/storage"
      )

func SignedUrl(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    fmt.Printf("status method %s is not allowd, only allowed post", r.Method)
    http.Error(w, "405 - Status Method Not Allowed - " + r.Method, http.StatusMethodNotAllowed)
    return
  }

  serviceAccount := os.Getenv("SERVICE_JSON_FILE")
  whiteListBucket := os.Getenv("BUCKET_NAME")

  objectName := r.FormValue("object")
  method := r.FormValue("method")

  bucketName, err := findBucket(whiteListBucket, r.FormValue("bucket"))
  if err != nil {
    log.Fatalln("bucket name is not valid, err: %v", err)
    http.Error(w, "403 - Status Forbidden - bucket name is not valid", http.StatusForbidden)
    return
  }

  jsonKey, err := ioutil.ReadFile(serviceAccount)
  if err != nil {
    log.Fatalln("cannot read the JSON key file, err: %v", err)
    http.Error(w, "403 - Status Forbidden - service key is missing", http.StatusForbidden)
    return
  }

  expireTime, err := strconv.Atoi(r.FormValue("time"))
  if err != nil {
    expireTime = 15
  }

  conf, err := google.JWTConfigFromJSON(jsonKey)
  if err != nil {
    log.Fatalln("google.JWTConfigFromJSON: %v", err)
    http.Error(w, "403 - Status Forbidden - service key is not valid", http.StatusForbidden)
    return
  }

  // https://github.com/googleapis/google-cloud-go/blob/master/storage/storage.go#L157
  opts := &storage.SignedURLOptions{
    Scheme:         storage.SigningSchemeV4,
    Method:         method,
    GoogleAccessID: conf.Email,
    PrivateKey:     conf.PrivateKey,
    Expires:        time.Now().Add(time.Duration(expireTime) * time.Minute),
  }

  u, err := storage.SignedURL(bucketName, objectName, opts)
  if err != nil {
    log.Fatalln("Unable to generate a signed URL: %v", err)
    http.Error(w, "403 - Status Forbidden - unable to generate a signed url", http.StatusForbidden)
    return
  }

  fmt.Fprintf(w, u)
}

func findBucket(whiteList, name string) (string, error) {
  lists := strings.Split(whiteList, ":")

  if name == "" {
    if len(lists) == 1 && lists[0] != "" {
      return lists[0], nil
    }

    return "", fmt.Errorf("bucket name is missing, bucket: %s, white list: %s", name, whiteList)
  }

  if len(lists) == 1 && lists[0] == "" {
    return name, nil
  }

  for _, v := range lists {
    if v == name {
      return name, nil
    }
  }

  return "", fmt.Errorf("bucket name is not available, bucket: %s, white list: %s", name, whiteList)
}



