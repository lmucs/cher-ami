package routes

import (
	"github.com/ant0ine/go-json-rest/rest"
	cheramiapi "github.com/rtoal/cher-ami/server/api"
)

func MakeHandler(api cheramiapi.Api, disableLogs bool) (rest.ResourceHandler, error) {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
		DisableLogger:            disableLogs,
	}

	err := handler.SetRoutes(
		&rest.Route{"POST", "/signup", api.Signup},
		&rest.Route{"POST", "/changepassword", api.ChangePassword},
		&rest.Route{"POST", "/sessions", api.Login},
		&rest.Route{"DELETE", "/sessions", api.Logout},
		//&rest.Route{"GET", "/users/:handle", api.GetUser},
		&rest.Route{"GET", "/users", api.SearchForUsers},
		&rest.Route{"DELETE", "/users/:handle", api.DeleteUser},
		&rest.Route{"GET", "/messages", api.GetAuthoredMessages},
		&rest.Route{"GET", "/messages/:id", api.GetMessageById},
		&rest.Route{"POST", "/messages", api.NewMessage},
		&rest.Route{"PATCH", "/messages/:id", api.EditMessage},
		&rest.Route{"DELETE", "/messages", api.DeleteMessage},
		&rest.Route{"POST", "/joindefault", api.JoinDefault},
		&rest.Route{"POST", "/join", api.Join},
		&rest.Route{"POST", "/block", api.BlockUser},
		&rest.Route{"POST", "/circles", api.NewCircle},
	)

	return handler, err
}
