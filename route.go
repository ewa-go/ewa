package egowebapi

import (
	"github.com/egovorukhin/egowebapi/security"
	"strconv"
)

type Route struct {
	isEmptyParam bool
	Handler      Handler
	isSession    bool
	isPermission bool
	sign         Sign
	Operation
}

// Map тип список
type Map map[string]interface{}

type Sign int

const (
	SignNone Sign = iota
	SignIn
	SignOut
)

// SetParameters указываем параметры маршрута
func (r *Route) SetParameters(isEmptyParam bool, params ...*Parameter) *Route {
	r.isEmptyParam = isEmptyParam
	r.Parameters = params
	return r
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

// SetDefaultResponse описываем варианты ответов для Swagger
func (r *Route) SetDefaultResponse(resp Response) *Route {
	r.Responses["default"] = resp
	return r
}

// SetResponse описываем варианты ответов для Swagger
func (r *Route) SetResponse(code int, resp Response) *Route {
	r.Responses[strconv.Itoa(code)] = resp
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
func (r *Route) SetSign(sign Sign) *Route {
	r.sign = sign
	return r
}

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

// Session вешаем получение аутентификации сессии,
func (r *Route) Session() *Route {
	r.isSession = true
	return r
}

// Permission ставим флаг для проверки маршрута на право доступа
func (r *Route) Permission() *Route {
	r.isPermission = true
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

// getHandler возвращаем обработчик основанный на параметрах конфигурации маршрута
func (r *Route) getHandler(config Config, view *View, swagger Swagger) Handler {

	return func(c *Context) error {

		c.View = view
		c.Swagger = swagger

		var (
			err    error
			isAuth bool
		)
		if len(r.Security) > 0 {
			isAuth = true
		}
		for _, sec := range r.Security {
			for key := range sec {
				switch key {
				case security.BasicAuth:
					if config.Authorization.Basic != nil {
						config.Authorization.Basic.SetHeader(c.Get(HeaderAuthorization))
						c.Identity, err = config.Authorization.Basic.Do()
						if err != nil {
							c.Set(HeaderWWWAuthenticate, err.Error())
						}
					}
					break
				case security.DigestAuth:
					if config.Authorization.Digest != nil {
						c.Identity, err = config.Authorization.Digest.Do()
					}
					break
				case security.ApiKeyAuth:
					if config.Authorization.ApiKey != nil {
						a := config.Authorization.ApiKey
						var value string
						switch a.Param {
						// Пытаемся получить из заголовка токен
						case security.ParamQuery:
							value = c.QueryParam(a.KeyName)
							break
						// Если не нашли в заголовке, то ищем в переменных запроса адресной строки
						case security.ParamHeader:
							value = c.Get(a.KeyName)
							break
						}
						a.SetValue(value)
						c.Identity, err = a.Do()
					}
					break
				}
			}
		}

		// Проверка на сессию
		if config.Session != nil {
			// Вход/Выход из сессии
			switch r.sign {
			case SignNone:
				break
			case SignIn:
				config.Session.signIn(c)
				break
			case SignOut:
				config.Session.signOut(c)
				return c.Redirect(config.Session.RedirectPath, config.Session.RedirectStatus)
			}
			if !isAuth && r.isSession {
				err = config.Session.check(c)
			}
		}

		// Проверка на ошибку авторизации и отправку кода 401
		if err != nil {
			if isAuth {
				if config.Authorization.Unauthorized != nil && config.Authorization.Unauthorized(err) {
					return c.SendString(StatusUnauthorized, err.Error())
				}
				return c.SendStatus(StatusUnauthorized)
			} else if r.isSession {
				// Если cookie не существует, то перенаправляем запрос условно на "/login"
				return c.Redirect(config.Session.RedirectPath, config.Session.RedirectStatus)
			}
		}

		if config.Session != nil {
			// Вход/Выход из сессии
			switch r.sign {
			case SignNone:
				break
			case SignIn:
				config.Session.signIn(c)
				break
			case SignOut:
				config.Session.signOut(c)
				return c.Redirect(config.Session.RedirectPath, config.Session.RedirectStatus)
			}
			if r.isSession {
				err = config.Session.check(c)
				if err != nil {
					// Если cookie не существует, то перенаправляем запрос условно на "/signIn"
					return c.Redirect(config.Session.RedirectPath, config.Session.RedirectStatus)
				}
			}
		}

		// Доступ к маршрутам
		if r.isPermission && config.Permission != nil {
			if c.Identity != nil {
				if !config.Permission.check(c.Identity.Username, c.Path()) {
					if config.Permission.NotPermissionHandler != nil {
						return config.Permission.NotPermissionHandler(c, StatusForbidden, "Forbidden")
					}
					return c.SendStatus(StatusForbidden)
				}
			}
		}

		// Обычный маршрут
		return r.Handler(c)
	}
}
