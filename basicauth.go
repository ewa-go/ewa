package egowebapi

type BasicAuth struct {
	Authorizer   Authorizer
	Unauthorized Handler
}

func NewBasicAuth(authorizer Authorizer, unauthorized Handler) *BasicAuth {
	return &BasicAuth{
		Authorizer:   authorizer,
		Unauthorized: unauthorized,
	}
}
