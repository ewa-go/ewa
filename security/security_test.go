package security

import (
	"fmt"
	"testing"
)

func TestAuthorization_Get(t *testing.T) {

	a := Authorization{
		Unauthorized: func(err error) bool {
			fmt.Println(err)
			return true
		},
		Basic: &Basic{
			Handler: func(user string, pass string) bool {
				fmt.Println(user, pass)
				return true
			},
		},
		ApiKey: &ApiKey{
			KeyName: "Token",
			Param:   ParamHeader,
			Handler: func(token string) (username string, err error) {
				fmt.Println(token)
				return "user", nil
			},
		},
	}

	def := a.Get(BasicAuth).Definition()
	fmt.Printf("%v+\n", def)
	def = a.Get(ApiKeyAuth).Definition()
	fmt.Printf("%v+\n", def)
}
