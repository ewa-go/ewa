package security

type BearerToken struct {
}

func (BearerToken) Do() (identity *Identity, err error) {
	return
}

func (BearerToken) Definition() Definition {
	return Definition{}
}
