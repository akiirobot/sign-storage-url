package sign

import (
        "fmt"
        "time"
        "strings"
        "strconv"
        "io/ioutil"
        "golang.org/x/oauth2/google"
        "cloud.google.com/go/storage"
      )

func Sign(serviceAccount, whiteListBucket, objectName, method, bucket, timeStamp string) (string, error) {
  bucketName, err := findBucket(whiteListBucket, bucket)
  if err != nil {
    return "", fmt.Errorf("bucket name is not valid, bucket name is not valid, err: %v", err)
  }

  jsonKey, err := ioutil.ReadFile(serviceAccount)
  if err != nil {
    return "", fmt.Errorf("service key is missing, cannot read the JSON key file, err: %v", err)
  }

  expireTime, err := strconv.Atoi(timeStamp)
  if err != nil {
    expireTime = 15
  }

  conf, err := google.JWTConfigFromJSON(jsonKey)
  if err != nil {
    return "", fmt.Errorf("service key is not valid, google.JWTConfigFromJSON: %v", err)
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
    return "", fmt.Errorf("Unable to generate a signed URL: %v", err)
  }

  return u, nil
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

