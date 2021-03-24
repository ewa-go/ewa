package user

import "fmt"

var users = map[string]*User{}

type User struct {
	Lastname  string
	Firstname string
}

func Set(id string, lastname, firstname string) {
	users[id] = &User{
		Lastname:  lastname,
		Firstname: firstname,
	}
}

func GetUsers() map[string]*User {
	return users
}

func Get(id string) *User {
	if v, ok := users[id]; ok {
		return v
	}
	return nil
}

func Delete(id string) {
	delete(users, id)
}

func (u *User) String() string {
	return fmt.Sprintf("lastname: %s, firstname: %s", u.Lastname, u.Firstname)
}
