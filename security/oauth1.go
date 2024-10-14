package security

type OAuth1 struct{}

func (OAuth1) Name() string {
	return OAuth1Auth
}

func (OAuth1) Do() (identity *Identity, err error) {
	return
}

func (OAuth1) Definition() Definition {
	return Definition{}
}
