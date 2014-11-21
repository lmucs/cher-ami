package routes

import (
	cheramiapi "../api"
	"github.com/ant0ine/go-json-rest/rest"
)

func MakeHandler(api cheramiapi.Api, disableLogs bool) (rest.ResourceHandler, error) {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		DisableLogger:            disableLogs,
	}

	err := handler.SetRoutes(
		&rest.Route{"POST", "/signup", api.Signup},
		&rest.Route{"POST", "/sessions", api.Login},
		&rest.Route{"DELETE", "/sessions", api.Logout},
		&rest.Route{"GET", "/users/:handle", api.GetUser},
		&rest.Route{"PUT", "/users/:handle", api.SetUser},
		&rest.Route{"GET", "/users", api.SearchForUsers},
		&rest.Route{"GET", "/messages", api.GetAuthoredMessages},
		&rest.Route{"GET", "/messages/:id", api.GetMessageById},
		&rest.Route{"POST", "/messages", api.NewMessage},
		&rest.Route{"PATCH", "/messages/:id", api.EditMessage},
		&rest.Route{"DELETE", "/messages", api.DeleteMessage},
		&rest.Route{"POST", "/joindefault", api.JoinDefault},
		&rest.Route{"POST", "/join", api.Join},
		&rest.Route{"POST", "/block", api.BlockUser},
		&rest.Route{"POST", "/circles", api.NewCircle},
		&rest.Route{"GET", "/circles", api.SearchCircles},
	)

	return handler, err
}
