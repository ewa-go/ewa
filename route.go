package ewa

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ewa-go/ewa/v2/consts"
)

type Route struct {
	emptyPathParam *EmptyPathParam
	session        SessionTurn
	isPermission   bool
	models         Models
	Handler        Handler
	Operation
}

type EmptyPathParam struct {
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`
}

// Map тип список
type Map map[string]interface{}

type SessionTurn int

const (
	None SessionTurn = iota
	Is
	On
	Off
)

// setResponse описываем варианты ответов для Swagger
func (e *EmptyPathParam) setResponse(code int, modelName string, isArray bool, headers Headers, desc ...string) {
	response := &Response{
		Headers: headers,
		Schema:  NewSchema(modelName, isArray),
	}
	if desc != nil {
		response.Description = desc[0]
	}
	e.Responses[strconv.Itoa(code)] = response
}

// SetResponse описываем варианты ответов для Swagger
func (e *EmptyPathParam) SetResponse(code int, modelName string, headers Headers, desc ...string) *EmptyPathParam {
	e.setResponse(code, modelName, false, headers, desc...)
	return e
}

// SetResponseArray описываем варианты ответов для Swagger
func (e *EmptyPathParam) SetResponseArray(code int, modelName string, headers Headers, desc ...string) *EmptyPathParam {
	e.setResponse(code, modelName, true, headers, desc...)
	return e
}

// SetEmptyParam указываем параметры маршрута
func (r *Route) SetEmptyParam(summary string, desc ...string) *EmptyPathParam {
	e := &EmptyPathParam{
		Summary:   summary,
		Responses: map[string]*Response{},
	}
	if desc != nil {
		e.Description = desc[0]
	}
	r.emptyPathParam = e
	return r.emptyPathParam
}

// SetParameters указываем параметры маршрута
func (r *Route) SetParameters(params ...*Parameter) *Route {
	for _, param := range params {
		r.Parameters = append(r.Parameters, param)
	}
	return r
}

// InitParametersByModel Формирование параметров на основе модели
func (r *Route) InitParametersByModel(name string) *Route {
	r.SetParameters(ModelToParameters(r.Model(name))...)
	return r
}

// Model Вернуть модель параметров
func (r *Route) Model(name string) interface{} {
	if model, ok := r.models[name]; ok {
		return model
	}
	return nil
}

// SetConsumes устанавливаем Content-Type запроса для Swagger
func (r *Route) SetConsumes(c ...string) *Route {
	r.Consumes = c
	return r
}

// SetProduces устанавливаем Content-Type ответа для Swagger
func (r *Route) SetProduces(p ...string) *Route {
	r.Produces = p
	return r
}

// SetOperationID устанавливаем идентификатор операции для Swagger
func (r *Route) SetOperationID(id string) *Route {
	r.ID = id
	return r
}

// setResponse описываем варианты ответов для Swagger
func (r *Route) setResponse(code int, modelName string, isArray bool, headers Headers, desc ...string) {
	response := &Response{
		Headers: headers,
		Schema:  NewSchema(modelName, isArray),
	}
	if desc != nil {
		response.Description = desc[0]
	}
	r.Responses[strconv.Itoa(code)] = response
}

// SetDefaultResponse описываем варианты ответов для Swagger
func (r *Route) SetDefaultResponse(modelName string, isArray bool, headers Headers, desc ...string) *Route {
	response := &Response{
		Schema:  NewSchema(modelName, isArray),
		Headers: headers,
	}
	if desc != nil {
		response.Description = desc[0]
	}
	r.Responses["default"] = response
	return r
}

// SetResponse описываем варианты ответов для Swagger
func (r *Route) SetResponse(code int, modelName string, headers Headers, desc ...string) *Route {
	r.setResponse(code, modelName, false, headers, desc...)
	return r
}

// SetResponseArray описываем варианты ответов для Swagger
func (r *Route) SetResponseArray(code int, modelName string, headers Headers, desc ...string) *Route {
	r.setResponse(code, modelName, true, headers, desc...)
	return r
}

// SetDescription описание операции
func (r *Route) SetDescription(desc string) *Route {
	r.Description = desc
	return r
}

// SetSummary резюме запроса
func (r *Route) SetSummary(s string) *Route {
	r.Summary = s
	return r
}

// SetSign устанавливаем вариант входа/выхода для маршрута
/*func (r *Route) SetSign(sign Sign) *Route {
	r.sign = sign
	return r
}*/

// SetSecurity указываем метод авторизации
func (r *Route) SetSecurity(security ...string) *Route {
	for _, sec := range security {
		r.Security = append(r.Security, map[string][]string{
			sec: {},
		})
	}
	//r.auth = auth
	return r
}

// BrowserSession вешаем получение аутентификации сессии,
func (r *Route) Session(t ...SessionTurn) *Route {
	if t == nil {
		r.session = Is
	} else {
		r.session = t[0]
	}
	return r
}

// Permission ставим флаг для проверки маршрута на право доступа
func (r *Route) Permission(b bool) *Route {
	r.isPermission = b
	return r
}

// EmptyHandler пустой обработчик
func (r *Route) EmptyHandler() {
	r.Handler = nil
}

// SetHandler устанавливаем обработчик
func (r *Route) SetHandler(handler Handler) *Route {
	r.Handler = handler
	return r
}

type ISession interface {
	GetSession(string) ISession
}

// getHandler возвращаем обработчик основанный на параметрах конфигурации маршрута
func (r *Route) getHandler(config Config, swagger *Swagger) Handler {

	return func(c *Context) (err error) {

		c.Swagger = *swagger

		var auth IAuthorization
		if len(r.Security) > 0 {
			auth = config.Authorization.ByHeader(c.Get(consts.HeaderAuthorization))
		}
		for _, sec := range r.Security {
			// Пытаемся найти авторизацию среди установок
			if auth != nil {
				if _, ok := sec[auth.Name()]; ok {
					break
				}
			}
			// BearerToken получение из query параметра
			if _, ok := sec[BearerTokenAuth]; ok {
				if config.Authorization.BearerToken != nil && config.Authorization.BearerToken.Param == ParamQuery {
					auth = config.Authorization.Get(BearerTokenAuth, c.QueryParam("token"))
					break
				}
			}
			// ApiKey получение токена
			if _, ok := sec[ApiKeyAuth]; ok {
				if config.Authorization.ApiKey != nil {
					var value string
					apiKey := config.Authorization.ApiKey
					switch apiKey.Param {
					// Если не нашли в заголовке, то ищем в переменных запроса адресной строки
					case ParamQuery:
						value = c.QueryParam(apiKey.KeyName)
					// Пытаемся получить из заголовка токен
					case ParamHeader:
						value = c.Get(apiKey.KeyName)
					}
					if len(value) > 0 {
						auth = config.Authorization.Get(ApiKeyAuth, value)
						break
					}
				}
			}
		}

		if auth != nil {
			// Получаем пользователя, если нет ошибок, то выходим
			c.Identity, err = Do(auth, c)
			if err != nil {
				switch auth.Name() {
				case BasicAuth:
					c.Set(consts.HeaderWWWAuthenticate, err.Error())
				}
			}
		} else if len(r.Security) > 0 {
			err = errors.New("unauthorized")
		}

		// Проверка на сессию
		if config.Session != nil && r.session != None {
			keyName := config.Session.KeyName
			value := c.Cookies(keyName)
			if len(value) > 0 {
				c.Identity, err = config.Session.Check(c, value)
				if err != nil {
					c.ClearCookie(keyName)
				} else {
					c.Session = newSession(keyName, value)
				}
			}

			switch r.session {
			case Is:
				if c.Session == nil {
					return c.Redirect(config.Session.RedirectPath, config.Session.RedirectStatus)
				}
			case On:
				value = config.Session.GenSessionIdHandler()
				cookie := &http.Cookie{
					Name:    keyName,
					Value:   value,
					Expires: time.Now().Add(config.Session.Expires),
				}
				c.SetCookie(cookie)
				c.Session = newSession(keyName, value)
			case Off:
				ok := true
				if c.Session != nil && config.Session.DeleteSessionHandler != nil {
					ok = config.Session.DeleteSessionHandler(c.Session.Value)
					c.ClearCookie(config.Session.KeyName)
					c.Session = nil
				}
				if ok {
					return c.Redirect(config.Session.RedirectPath, config.Session.RedirectStatus)
				}
			}
		}

		// Проверка на ошибку авторизации и отправку кода 401
		if r.session == None && err != nil {
			if config.Authorization.Unauthorized != nil && config.Authorization.Unauthorized(err) {
				return c.SendString(consts.StatusUnauthorized, err.Error())
			}
			return c.SendStatus(consts.StatusUnauthorized)
		}

		// Доступ к маршрутам
		if r.isPermission && config.Permission != nil {
			identity := c.Identity
			if config.Permission.Handler != nil && !config.Permission.Handler(c, identity, c.Method(), c.Path()) {
				if config.Permission.NotPermissionHandler != nil {
					return config.Permission.NotPermissionHandler(c, consts.StatusForbidden, "Forbidden")
				}
				return c.SendStatus(consts.StatusForbidden)
			}
			c.Identity = identity
		}

		// Обычный маршрут
		return r.Handler(c)
	}
}

func newSession(keyName, value string) *BrowserSession {
	return &BrowserSession{
		Key:      keyName,
		Value:    value,
		LastTime: time.Now(),
	}
}
