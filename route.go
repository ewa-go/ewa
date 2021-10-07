package egowebapi

type Route struct {
	Params        []string
	Description   string
	IsBasicAuth   bool
	IsJWT         bool
	IsSession     bool
	IsPermission  bool
	Handler       Handler
	WebHandler    WebHandler
	LoginHandler  AuthHandler
	LogoutHandler AuthHandler
	WsHandler
}

func (r *Route) SetHandler(handler Handler) *Route {
	r.Handler = handler
	return r
}

func (r *Route) SetParams(params ...string) *Route {
	r.Params = params
	return r
}

func (r *Route) SetDescription(s string) *Route {
	r.Description = s
	return r
}

// BasicAuth Вешаем флаг авторизации Basic
func (r *Route) BasicAuth() *Route {
	r.IsBasicAuth = true
	return r
}

// JWT Вешаем флаг авторизации JW Token
func (r *Route) JWT() *Route {
	r.IsBasicAuth = true
	return r
}

// Session Вешаем получение аутентификации сессии, IsPermission ставим флаг для проверки маршрута на право доступа
func (r *Route) Session(IsPermission bool) *Route {
	r.IsSession = true
	r.IsPermission = IsPermission
	return r
}
