package helper

import (
	b "bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Execute(httpMethod string, url string, m map[string]interface{}) (*http.Response, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest(httpMethod, url, b.NewReader(bytes))
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}
