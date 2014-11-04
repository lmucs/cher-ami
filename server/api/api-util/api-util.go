package responses

import (
	"github.com/ant0ine/go-json-rest/rest"
)

//
// Responses package
// Used to send common Json responses
//

type Resp struct{}

func (r Resp) SimpleJsonResponse(w rest.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.WriteJson(map[string]interface{}{
		"Response": message,
	})
}

func (r Resp) FailedToAuthenticate(w rest.ResponseWriter) {
	w.WriteHeader(401)
	w.WriteJson(map[string]interface{}{
		"response": "Failed to authenticate user request",
		"reason":   "Missing, illegal or expired token",
	})
}
