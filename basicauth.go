package egowebapi

import "github.com/gofiber/fiber/v2/middleware/basicauth"

type BasicAuth basicauth.Config

type Authorizer func(user string, pass string) bool

func NewBasicAuth(users map[string]string, authorizer Authorizer, unauthorized Handler) BasicAuth {
	return BasicAuth{
		Users:        users,
		Realm:        "Forbidden",
		Authorizer:   authorizer,
		Unauthorized: unauthorized,
	}
}
