package ewa

import (
	"github.com/ewa-go/ewa/consts"
	"github.com/ewa-go/ewa/security"
	"github.com/ewa-go/ewa/session"
	"testing"
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
		Authorization: security.Authorization{},
		Session: &session.Config{
			RedirectPath: "/login",
			SessionHandler: func(value string) (user string, err error) {
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
		Authorization: security.Authorization{},
		Session: &session.Config{
			RedirectPath: "/login",
			SessionHandler: func(value string) (user string, err error) {
				return "username", nil
			},
		},
		Permission: &Permission{
			AllRoutes: true,
			Handler: func(c *Context, identity *security.Identity, method, path string) bool {
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
