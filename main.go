package main

import (
	"encoding/json"
	"log"
	"net/http"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
  "time"
  "strings"
  "io/ioutil"
  "bytes"
)

const auth_username = "ashish"
const auth_password = "34b5db01-8bc3-4960-92d7-d83155a42dfe"

type (
  SMS struct {
    From string `json:"from"`
    To   string `json:"to"`
    Message string `json:"message"`
  }
)

func buildHash(data []byte) string{
  mac := hmac.New(sha256.New, []byte(auth_password))
  mac.Write(data)
  return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func exitWithJson(w http.ResponseWriter, message string, error string) {
  _output := map[string]string { "message": message, "error": error}
  output, _ := json.Marshal(_output)
  w.Header().Set("Content-Type", "application/json")
  w.Write(output)
}

func sendToKannel(data SMS) bool{
  //simulating 1 sec delay
  time.Sleep(1 * time.Second)
  return true
}

func getBodyJson(r *http.Request) []byte {
  bodyBytes, _ := ioutil.ReadAll(r.Body)
  r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
  return bodyBytes
}

func SMSHandler(w http.ResponseWriter, r *http.Request) {

  if r.Header.Get("Content-Type") != "application/json" {
    http.Error(w, "only json data accepted", 400)
    return
  }
  
  var data SMS
  _ = json.Unmarshal(getBodyJson(r), &data)
  log.Println("sms data", data)

  var sms_map map[string]interface{}
  err := json.Unmarshal(getBodyJson(r),&sms_map)
  if err != nil {
    http.Error(w, "invalid json", 400)
    return
  }
  log.Println("sms_map", sms_map)

    
  hash := r.FormValue("hash")
  if hash == "" {
    http.Error(w, "hash value as parameter expected", 400)
    return
  }
  
  b, _ := json.Marshal(data)
  expectedHash := buildHash(b)
  log.Println("hash value", hash)
  log.Println("expt value", expectedHash)

  //check hash
  if strings.Compare(hash, expectedHash) != 0{
    http.Error(w, "authentication failed", 400)
    return
  }
  

  //check for elements presence
  _, ok := sms_map["from"]
  if ok == false {
    exitWithJson(w,"","from field absent")
    return
  }
  _, ok = sms_map["to"]
  if ok == false {
    exitWithJson(w,"","to field absent")
    return
  }
  _, ok = sms_map["message"]
  if ok == false {
    exitWithJson(w,"","message field absent")
    return
  }

  //validate elements - not sure why the count yields one more digit
  if strings.Count(data.From,"") != 13 {
    exitWithJson(w,"","invalid from field")
    return
  }
  if strings.Count(data.To,"") != 13 {
    exitWithJson(w,"","invalid to field")
    return
  }
  //no idea how to check if message text is utf16 or gsm7 supported
  
  //everything looks fine - send data to kannel
  if sendToKannel(data)  {
    exitWithJson(w,"outbound sms ok","")
  } else {
    exitWithJson(w,"","unknown error occured")
  }
 
}

func main() {
  http.HandleFunc("/outbound/sms", SMSHandler)

  log.Println("starting server")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
