package egowebapi

// Permission Структура описывает разрешения на запрос
/*type Permission struct {
	Check PermissionHandler
	Error ErrorHandler
}

const StatusForbidden = "Доступ запрещен (Permission denied)"

// Проверяем запрос на разрешения
func (p *Permission) check(handler WebHandler) Handler {
	return func(ctx *fiber.Ctx) error {

		if p.Check != nil {
			key := ctx.Cookies(sessionId)
			route := ctx.Route()
			if !p.Check(key, route.Path) {
				if p.Error != nil {
					return p.Error(ctx, 403, StatusForbidden)
				}
				return ctx.Status(403).SendString(StatusForbidden)
			}
		}

		return handler(ctx, nil)
	}
}*/
