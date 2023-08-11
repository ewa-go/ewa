package security

type Digest struct {
	Handler DigestAuthHandler
}

type Advanced struct {
	Realm       string
	Nonce       string
	Algorithm   string
	Qop         string
	NonceCount  string
	ClientNonce string
	Opaque      string
}

type DigestAuthHandler func(user string, pass string, advanced Advanced) bool

const (
	Md5Algorithm           = "MD5"
	Md5SessAlgorithm       = "MD5-sess"
	Sha256Algorithm        = "SHA-256"
	Sha256SessAlgorithm    = "SHA-256-sess"
	Sha512256Algorithm     = "SHA-512-256"
	Sha512256SessAlgorithm = "SHA-512-256-sess"
)

func (Digest) Do() (identity *Identity, err error) {
	return
}

func (Digest) Definition() Definition {
	return Definition{}
}
