package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
  "github.com/pocha/sms-gateway-go-sample/main"
	"github.com/stretchr/testify/assert"
)

type Data struct {
  from, to uint64
  message string
}



func TestSMSHandler(t *testing.T) {
	assert := assert.New(t)


	tests := []struct {
		description        string
		url                string
    input              Data
		expectedStatusCode int
    output             map[string]string
	}{
		{
			description:        "valid test data",
			url:                "/outbound/sms",
      input:               {
                              from: 919538384545,
                              to:   919845350048,
                              message:  "hello how are you",
                            },
			expectedStatusCode: 200,
      output:              { "message" : "outbound sms ok", "error" : "" },
		}, 	
  }

	for _, tc := range tests {

    input, err := json.Marshal(tc.input)
    output, err := json.Marshal(tc.output)

		req, err := http.NewRequest("POST", tc.url, bytes.NewBuffer(input))
		assert.NoError(err)
    req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		main.SMSHandler(w, req)

		assert.Equal(tc.expectedStatusCode, w.Code, tc.description)
		assert.Equal(tc.output, w.Body, tc.description)
	}
}
