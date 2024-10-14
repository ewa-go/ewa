package security

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
