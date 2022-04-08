package egowebapi

import "github.com/egovorukhin/egowebapi/swagger/v2"

type Route struct {
	params       []string
	auth         []string
	Handler      Handler
	isSession    bool
	isPermission bool
	sign         Sign
	path         *v2.Path
}

type Map map[string]interface{}

type Sign int

const (
	SignNone Sign = iota
	SignIn
	SignOut
)

// SetParams указываем параметры маршрута
func (r *Route) SetParams(params ...string) *Route {
	r.params = params
	return r
}

// SetParameters указываем информацию о параметрах адресной строки для Swagger
func (r *Route) SetParameters(params ...v2.Parameter) *Route {
	if r.path == nil {
		r.path = &v2.Path{}
	}
	r.path.Parameters = params
	return r
}

// SetConsumes устанавливаем Content-Type запроса для Swagger
func (r *Route) SetConsumes(c ...string) *Route {
	if r.path == nil {
		r.path = &v2.Path{}
	}
	r.path.Consumers = c
	return r
}

// SetProduces устанавливаем Content-Type ответа для Swagger
func (r *Route) SetProduces(p ...string) *Route {
	if r.path == nil {
		r.path = &v2.Path{}
	}
	r.path.Produces = p
	return r
}

// SetResponse описываем варианты ответов для Swagger
func (r *Route) SetResponse(code int, resp v2.Response) *Route {
	if r.path == nil {
		r.path = &v2.Path{}
	}
	if r.path.Responses == nil {
		r.path.Responses = map[int]v2.Response{}
	}
	r.path.Responses[code] = resp
	return r
}

// SetSign устанавливаем вариант входа/выхода для маршрута
func (r *Route) SetSign(sign Sign) *Route {
	r.sign = sign
	return r
}

// Auth указываем метод авторизации
func (r *Route) Auth(auth ...string) *Route {
	r.auth = auth
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
func (r *Route) getHandler(config Config, view *View, swagger *v2.Swagger) Handler {

	return func(c *Context) error {

		c.View = view
		c.Swagger = swagger

		var (
			err    error
			isAuth bool
		)
		if len(r.auth) > 0 {
			isAuth = true
		}
		for _, auth := range r.auth {
			switch auth {
			case BasicAuth:
				if config.Authorization.Basic != nil {
					err = config.Authorization.Basic.Do(c)
					if err != nil {
						c.Set(HeaderWWWAuthenticate, err.Error())
					}
				}
				break
			case DigestAuth:
				if config.Authorization.Digest != nil {
					err = config.Authorization.Digest.Do(c)
				}
				break
			case ApiKeyAuth:
				if config.Authorization.ApiKey != nil {
					err = config.Authorization.ApiKey.Do(c)
				}
				break
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
				if config.Authorization.Unauthorized != nil {
					return config.Authorization.Unauthorized(c, StatusUnauthorized, err)
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
