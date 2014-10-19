package routes

import (
	cheramiapi "../api"
	"github.com/ant0ine/go-json-rest/rest"
)

func MakeHandler(a cheramiapi.Api) (rest.ResourceHandler, error) {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}

	err := handler.SetRoutes(
		&rest.Route{"POST", "/signup", a.Signup},
		&rest.Route{"POST", "/login", a.Login},
		&rest.Route{"POST", "/logout", a.Logout},
		&rest.Route{"POST", "/changepassword", a.ChangePassword},
		&rest.Route{"GET", "/users/user", a.GetUser},
		&rest.Route{"GET", "/users", a.GetUsers},
		&rest.Route{"DELETE", "/users/user", a.DeleteUser},
		&rest.Route{"GET", "/messages", a.GetAuthoredMessages},
		&rest.Route{"GET", "/messages/{author}", a.GetMessagesByHandle},
		&rest.Route{"POST", "/messages", a.NewMessage},
		&rest.Route{"DELETE", "/messages", a.DeleteMessage},
		&rest.Route{"POST", "/publish", a.PublishMessage},
		&rest.Route{"POST", "/joindefault", a.JoinDefault},
		&rest.Route{"POST", "/join", a.Join},
		&rest.Route{"POST", "/block", a.BlockUser},
		&rest.Route{"POST", "/circles", a.NewCircle},
	)

	return handler, err
}
