package swagger

import v2 "github.com/egovorukhin/egowebapi/swagger/v2"

type Config struct {
	//Version     string
	Host string
	Info v2.Info
}

const (
	V2 = "v2"
	V3 = "v3"
)
