package ldap_authentication

import "github.com/common-go/auth"

type AuthenticationService interface {
	Authenticate(user auth.AuthInfo) (auth.AuthResult, error)
}
