// gcloud functions deploy sign --entry-point=SignedUrl --runtime=go111 --trigger-http --quiet

package main

import (
        "os"
        "fmt"
        "log"
        "net/http"
        "github.com/akiirobot/sign-storage-url/src/sign"
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

  bucket := r.FormValue("bucket")
  timeStamp := r.FormValue("time")

  url, err := sign.Sign(serviceAccount, whiteListBucket, objectName, method, bucket, timeStamp)
  if err != nil {
    log.Fatalln(err)
    http.Error(w, "403 - Status Forbidden - " + err.Error(), http.StatusForbidden)
  }

  fmt.Fprintf(w, url)
}



