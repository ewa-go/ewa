package egowebapi

// Permission структура описывает разрешения на запрос
type Permission struct {
	Handler PermissionHandler
}

// Check Проверяем запрос на разрешения
func (p *Permission) Check(id interface{}, path string) bool {
	if p.Handler == nil {
		return true
	}
	return p.Handler(id, path)
}
