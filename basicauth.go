package egowebapi

import "github.com/gofiber/fiber/v2/middleware/basicauth"

type BasicAuth struct {
	basicauth.Config
}

type Authorizer func(user string, pass string) bool

func NewBasicAuth(users map[string]string, authorizer Authorizer, unauthorized Handler) BasicAuth {
	return BasicAuth{
		Config: basicauth.Config{
			Users:        users,
			Realm:        "Forbidden",
			Authorizer:   authorizer,
			Unauthorized: unauthorized,
		},
	}
}
