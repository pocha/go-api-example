package main

import (
	"encoding/json"
	"log"
	"net/http"
  "crypto/hmac"
  "crypto/sha256"
)

const auth_username = "ashish"
const auth_password = "34b5db01-8bc3-4960-92d7-d83155a42dfe"

type (
  SMS struct {
    from, to uint64
    message string
  }

  Output struct {
    error, message string
  }
)


func SMSHandler(w http.ResponseWriter, r *http.Request) {

  if r.Header.Get("Content-Type") != "application/json" {
    http.Error(w, "only json data accepted", 400)
    return
  }
  
  var data []byte 
  err := json.NewDecoder(r.Body).Decode(data)
  if err != nil {
    http.Error(w, "invalid json", 400)
    return
  }


  hash := r.FormValue("hash")
  if hash == "" {
    http.Error(w, "hash value as parameter expected", 400)
    return
  }

  mac := hmac.New(sha256.New, []byte(auth_password))
  mac.Write(data)
  expectedMAC := mac.Sum(nil)

  //check hash
  if hmac.Equal([]byte(hash),expectedMAC) == false {
    http.Error(w, "authentication failed", 400)
    return
  }

  output, err := json.Marshal( map[string]string { "message" : "outbound sms ok", "error": "" } )

  w.Header().Set("Content-Type", "application/json")
  w.Write(output)
}

func main() {
  http.HandleFunc("/outbound/sms", SMSHandler)

  log.Println("starting server")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
