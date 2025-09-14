package ewa

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
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
	Do(c *Context) (*Identity, error)
	Definition() Definition
	Name() string
}

const (
	ParamQuery  = "query"
	ParamHeader = "header"
)

func Do(a IAuthorization, c *Context) (*Identity, error) {
	return a.Do(c)
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

// Basic авторизация
type Basic struct {
	header  string
	Handler BasicAuthHandler
}

type BasicAuthHandler func(c *Context, user string, pass string) error

func (b *Basic) parse() (username, password string, ok bool) {
	const prefix = "Basic "
	if len(b.header) < len(prefix) || !strings.EqualFold(b.header[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(b.header[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	i := strings.IndexByte(cs, ':')
	if i < 0 {
		return
	}
	return cs[:i], cs[i+1:], true
}

func (b *Basic) SetHeader(header string) *Basic {
	b.header = header
	return b
}

func (b *Basic) Name() string {
	return BasicAuth
}

func (b *Basic) Do(c *Context) (*Identity, error) {

	err := errors.New(`basic realm="Необходимо указать имя пользователя и пароль"`)
	if b.header == "" {
		return nil, err
	}

	username, password, ok := b.parse()
	if !ok {
		return nil, err
	}

	if b.Handler == nil {
		return nil, fmt.Errorf("[%s] handler not initialized]", BasicAuth)
	}
	err = b.Handler(c, username, password)
	if err != nil {
		return nil, err
	}

	identity := &Identity{
		Username: username,
		AuthName: BasicAuth,
		Datetime: time.Now(),
	}

	return identity, nil
}

func (b *Basic) Definition() Definition {
	return Definition{
		Type:        TypeBasic,
		Description: "Basic Authorization",
	}
}

// ApiKey авторизация
type ApiKey struct {
	KeyName string
	Param   string
	value   string
	Handler ApiKeyAuthHandler
}

type ApiKeyAuthHandler func(c *Context, token string) (username string, err error)

func (a *ApiKey) SetValue(value string) *ApiKey {
	a.value = value
	return a
}

func (a *ApiKey) Name() string {
	return ApiKeyAuth
}

func (a *ApiKey) Do(c *Context) (identity *Identity, err error) {

	if a.value == "" {
		return nil, errors.New(fmt.Sprintf("Not found token by [%s]", a.Param))
	}

	username := ""
	if a.Handler != nil {
		username, err = a.Handler(c, a.value)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("[%s] handler not initialized]", ApiKeyAuth)
	}

	identity = &Identity{
		Username: username,
		AuthName: ApiKeyAuth,
		Datetime: time.Now(),
	}

	return
}

func (a *ApiKey) Definition() Definition {
	return Definition{
		Type:        TypeApiKey,
		In:          a.Param,
		Name:        a.KeyName,
		Description: fmt.Sprintf("Api Key Authorization. Set name: %s, in: %s", a.KeyName, a.Param),
	}
}

// Digest авторизация
type Digest struct {
	Handler DigestAuthHandler

	header string
}

type Advanced struct {
	Realm       string
	Nonce       string
	Algorithm   string
	Qop         string
	NonceCount  string
	ClientNonce string
	Opaque      string
}

type DigestAuthHandler func(c *Context, user string, pass string, advanced Advanced) bool

const (
	Md5Algorithm           = "MD5"
	Md5SessAlgorithm       = "MD5-sess"
	Sha256Algorithm        = "SHA-256"
	Sha256SessAlgorithm    = "SHA-256-sess"
	Sha512256Algorithm     = "SHA-512-256"
	Sha512256SessAlgorithm = "SHA-512-256-sess"
)

func (d *Digest) Name() string {
	return DigestAuth
}

func (d *Digest) Do(c *Context) (identity *Identity, err error) {
	return
}

func (d *Digest) Definition() Definition {
	return Definition{
		Type:        TypeDigest,
		Description: fmt.Sprintf("Digest Authorization"),
	}
}

func (d *Digest) SetHeader(header string) *Digest {
	d.header = header
	return d
}

// BearerToken авторизация
type BearerToken struct {
	Param   string
	Handler BearerTokenHandler

	value string
}

type BearerTokenHandler func(c *Context, token string, isJWT bool) (username string, err error)

func (b *BearerToken) Name() string {
	return BearerTokenAuth
}

func (b *BearerToken) Do(c *Context) (identity *Identity, err error) {

	if b.value == "" {
		return nil, errors.New("header is required")
	}

	token, jwt, ok := b.parse()
	if !ok {
		return nil, errors.New("invalid token")
	}

	var username string
	if b.Handler != nil {
		username, err = b.Handler(c, token, jwt)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("[%s] handler not initialized]", BearerTokenAuth)
	}

	identity = &Identity{
		Username: username,
		AuthName: BearerTokenAuth,
		Datetime: time.Now(),
	}

	return
}

func (b *BearerToken) Definition() Definition {
	return Definition{
		Type:        TypeBearerToken,
		Description: "Bearer Token Authorization",
	}
}

func (b *BearerToken) SetValue(value string) *BearerToken {
	b.value = value
	return b
}

func (b *BearerToken) parse() (token string, isJWT, ok bool) {
	const prefix = "Bearer "
	if b.value[:len(prefix)] == prefix {
		token = b.value[len(prefix):]
		ok = true
		i := strings.Index(b.value, ".")
		if i > 0 {
			isJWT = true
		}
	}
	return
}

// OAuth1 авторизация
type OAuth1 struct{}

func (OAuth1) Name() string {
	return OAuth1Auth
}

func (OAuth1) Do(c *Context) (identity *Identity, err error) {
	return
}

func (OAuth1) Definition() Definition {
	return Definition{}
}

// OAuth2 авторизация
type OAuth2 struct {
	HeaderPrefix string
	Param        string
	Flow         Flow
	Handler      OAuth2AuthHandler

	value string
}

type Flow struct {
	Type   string
	Url    string
	Scopes map[string]string
}

const (
	FlowImplicit    = "implicit"
	FlowPassword    = "password"
	FlowAccessCode  = "access_code"
	FlowApplication = "application"
)

type OAuth2AuthHandler func(c *Context, token string) (username string, err error)

func (o *OAuth2) SetValue(value string) {
	o.value = value
}

func (o *OAuth2) Name() string {
	return OAuth2Auth
}

func (o *OAuth2) Do(c *Context) (identity *Identity, err error) {

	if o.value == "" {
		return nil, errors.New(fmt.Sprintf("Not found token by [%s]", o.Param))
	}

	username := ""
	if o.Handler != nil {
		username, err = o.Handler(c, o.value)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("[%s] handler not initialized]", OAuth2Auth)
	}

	identity = &Identity{
		Username: username,
		AuthName: OAuth2Auth,
		Datetime: time.Now(),
	}

	return
}

func (o *OAuth2) Definition() (d Definition) {
	d.Type = TypeOAuth2
	d.Description = "OAuth2 Authorization"
	d.Flow = o.Flow.Type
	d.Scopes = o.Flow.Scopes
	switch o.Flow.Type {
	case FlowImplicit, FlowAccessCode:
		d.AuthorizationURL = o.Flow.Url
	case FlowPassword, FlowApplication:
		d.TokenURL = o.Flow.Url
	}
	return
}

func (o *OAuth2) parse() (token string, ok bool) {
	switch o.Param {
	case ParamHeader:
		if o.value[:len(o.HeaderPrefix)] == o.HeaderPrefix {
			return o.value[len(o.HeaderPrefix):], true
		}
	case ParamQuery:
		return o.value, true
	}
	return
}

// JWTBearer авторизация
/*type JWTBearer struct {
	payload map[string]string
	Param   string
	value   string
	Handler JWTBearerAuthHandler
}

type JWTBearerAuthHandler func(jot *jwt.JWT) (username string, err error)

func (j *JWTBearer) Name() string {
	return JWTBearerAuth
}

func (j *JWTBearer) Do() (identity *Identity, err error) {

	if j.value == "" {
		return nil, errors.New(fmt.Sprintf("Not found token by [%s]", j.Param))
	}

	jot, ok := j.parse()
	if !ok {
		return nil, errors.New("invalid token")
	}

	username := ""
	if j.Handler != nil {
		username, err = j.Handler(jot)
		if err != nil {
			return nil, err
		}
	}

	identity = &Identity{
		Username: username,
		AuthName: JWTBearerAuth,
		Datetime: time.Now(),
	}

	return
}

func (j *JWTBearer) Definition() Definition {
	return Definition{
		Type:        TypeJWTBearer,
		In:          j.Param,
		Description: fmt.Sprintf("JWTBearer Authorization"),
	}
}

func (j *JWTBearer) SetValue(value string) *JWTBearer {
	j.value = value
	return j
}

func (j *JWTBearer) parse() (*jwt.JWT, bool) {
	const prefix = "Bearer "
	if j.value[:len(prefix)] == prefix {
		jot, err := jwt.FromString(j.value[len(prefix):])
		if err != nil {
			return nil, false
		}
		return jot, true
	}
	return nil, false
}*/
