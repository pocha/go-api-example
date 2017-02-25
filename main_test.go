package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
  "encoding/json"
  "bytes"
  "fmt"
)



func TestSMSHandler(t *testing.T) {
	assert := assert.New(t)

  validData := SMS {
                  From: "919538384545",
                  To:   "919845350048",
                  Message:  "hello how are you",
                }


  jsonData, _ := json.Marshal(validData)
  validHash := buildHash(jsonData)
  fmt.Println("validHash", validHash)

	successTests := []struct {
		description        string
		url                string
    input              SMS
		expectedStatusCode int
    output             map[string]string
	}{
		{
			description:        "valid test data",
			url:                "/outbound/sms?hash=" + validHash,
      input:              validData,
			expectedStatusCode: 200,
      output:             map[string]string{ "message" : "outbound sms ok", "error" : "" },
		}, 	
  }

	for _, tc := range successTests {
    
    input, err := json.Marshal(tc.input)
    if err != nil {
      fmt.Println("error encounted", err)
    }
     
		req, err := http.NewRequest("POST", tc.url, bytes.NewBuffer(input))
		assert.NoError(err)
    req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		SMSHandler(w, req)

		assert.Equal(tc.expectedStatusCode, w.Code, tc.description)
    
    output, _ := json.Marshal(tc.output)
		assert.Equal(string(output), w.Body.String(), tc.description)
	}
}
