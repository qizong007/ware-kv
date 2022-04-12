package authentication

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	rootName := "qizong007"
	option := &AuthCenterOption{
		Username: rootName,
		Password: "123456",
	}
	center := NewAuthCenter(option)
	registerOption1 := &AuthRegisterOption{
		ParentUser: rootName,
		Username:   "user1",
		Password:   "pswd1",
		Role:       Admin,
	}
	if err := center.Register(registerOption1); err != nil {
		fmt.Println(err)
		return
	}
	registerOption2 := &AuthRegisterOption{
		ParentUser: "user1",
		Username:   "user2",
		Password:   "pswd2",
		Role:       Reader,
	}
	if err := center.Register(registerOption2); err != nil {
		fmt.Println(err)
		return
	}
	registerOption3 := &AuthRegisterOption{
		ParentUser: "user2",
		Username:   "user3",
		Password:   "pswd3",
		Role:       Writer,
	}
	if err := center.Register(registerOption3); err != nil {
		fmt.Println(err)
	}
	registerOption4 := &AuthRegisterOption{
		ParentUser: "user1",
		Username:   "user3",
		Password:   "pswd3",
		Role:       Writer,
	}
	if err := center.Register(registerOption4); err != nil {
		fmt.Println(err)
	}
	registerOption5 := &AuthRegisterOption{
		ParentUser: "user1",
		Username:   "user2",
		Password:   "pswd5",
		Role:       Writer,
	}
	if err := center.Register(registerOption5); err != nil {
		fmt.Println(err)
	}
	for _, user := range center.userTable {
		fmt.Println(user.role, user.username, user.password)
	}
}
