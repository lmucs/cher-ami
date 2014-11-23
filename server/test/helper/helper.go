package helper

import (
	"../../types"
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
	token := ""
	if str, ok := m["token"].(string); ok && str != "" {
		token = str
		delete(m, "token")
	}

	if bytes, err := json.Marshal(m); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		request, err := http.NewRequest(httpMethod, url, b.NewReader(bytes))
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Authorization", token)
		return http.DefaultClient.Do(request)
	}
}

func ExecutePatch(token string, url string, m types.JsonArray) (*http.Response, error) {
	if bytes, err := json.Marshal(m); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		request, err := http.NewRequest("PATCH", url, b.NewReader(bytes))
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Authorization", token)
		return http.DefaultClient.Do(request)
	}
}

func GetWithQueryParams(url string, m map[string]interface{}) (*http.Response, error) {
	token := ""
	str, ok := m["token"].(string)
	if ok && str != "" {
		token = str
		delete(m, "token")
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
		request.Header.Add("Authorization", token)
		return http.DefaultClient.Do(request)
	}
}

//
// Read Body of Response:
//

func GetJsonResponseMessage(response *http.Response) string {
	var message struct {
		Response string
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
	}

	return message.Response
}

func GetJsonReasonMessage(response *http.Response) string {
	var message struct {
		Reason string
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
	}

	return message.Reason
}

func GetJsonValidationReasonMessage(response *http.Response) []string {
	var message struct {
		Reason []string
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
	}

	return message.Reason
}

func GetJsonPatchValidationReasonMessage(response *http.Response) ([]string, int) {
	var message struct {
		Reason []string
		Index  int
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
	}

	return message.Reason, message.Index
}

func GetJsonUserData(response *http.Response) string {
	var user struct {
		Handle string
		Name   string
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &user); err != nil {
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

func GetAuthTokenFromResponse(response *http.Response) string {
	var authentication struct {
		Response string
		Token    string
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &authentication); err != nil {
		log.Fatal(err)
	}

	return authentication.Token
}

func GetIdFromResponse(response *http.Response) string {
	var res struct {
		Id string
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &res); err != nil {
		log.Fatal(err)
	}

	return res.Id
}
