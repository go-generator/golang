package ldap_authentication

import "github.com/common-go/auth"

type BasicAuthenticator interface {
	Authenticate(auth auth.AuthInfo) (auth.AuthResult, error)
}
