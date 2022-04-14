package models

import "errors"

var users Users

type User struct {
	Id        string `json:"id,omitempty"`
	Firstname string `json:"firstname" jsonschema:"required,format=string,description=Имя"`
	Lastname  string `json:"lastname" jsonschema:"description=Фамилия"`
}

type Users map[string]*User

func GetUser(id string) *User {
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

func GetUsers() Users {
	if users == nil {
		return nil
	}
	return users
}

func (u User) Set() {
	if users == nil {
		return
	}
	users[u.Id] = &u
}

func (u User) Update(id string) error {
	if users == nil {
		return nil
	}
	if _, ok := users[id]; ok {
		users[id] = &u
		return nil
	}
	return errors.New("Запись не найдена - " + id)
}

func (u User) Delete() {
	if users == nil {
		return
	}
	delete(users, u.Id)
}
