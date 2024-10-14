package security

import "fmt"

type Digest struct {
	Handler DigestAuthHandler

	header string
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

func (d *Digest) Name() string {
	return DigestAuth
}

func (d *Digest) Do() (identity *Identity, err error) {
	return
}

func (d *Digest) Definition() Definition {
	return Definition{
		Type:        TypeDigest,
		Description: fmt.Sprintf("Digest Authorization"),
	}
}

func (d *Digest) SetHeader(header string) *Digest {
	d.header = header
	return d
}
