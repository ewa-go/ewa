package ewa

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gbrlsnchs/jwt"
)

func equal(t *testing.T, values ...interface{}) {
	var a, b interface{}
	for i, v := range values {
		if i%2 == 0 {
			a = v
		} else {
			b = v
			if a != b {
				t.Fatalf("Values did not match, a: %v, b: %v\n", a, b)
			}
		}
	}
}

func marshal(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

type SecurityUser struct {
	Login string `json:"login"`
	m     map[string]interface{}
}

func (u *SecurityUser) Username() string {
	return u.Login
}

func (u *SecurityUser) Get(name string) any {
	if v, ok := u.m[name]; ok {
		return v
	}
	return nil
}

func getAuthorization() Authorization {
	return Authorization{
		Unauthorized: func(err error) bool {
			if err != nil {
				return false
			}
			return true
		},
		Basic: &Basic{
			header: "Basic dXNlcjpRcTEyMzQ1Ng==",
			Handler: func(c *Context, user string, pass string) error {
				if user == "user" && pass == "Qq123456" {
					return nil
				}
				return errors.New("Unauthorized")
			},
		},
		ApiKey: &ApiKey{
			KeyName: "Token",
			Param:   ParamHeader,
			value:   "e6217493-0f45-447c-bd6d-7331544a9e5e",
			Handler: func(c *Context, token string) (username string, err error) {
				if token != "e6217493-0f45-447c-bd6d-7331544a9e5e" {
					return "", errors.New("invalid token")
				}
				return "apiKeyUser", nil
			},
		},
		BearerToken: &BearerToken{
			value: "Bearer e6217493-0f45-447c-bd6d-7331544a9e5e",
			Handler: func(c *Context, token string, isJWT bool) (username string, err error) {
				if !isJWT {
					if token != "e6217493-0f45-447c-bd6d-7331544a9e5e" {
						return "", errors.New("invalid token")
					}
					return "bearerUser", nil
				}
				jot, err := jwt.FromString(token)
				if err != nil {
					return "", errors.New("invalid token")
				}
				public := jot.Public()
				if public != nil {
					if value, ok := public["name"]; ok {
						return value.(string), nil
					}
				}
				return "", errors.New("username not found")
			},
		},
		OAuth2: &OAuth2{
			HeaderPrefix: "Bearer",
			Param:        ParamHeader,
			Flow: Flow{
				Type: FlowImplicit,
				Url:  "https://www.googleapis.com/auth/userinfo.email",
			},
			Handler: func(c *Context, token string) (username string, err error) {
				return "oAuth2User", nil
			},
			value: "",
		},
	}
}

func TestAuthorization_Definition(t *testing.T) {

	a := getAuthorization()

	equal(t, marshal(a.Get(BasicAuth).Definition()), `{"description":"Basic Authorization","type":"basic"}`)
	equal(t, marshal(a.Get(ApiKeyAuth).Definition()), `{"description":"Api Key Authorization. Set name: Token, in: header","type":"apiKey","name":"Token","in":"header"}`)
	equal(t, marshal(a.Get(BearerTokenAuth).Definition()), `{"description":"Bearer Token Authorization","type":"bearerToken"}`)
	equal(t, marshal(a.Get(OAuth2Auth).Definition()), `{"description":"OAuth2 Authorization","type":"oauth2","flow":"implicit","authorizationUrl":"https://www.googleapis.com/auth/userinfo.email"}`)
}

func TestAuthorization_ByHeader(t *testing.T) {

	a := getAuthorization()

	s := a.ByHeader("Basic dXNlcjpRcTEyMzQ1Ng==")
	if s == nil {
		t.Fatal("Authorization by header failed")
	}
	equal(t, s.Name(), BasicAuth)

	s = a.ByHeader("Bearer e6217493-0f45-447c-bd6d-7331544a9e5e")
	if s == nil {
		t.Fatal("Authorization by header failed")
	}
	equal(t, s.Name(), BearerTokenAuth)

	s = a.ByHeader(`Digest username="user", realm="test@mail.ru", nonce="Nonce", uri="/api/bridge/adMember", algorithm="MD5", qop=auth-int, nc=0000001, cnonce="0a4f113b", response="8faf1d6b11c89d99dea1b50b6f2cf68d", opaque="Opaque"`)
	if s == nil {
		t.Fatal("Authorization by header failed")
	}
	equal(t, s.Name(), DigestAuth)
}

/* Basic */
func TestBasicParse(t *testing.T) {
	b := Basic{
		header: "Basic dXNlcjpRcTEyMzQ1Ng==",
	}
	username, pass, ok := b.parse()
	if !ok {
		t.Fatal("Basic parse failed")
	}
	equal(t, username, "user", pass, "Qq123456")
}

func TestBasicHandler(t *testing.T) {
	c := new(Context)
	i, err := Do(getAuthorization().Get(BasicAuth), c)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, i.Username, "user")
}

/* BearerToken */
func TestBearerTokenParse(t *testing.T) {
	b := BearerToken{
		value: "Bearer e6217493-0f45-447c-bd6d-7331544a9e5e",
	}
	token, jwt, ok := b.parse()
	if !ok {
		t.Fatal("Bearer parse failed")
	}
	equal(t, token, "e6217493-0f45-447c-bd6d-7331544a9e5e", jwt, false)
	b.value = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	token, jwt, ok = b.parse()
	if !ok {
		t.Fatal("Bearer parse failed")
	}
	equal(t, token,
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		jwt, true)
}

func TestBearerTokenHandler(t *testing.T) {
	c := new(Context)
	i, err := Do(getAuthorization().Get(BearerTokenAuth, "Bearer e6217493-0f45-447c-bd6d-7331544a9e5e"), c)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, i.Username, "bearerUser")
	i, err = Do(getAuthorization().Get(BearerTokenAuth, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"), c)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, i.Username, "John Doe")
}

/* ApiKey */
func TestApiKeyHandler(t *testing.T) {
	c := new(Context)
	i, err := Do(getAuthorization().Get(ApiKeyAuth), c)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, i.Username, "apiKeyUser")
}
