package helper

import (
	"../../types"
	b "bytes"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	u "net/url"
	"strconv"
	"strings"
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

type fields struct {
	Response string `json:"response"`
	Reason   string `json:"reason"`
	Handle   string `json:"handle"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	Id       string `json:"id"`
	Url      string `json:"url"`
}

func Unmarshal(r *http.Response, v interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &v); err != nil {
		log.Fatal(err)
	}
}

func GetJsonResponseMessage(r *http.Response) string {
	f := fields{}
	Unmarshal(r, &f)
	return f.Response
}

func GetJsonReasonMessage(r *http.Response) string {
	f := fields{}
	Unmarshal(r, &f)
	return f.Reason
}

func GetJsonValidationReasonMessage(r *http.Response) []string {
	var message struct {
		Reason []string `json:"reason"`
	}
	Unmarshal(r, &message)
	return message.Reason
}

func GetJsonPatchValidationReasonMessage(r *http.Response) ([]string, int) {
	var message struct {
		Reason []string `json:"reason"`
		Index  int      `json:"index"`
	}
	Unmarshal(r, &message)
	return message.Reason, message.Index
}

func GetJsonUserData(r *http.Response) string {
	f := fields{}
	Unmarshal(r, &f)
	return f.Handle
}

func GetUrlFromResponse(r *http.Response) string {
	f := fields{}
	Unmarshal(r, &f)
	return f.Url
}

func GetIdFromUrlField(r *http.Response) string {
	return GetIdFromUrlString(GetUrlFromResponse(r))
}

func GetIdFromUrlString(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}

func GetAuthTokenFromResponse(r *http.Response) string {
	f := fields{}
	Unmarshal(r, &f)
	return f.Token
}
