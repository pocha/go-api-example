package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
    "regexp"
)

const auth_username = "ashish"
const auth_password = "34b5db01-8bc3-4960-92d7-d83155a42dfe"

type (
	SMS struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Message string `json:"message"`
	}
)

func buildHash(data []byte) string {
	mac := hmac.New(sha256.New, []byte(auth_password))
	mac.Write(data)
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func exitWithJson(w http.ResponseWriter, message string, error string) {
	_output := map[string]string{"message": message, "error": error}
	output, _ := json.Marshal(_output)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func sendToKannel(data SMS) bool {
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
	err := json.Unmarshal(getBodyJson(r), &data)
	if err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	log.Println("sms data", data)


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
	if strings.Compare(hash, expectedHash) != 0 {
		http.Error(w, "authentication failed", 400)
		return
	}

	//check for elements presence
	var smsMap map[string]interface{}
	err = json.Unmarshal(getBodyJson(r), &smsMap)
	log.Println("smsMap", smsMap)

	_, ok := smsMap["from"]
	if ok == false {
		exitWithJson(w, "", "from field absent")
		return
	}
	_, ok = smsMap["to"]
	if ok == false {
		exitWithJson(w, "", "to field absent")
		return
	}
	_, ok = smsMap["message"]
	if ok == false {
		exitWithJson(w, "", "message field absent")
		return
	}

	//validate from, to elements
    regex, _ := regexp.Compile("^([0-9]){12}$")

	if regex.MatchString(data.From) == false {
		exitWithJson(w, "", "invalid from field")
		return
	}
	if regex.MatchString(data.To) == false {
		exitWithJson(w, "", "invalid to field")
		return
	}
	//no idea how to check if message text is utf16 or gsm7 supported

	//everything looks fine - send data to kannel
	if sendToKannel(data) {
		exitWithJson(w, "outbound sms ok", "")
	} else {
		exitWithJson(w, "", "unknown error occured")
	}

}

func main() {
	http.HandleFunc("/outbound/sms", SMSHandler)

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
