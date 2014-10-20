package helper

import (
	b "bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Execute(httpMethod string, url string, m map[string]interface{}) (*http.Response, error) {
	var sessionid string
	str, ok := m["sessionid"].(string)
	if ok && str != "" {
		sessionid = str
		delete(m, "sessionid")
	}
	if bytes, err := json.Marshal(m); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		request, err := http.NewRequest(httpMethod, url, b.NewReader(bytes))
		if err != nil {
			log.Fatal(err)
		}
		if ok && str != "" {
			request.Header.Add("Authentication", sessionid)
		}
		return http.DefaultClient.Do(request)
	}
}
