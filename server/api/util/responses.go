package util

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/rtoal/cher-ami/server/types"
)

//
// Responses package
// Used to send common Json responses
//

type Util struct{}

func (u Util) SimpleJsonResponse(w rest.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.WriteJson(types.Json{
		"response": message,
	})
}

func (u Util) SimpleJsonReason(w rest.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.WriteJson(types.Json{
		"reason": message,
	})
}

func (u Util) FailedToAuthenticate(w rest.ResponseWriter) {
	w.WriteHeader(401)
	w.WriteJson(types.Json{
		"response": "Failed to authenticate user request",
		"reason":   "Missing, illegal or expired token",
	})
}

func (u Util) FailedToDetermineHandleFromSession(w rest.ResponseWriter) {
	w.WriteHeader(400)
	w.WriteJson(types.Json{
		"response": "Unexpected failure to retrieve owner of session",
	})
	return
}
