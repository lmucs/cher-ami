package helper

import (
	b "bytes"
	"encoding/json"
	"error"
	"log"
	"net/http"
	u "net/url"
	"os"
)

type ClosingBuffer struct {
	Buffer *b.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

func Execute(httpMethod string, url string, m map[string]interface{}) (*http.Response, error) {
	request := &http.Request{}

	// Pull sessionid and put it into the header
	var sessionid string
	if str, ok := m["sessionid"].(string); ok && str != "" {
		sessionid = str
		delete(m, "sessionid")
		request.Header.Set("Authentication", sessionid)
	}

	if bytes, err := json.Marshal(m); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		request.Method = httpMethod
		if request.URL, err = u.Parse(url); err != nil {
			return nil, err
		}
		request.Body = &ClosingBuffer{b.NewBuffer(bytes)}
		return http.DefaultClient.Do(request)
	}

}
