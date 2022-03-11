package egowebapi

type Route struct {
	params       []string
	auth         Auth
	Handler      Handler
	isSession    bool
	isPermission bool
	Sign         Sign
	//webSocket    *WebSocket
	//option       Option
}

/*type WebSocket struct {
	UpgradeHandler fiber.Handler
}*/

type Option struct {
	Headers     []string `json:"headers,omitempty"`
	Method      string   `json:"method,omitempty"`
	Body        string   `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Auth string

const (
	NoAuth     Auth = "NoAuth"
	BasicAuth  Auth = "BasicAuth"
	DigestAuth Auth = "DigestAuth"
	ApiKeyAuth Auth = "ApiKeyAuth"
)

type Sign int

const (
	SignNone Sign = iota
	SignIn
	SignOut
	SignUp
)

// SetParams указываем параметры маршрута
func (r *Route) SetParams(params ...string) *Route {
	r.params = params
	return r
}

func (r *Route) SetSign(sign Sign) *Route {
	r.Sign = sign
	return r
}

// SetDescription устанавливаем описание маршрута
/*func (r *Route) SetDescription(s string) *Route {
	r.option.Description = s
	return r
}

// SetBody устанавливаем описание тела маршрута
func (r *Route) SetBody(s string) *Route {
	r.option.Body = s
	return r
}*/

// Auth указываем метод авторизации
func (r *Route) Auth(auth Auth) *Route {
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

// WebSocket Устанавливаем web socket соединение
/*func (r *Route) WebSocket(upgrade fiber.Handler) *Route {
	r.webSocket = &WebSocket{
		UpgradeHandler: upgrade,
	}
	return r
}*/

func (r *Route) EmptyHandler() {
	r.Handler = nil
}

// GetHandler возвращаем обработчик основанный на параметрах конфигурации маршрута
func (r *Route) getHandler(config Config, swagger *Swagger) Handler {

	return func(c *Context) error {

		var err error
		switch r.auth {
		case BasicAuth:
			if config.Authorization.Basic != nil {
				c.Identity, err = config.Authorization.Basic.Do(c)
			}
			break
		case DigestAuth:
			if config.Authorization.Digest != nil {
				c.Identity, err = config.Authorization.Digest.Do(c)
			}
			break
		case ApiKeyAuth:
			if config.Authorization.ApiKey != nil {
				c.Identity, err = config.Authorization.ApiKey.Do(c)
			}
			break
		}

		// Проверка на авторизацию
		if err != nil {
			if config.Authorization.Unauthorized != nil {
				return config.Authorization.Unauthorized(c, StatusUnauthorized, err)
			}
			c.Set(HeaderWWWAuthenticate, err.Error())
			return c.SendStatus(StatusUnauthorized)
		}

		// Проверяем маршрут на актуальность сессии
		if /*r.isSession &&*/ config.Session != nil /*|| r.isSession*/ {
			// Маршрут для входа/выхода
			switch r.Sign {
			case SignIn:
				config.Session.login(c)
				break
			case SignOut:
				config.Session.logout(c)
				return c.Redirect(config.Session.RedirectPath, StatusFound)
			}
			if r.isSession {
				// Проверка сессии
				err = config.Session.check(c)
				if err != nil {
					// Если cookie не существует, то перенаправляем запрос на условно "/login"
					return c.Redirect(config.Session.RedirectPath, StatusFound)
				}
			}
		}

		// Доступ к маршрутам
		if r.isPermission && config.Permission != nil {
			if c.Identity != nil {
				if !config.Permission.check(c.Identity.SessionId, c.Path()) {
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
	/*switch h := r.Handler.(type) {
	// handler для маршрутов с identity
	case func(*fiber.Ctx, *Identity) error:
		// Авторизация
		switch r.Authorization {
		case BasicAuth:
			if config.Authorization.Basic != nil {
				return config.Authorization.Basic.Do(h, r.IsPermission, config.Permission)
			}
			break
		case DigestAuth:
			if config.Authorization.Digest != nil {
				return config.Authorization.Digest.Do(h, r.IsPermission, config.Permission)
			}
			break
		case ApiKeyAuth:
			if config.Authorization.ApiKey != nil {
				return config.Authorization.ApiKey.Do(h, r.IsPermission, config.Permission)
			}
			break
		}

		// Проверяем маршрут на актуальность сессии
		if (r.IsSession && config.Session != nil) || r.IsSession {
			return config.Session.check(h, r.IsPermission, config.Permission)
		}
		return func(ctx *fiber.Ctx) error {
			return h(ctx, nil)
		}

	// Swagger handler для добавления описания маршрутов
	case func(*fiber.Ctx, *Swagger) error:
		return func(ctx *fiber.Ctx) error {
			return h(ctx, swagger)
		}

	// LoginHandler для маршрута web авторизации Login
	case func(*fiber.Ctx, string) error:
		if config.Session != nil {
			return config.Session.login(h)
		}
		break
	// LogoutHandler для маршрута web авторизации Logout
	case func(*fiber.Ctx, *Identity, string) error:
		if config.Session != nil {
			return config.Session.logout(h)
		}
		break

	// Handler для маршрут WebSocket соединения
	case func(*websocket.Conn):
		return websocket.New(h)

	// Обычный обработчик без ништяков
	case func(*fiber.Ctx) error:
		return h
	}

	// Ну если ни один из обработчиков не удовлетворяет требованиям, то вернем ответ с кодом 404
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("%s %s", ctx.Route().Method, ctx.Route().Path))
	}*/
}
