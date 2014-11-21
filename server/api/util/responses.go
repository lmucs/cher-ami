package util

import (
	"../../types"
	"github.com/ant0ine/go-json-rest/rest"
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

func (u Util) SimpleJsonValidationReason(w rest.ResponseWriter, code int, err []error) {
	errorMessage := make([]string, len(err))
	for i := range err {
		errorMessage[i] = err[i].Error()
	}

	w.WriteHeader(code)
	w.WriteJson(types.Json{
		"reason": errorMessage,
	})
}

func (u Util) FailedToAuthenticate(w rest.ResponseWriter) {
	w.WriteHeader(401)
	w.WriteJson(types.Json{
		"response": "Failed to authenticate user request",
		"reason":   "Missing, illegal or expired token",
	})
}

func (u Util) FailedToDetermineHandleFromAuthToken(w rest.ResponseWriter) {
	w.WriteHeader(500)
	w.WriteJson(types.Json{
		"reason": "Unexpected failure to retrieve owner of Authentication token",
	})
}
