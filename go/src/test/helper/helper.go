package helper

import (
	b "bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	u "net/url"
	"strconv"
)

//
// Execute Requests:
//

func Execute(httpMethod string, url string, m map[string]interface{}) (*http.Response, error) {
	sessionid := ""
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
			request.Header.Add("Authorization", sessionid)
		}
		return http.DefaultClient.Do(request)
	}
}

func GetWithQueryParams(url string, m map[string]interface{}) (*http.Response, error) {
	sessionid := ""
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

		request, err := http.NewRequest("GET", queryUrl, nil)
		if err != nil {
			log.Fatal(err)
		}
		if ok && str != "" {
			request.Header.Add("Authorization", sessionid)
		}
		return http.DefaultClient.Do(request)
	}
}

//
// Read Body of Response:
//

func GetJsonResponseMessage(response *http.Response) string {
	type Json struct {
		Response string
	}

	var message Json

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
	}

	return message.Response
}

func GetJsonUserData(response *http.Response) string {
	type Json struct {
		Handle string
		Name   string
	}

	var user Json

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	return user.Handle
}

func Unmarshal(response *http.Response, v interface{}) {
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &v); err != nil {
		log.Fatal(err)
	}
}

//
// Read info from headers:
//

func GetSessionFromResponse(response *http.Response) string {
	authentication := struct {
		Response  string
		Sessionid string
	}{}
	var (
		body []byte
		err  error
	)
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &authentication); err != nil {
		log.Fatal(err)
	}

	return authentication.Sessionid
}
