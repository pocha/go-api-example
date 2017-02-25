package main

import (
	"encoding/json"
	"log"
	"net/http"
  "crypto/hmac"
  "crypto/sha256"
  "github.com/pocha/sms-gateway-go-sample/sms"
)

const auth_username = "ashish"
const auth_password = "34b5db01-8bc3-4960-92d7-d83155a42dfe"



type Output struct {
  error, message string
}


func SMSHandler(w http.ResponseWriter, r *http.Request) {
	
  if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "only json data accepted", 400)
		return
	}

  data := r.Body
  hash = r.FormValue("hash")
  mac := hmac.New(sha256.New, auth_password)
  mac.Write(data)
  expectedMAC := mac.Sum(nil)
  
  //check hash
  if hmac.Equal(hash,expectedMAC) == false {
		http.Error(w, "authentication failed", 400)
		return
	}
  
	err := json.NewDecoder(data).Decode(models.SMS)
	if err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
  
  output := json.Marshal(map[string]string { "message" : "outbound sms ok", "error": "" } )

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func main() {
	http.HandleFunc("/outbound/sms", SMSHandler)

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":80", nil))
}
