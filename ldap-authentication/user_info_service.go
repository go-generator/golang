package ldap_authentication

import "github.com/common-go/auth"

type UserInfoService interface {
	GetUserInfo(auth auth.AuthInfo) (*auth.UserInfo, error)
	PassAuthentication(user auth.UserInfo) (bool, error)
	HandleWrongPassword(user auth.UserInfo) (bool, error)
}
