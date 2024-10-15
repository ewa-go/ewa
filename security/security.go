package security

import (
	"strings"
)

const (
	NoAuth          = ""
	BasicAuth       = "Basic"
	DigestAuth      = "Digest"
	ApiKeyAuth      = "ApiKey"
	OAuth1Auth      = "OAuth1"
	OAuth2Auth      = "OAuth2"
	BearerTokenAuth = "Bearer"
	JWTBearerAuth   = "JWTBearer"
)

const (
	TypeBasic       = "basic"
	TypeApiKey      = "apiKey"
	TypeOAuth1      = "oauth1"
	TypeOAuth2      = "oauth2"
	TypeDigest      = "digest"
	TypeBearerToken = "bearerToken"
	//TypeJWTBearer   = "jwtBearer"
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
	//JWTBearer    *JWTBearer
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
	Name() string
}

const (
	ParamQuery  = "query"
	ParamHeader = "header"
)

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
		if len(values) > 0 {
			a.OAuth2.SetValue(values[0].(string))
		}
		return a.OAuth2
	case BearerTokenAuth:
		if len(values) > 0 {
			a.BearerToken.SetValue(values[0].(string))
		}
		return a.BearerToken
		/*case JWTBearerAuth:
		return a.JWTBearer*/
	}
	return nil
}

func (a Authorization) ByHeader(header string) IAuthorization {

	if len(header) == 0 {
		return nil
	}
	index := strings.Index(header, " ")
	// Проверка заголовка Authorization
	switch header[:index] {
	case BasicAuth:
		if a.Basic != nil {
			return a.Basic.SetHeader(header)
		}
	case BearerTokenAuth:
		// BearerToken
		if a.BearerToken != nil {
			return a.BearerToken.SetValue(header)
		}
	case DigestAuth:
		if a.Digest != nil {
			return a.Digest.SetHeader(header)
		}
	}

	return nil
}
