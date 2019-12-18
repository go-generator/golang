package auth_service

import "github.com/common-go/auth"

type AuthenticationService interface {
	Authenticate(user auth.AuthInfo) (auth.AuthResult, error)
}
