package security

type OAuth2 struct {
}

func (OAuth2) Do() (identity *Identity, err error) {
	return
}

func (OAuth2) Definition() Definition {
	return Definition{}
}
