package security

const (
	NoAuth     = ""
	BasicAuth  = "Basic"
	DigestAuth = "Digest"
	ApiKeyAuth = "ApiKey"
	OAuth2Auth = "OAuth2"
)

const (
	TypeBasic  = "basic"
	TypeApiKey = "apiKey"
	TypeOAuth2 = "oauth2"
)

type Authorization struct {
	AllRoutes    string
	Unauthorized UnauthorizedHandler
	Basic        *Basic
	Digest       *Digest
	ApiKey       *ApiKey
	OAuth2       *OAuth2
}

type Definition struct {
	Description      string            `json:"description,omitempty"`
	Type             string            `json:"type"`
	Name             string            `json:"name,omitempty"`             // api key
	In               string            `json:"in,omitempty"`               // api key
	Flow             string            `json:"flow,omitempty"`             // oauth2
	AuthorizationURL string            `json:"authorizationUrl,omitempty"` // oauth2
	TokenURL         string            `json:"tokenUrl,omitempty"`         // oauth2
	Scopes           map[string]string `json:"scopes,omitempty"`           // oauth2
}

type UnauthorizedHandler func(err error) bool

type IAuthorization interface {
	Do() (*Identity, error)
	Definition() Definition
}

func (a Authorization) Get(auth string) IAuthorization {
	switch auth {
	case BasicAuth:
		return a.Basic
	case ApiKeyAuth:
		return a.ApiKey
	case DigestAuth:
		return a.Digest
	case OAuth2Auth:
		return a.OAuth2
	}
	return nil
}
