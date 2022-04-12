package authentication

import (
	"fmt"
	"github.com/qizong007/ware-kv/util"
	"sync"
)

// auth role
const (
	Admin = iota // Reader + Writer + register auth
	Reader
	Writer
)

func GetAuthRoleIntFromStr(str string) int {
	switch str {
	case "Admin":
		return Admin
	case "Reader":
		return Reader
	case "Writer":
		return Writer
	default:
		return Reader // default is Reader
	}
}

var authCenter *AuthCenter

type AuthCenter struct {
	userTable map[string]*WareUser
	rw        sync.RWMutex
}

type WareUser struct {
	username string
	password string
	role     int
}

type AuthCenterOption struct {
	Username string
	Password string
}

type WareAuthOption struct {
	Open bool `yaml:"Open"`
	Root struct {
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"Root"`
	Others []struct {
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		Role     string `yaml:"Role"`
	} `yaml:"Others"`
}

func DefaultOption() *WareAuthOption {
	return &WareAuthOption{
		Open: false,
	}
}

func NewAuthCenter(option *AuthCenterOption) *AuthCenter {
	rootAdmin := &WareUser{
		username: option.Username,
		password: option.Password,
		role:     Admin,
	}
	authCenter = &AuthCenter{
		userTable: map[string]*WareUser{rootAdmin.username: rootAdmin},
	}
	return authCenter
}

func GetAuthCenter() *AuthCenter {
	return authCenter
}

type AuthRegisterOption struct {
	ParentUser string // parent's username (role should be 'Admin')
	Username   string
	Password   string
	Role       int
}

type AuthCancelOption struct {
	ParentUser string // parent's username (role should be 'Admin')
	Username   string
}

func (a *AuthCenter) Register(option *AuthRegisterOption) error {
	a.rw.Lock()
	defer a.rw.Unlock()
	var (
		parentUser *WareUser
		ok         bool
	)

	// parent not existed
	if parentUser, ok = a.userTable[option.ParentUser]; !ok {
		return fmt.Errorf(util.ErrCode2Msg[util.UserNotExist])
	}

	// only 'Admin' can register for others
	if parentUser.role != Admin {
		return fmt.Errorf(util.ErrCode2Msg[util.PermissionDenied])
	}

	newName := option.Username
	// user existed
	if _, ok = a.userTable[newName]; ok {
		return fmt.Errorf(util.ErrCode2Msg[util.UserExisted])
	}

	a.userTable[newName] = &WareUser{
		username: newName,
		password: option.Password,
		role:     option.Role,
	}
	return nil
}

func (a *AuthCenter) Cancel(option *AuthCancelOption) error {
	a.rw.Lock()
	defer a.rw.Unlock()
	var (
		parentUser *WareUser
		ok         bool
	)

	// parent not existed
	if parentUser, ok = a.userTable[option.ParentUser]; !ok {
		return fmt.Errorf(util.ErrCode2Msg[util.UserNotExist])
	}

	// only 'Admin' can register for others
	if parentUser.role != Admin {
		return fmt.Errorf(util.ErrCode2Msg[util.PermissionDenied])
	}

	delete(a.userTable, option.Username)
	return nil
}

func (a *AuthCenter) GetReaders() map[string]string {
	a.rw.RLock()
	defer a.rw.RUnlock()
	readers := a.getRoles([]int{Admin, Reader})
	return readers
}

func (a *AuthCenter) GetWriters() map[string]string {
	a.rw.RLock()
	defer a.rw.RUnlock()
	writers := a.getRoles([]int{Admin, Writer})
	return writers
}

func (a *AuthCenter) getRoles(roleList []int) map[string]string {
	// username -> password
	roleTable := make(map[string]string)
	for _, user := range a.userTable {
		if util.IsIntInList(roleList, user.role) {
			roleTable[user.username] = user.password
		}
	}
	return roleTable
}
