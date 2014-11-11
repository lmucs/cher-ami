package api_util

import (
	"github.com/ant0ine/go-json-rest/rest"
)

//
// Responses package
// Used to send common Json responses
//

type Util struct{}
type json map[string]interface{}

func (u Util) SimpleJsonResponse(w rest.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.WriteJson(map[string]interface{}{
		"Response": message,
	})
}

func (u Util) FailedToAuthenticate(w rest.ResponseWriter) {
	w.WriteHeader(401)
	w.WriteJson(map[string]interface{}{
		"response": "Failed to authenticate user request",
		"reason":   "Missing, illegal or expired token",
	})
}

func (u Util) FailedToDetermineHandleFromSession(w rest.ResponseWriter) {
	w.WriteHeader(400)
	w.WriteJson(json{
		"Response": "Unexpected failure to retrieve owner of session",
	})
	return
}
