package ewa

import (
	"fmt"
	"testing"
)

func TestNewSuffix(t *testing.T) {

	hostname := Suffix{
		Index:       2,
		Value:       "hostname",
		isParam:     false,
		Description: "Hostname",
	}

	client := Suffix{
		Index:       2,
		Value:       "client",
		isParam:     false,
		Description: "Client",
	}

	s := NewSuffix(hostname, client)
	fmt.Printf("%#v", s)
}
