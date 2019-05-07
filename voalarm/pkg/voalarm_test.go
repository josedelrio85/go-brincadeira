package voalarm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")
		switch req.Method {
		case http.MethodGet:
			rw.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodPost:
			if contentType != "application/json" {
				rw.WriteHeader(http.StatusUnsupportedMediaType)
				break
			}

			var r Alarm
			data, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(data, &r)
			if (Alarm{}) == r {
				rw.WriteHeader(http.StatusInternalServerError)
			}

			rw.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	a := errors.New("test error")
	tests := []struct {
		Description    string
		ExpectedStatus int
		TypeRequest    string
		Contentype     string
		Params         Alarm
	}{
		{
			Description:    "Get method with no headers",
			ExpectedStatus: http.StatusMethodNotAllowed,
			TypeRequest:    http.MethodGet,
		},
		{
			Description:    "Post method with no content-type",
			ExpectedStatus: http.StatusUnsupportedMediaType,
			TypeRequest:    http.MethodPost,
		},
		{
			Description:    "Post method with content-type",
			ExpectedStatus: http.StatusOK,
			TypeRequest:    http.MethodPost,
			Contentype:     "application/json",
			Params: Alarm{
				MessageType:       Acknowledgement,
				EntityState:       Acknowledgement,
				EntityID:          "go! exception",
				EntityDisplayName: "go! exception",
				StateMessage:      a.Error(),
				StateStartTime:    time.Now().Format("2006-01-02 15:04:05"),
			},
		},
		{
			Description:    "Post method with content-type and malformed parameters",
			ExpectedStatus: http.StatusInternalServerError,
			TypeRequest:    http.MethodPost,
			Contentype:     "application/json",
			Params: Alarm{
				// MessageType:       Acknowledgement,
				EntityState:       Acknowledgement,
				EntityID:          "go! exception",
				EntityDisplayName: "go! exception",
				StateMessage:      "hola",
				StateStartTime:    time.Now().Format("2006-01-02 15:04:05"),
			},
		},
	}

	for _, test := range tests {

		bytevalues, berr := json.Marshal(test.Params)
		if berr != nil {
			t.Errorf("Error marshaling alarm %s", berr)
		}

		req, err := http.NewRequest(test.TypeRequest, server.URL, bytes.NewBuffer(bytevalues))
		req.Header.Add("Content-Type", test.Contentype)

		http := &http.Client{}
		resp, err := http.Do(req)
		if err != nil {
			t.Errorf("error sending test Request: Err: %v", err)
		}

		assert.NoError(err)
		assert.Equal(test.ExpectedStatus, resp.StatusCode)
	}
}
