package helper

import (
	Bytes "bytes"
	"encoding/json"
	"log"
	"net/http"
)

func MakeRequest(payload map[string]interface{}, url string) *http.Request {
	bytes, err := json.Marshal(payload)
	if err != nil {
		panic("Error marshalling payload")
	}

	reader := Bytes.NewReader(bytes)

	if request, err := http.NewRequest("POST", url, reader); err != nil {
		log.Fatal(err)
		return nil
	} else {
		return request
	}
}
