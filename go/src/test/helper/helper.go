package helper

import (
	b "bytes"
	"encoding/json"
	"log"
	"net/http"
	u "net/url"
	"strconv"
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

func GetWithQueryParams(httpMethod string, url string, m map[string]interface{}) (*http.Response, error) {
	var sessionid string
	str, ok := m["sessionid"].(string)
	if ok && str != "" {
		sessionid = str
		delete(m, "sessionid")
	}

	if baseUrl, err := u.Parse(url); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		params := u.Values{}
		// Respect type, then add as string
		for key, val := range m {
			if _, ok := val.(int); ok {
				params.Add(key, strconv.Itoa(val.(int)))
			} else {
				params.Add(key, val.(string))
			}
		}
		baseUrl.RawQuery = params.Encode()

		queryUrl := baseUrl.String()

		request, err := http.NewRequest(httpMethod, queryUrl, nil)
		if err != nil {
			log.Fatal(err)
		}
		if ok && str != "" {
			request.Header.Add("Authentication", sessionid)
		}
		return http.DefaultClient.Do(request)
	}
}
