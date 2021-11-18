package auth

const (
	Md5Algorithm           = "MD5"
	Md5SessAlgorithm       = "MD5-sess"
	Sha256Algorithm        = "SHA-256"
	Sha256SessAlgorithm    = "SHA-256-sess"
	Sha512256Algorithm     = "SHA-512-256"
	Sha512256SessAlgorithm = "SHA-512-256-sess"
)

type Authorization struct {
	Basic  *Basic
	Digest *Digest
	ApiKey *ApiKey
}

type BasicAuthHandler func(user string, pass string) bool
type DigestAuthHandler func(user string, pass string, advanced Advanced) bool
type ApiKeyHandler func(key string, value string) bool
