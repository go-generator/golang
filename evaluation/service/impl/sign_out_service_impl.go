package impl

import (
	"time"
)

type SignOutServiceImpl struct {
	// BlacklistTokenService    TokenBlacklistService
}

func NewSignOutServiceImpl() *SignOutServiceImpl {
	//return &SignOutServiceImpl{blacklistTokenService}
	return &SignOutServiceImpl{}
}

func (s *SignOutServiceImpl) SignOut(token string, timeExpires time.Time) (bool, error) {
	//return s.BlacklistTokenService.Revoke(token, "The token has signed out.", timeExpires), nil
	return true, nil
}
