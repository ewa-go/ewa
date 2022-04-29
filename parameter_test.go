package egowebapi

import (
	"fmt"
	"testing"
)

func TestParameters(t *testing.T) {

	type Person struct {
		Id        int    `json:"id"`
		Firstname string `json:"firstname"`
	}

	param := NewInBody(true, NewSchema(Person{}), "Описание")
	fmt.Printf("In Body: %+v\n", param)

	param = NewInPath("/{id}", true, "Описание").SetType(TypeInteger)
	fmt.Printf("In Path: %+v\n", param)

	param = NewInQuery("id", false, "Описание")
	fmt.Printf("In Query: %+v\n", param)

	param = NewInHeader("id", false, "Описание")
	fmt.Printf("In Header: %+v\n", param)

	param = NewInFormData("file", TypeFile, true, "Описание")
	fmt.Printf("In FormData: %+v\n", param)

	param = NewInFormData("id", TypeString, true, "Описание")
	fmt.Printf("In FormData: %+v\n", param)
}

func TestParameter(t *testing.T) {

	param := NewParameter("id").SetType(TypeInteger).SetDescription("Описание").SetRequired(true).SetFormat("int32")
	fmt.Printf("Param created: %+v\n", param)
}

func TestParameter_SetTypeFormat(t *testing.T) {

	param := NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(0)
	fmt.Printf("Param int: %+v\n", param)

	var i64 int64 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i64)
	fmt.Printf("Param int64: %+v\n", param)

	var i32 int32 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i32)
	fmt.Printf("Param int32: %+v\n", param)

	var i16 int16 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i16)
	fmt.Printf("Param int16: %+v\n", param)

	var i8 int8 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i8)
	fmt.Printf("Param int8: %+v\n", param)

	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat("")
	fmt.Printf("Param string: %+v\n", param)

	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(true)
	fmt.Printf("Param boolean: %+v\n", param)
}
