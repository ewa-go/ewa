package security

import (
	"fmt"
	"testing"
)

func TestBasicParse(t *testing.T) {
	b := Basic{
		header:  "Basic dXNlcjpRcTEyOjM0NTY=",
		Handler: nil,
	}
	username, pass, ok := b.parseBasicAuth()
	fmt.Println(username, pass, ok)
}
