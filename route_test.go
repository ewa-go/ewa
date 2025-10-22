package ewa

import (
	"testing"

	"github.com/ewa-go/ewa/v2/consts"
)

func TestRoute_Session(t *testing.T) {

	route := Route{
		session: On,
		Handler: func(c *Context) error {
			return nil
		},
	}

	_ = route.getHandler(Config{
		Port:          0,
		Secure:        nil,
		Authorization: Authorization{},
		Session: &Session{
			RedirectPath: "/login",
			SessionHandler: func(c *Context, value string) (user string, err error) {
				return "username", nil
			},
		},
		Permission:     nil,
		Static:         nil,
		NotFoundPage:   "",
		Views:          nil,
		ContextHandler: nil,
		ErrorHandler:   nil,
	}, nil)

}

func TestRoute_Permission(t *testing.T) {

	route := Route{
		session: On,
		Handler: func(c *Context) error {
			return nil
		},
	}

	_ = route.getHandler(Config{
		Port:          0,
		Secure:        nil,
		Authorization: Authorization{},
		Session: &Session{
			RedirectPath: "/login",
			SessionHandler: func(c *Context, value string) (user string, err error) {
				return "username", nil
			},
		},
		Permission: &Permission{
			AllRoutes: true,
			Handler: func(c *Context, identity *Identity, method, path string) bool {
				if identity != nil && identity.Username == "username" {
					switch method {
					// ReadOnly
					case consts.MethodPost, consts.MethodPut, consts.MethodDelete:
						return false
					}
				}
				return true
			},
			NotPermissionHandler: nil,
		},
		Static:         nil,
		NotFoundPage:   "",
		Views:          nil,
		ContextHandler: nil,
		ErrorHandler:   nil,
	}, nil)
}
