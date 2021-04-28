package egowebapi

type BasicAuth struct {
	Authorizer   Authorizer
	Unauthorized Handler
}

type Authorizer func(user string, pass string) bool

func NewBasicAuth(authorizer Authorizer, unauthorized Handler) *BasicAuth {
	return &BasicAuth{
		Authorizer:   authorizer,
		Unauthorized: unauthorized,
	}
}
