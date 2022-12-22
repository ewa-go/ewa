package ewa

import (
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
