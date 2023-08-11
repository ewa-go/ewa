package security

const (
	NoAuth          = ""
	BasicAuth       = "Basic"
	DigestAuth      = "Digest"
	ApiKeyAuth      = "ApiKey"
	OAuth1Auth      = "OAuth1"
	OAuth2Auth      = "OAuth2"
	BearerTokenAuth = "BearerToken"
	JWTBearerAuth   = "JWTBearer"
)

const (
	TypeBasic       = "basic"
	TypeApiKey      = "apiKey"
	TypeOAuth1      = "oauth1"
	TypeOAuth2      = "oauth2"
	TypeDigest      = "digest"
	TypeBearerToken = "bearerToken"
	TypeJWTBearer   = "jwtBearer"
)

type Authorization struct {
	AllRoutes    string
	Unauthorized UnauthorizedHandler
	Basic        *Basic
	Digest       *Digest
	ApiKey       *ApiKey
	OAuth1       *OAuth1
	OAuth2       *OAuth2
	BearerToken  *BearerToken
	JWTBearer    *JWTBearer
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

func Do(a IAuthorization) (*Identity, error) {
	return a.Do()
}

func (a Authorization) Get(auth string, values ...interface{}) IAuthorization {
	switch auth {
	case BasicAuth:
		if len(values) > 0 {
			a.Basic.SetHeader(values[0].(string))
		}
		return a.Basic
	case ApiKeyAuth:
		if len(values) > 0 {
			a.ApiKey.SetValue(values[0].(string))
		}
		return a.ApiKey
	case DigestAuth:
		return a.Digest
	case OAuth1Auth:
		return a.OAuth1
	case OAuth2Auth:
		return a.OAuth2
	case BearerTokenAuth:
		return a.BearerToken
	case JWTBearerAuth:
		return a.JWTBearer
	}
	return nil
}
