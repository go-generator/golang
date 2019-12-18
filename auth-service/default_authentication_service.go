package auth_service

import (
	"github.com/common-go/auth"
	"strings"
	"time"
)

type DefaultAuthenticationService struct {
	UserInfoService    UserInfoService
	PrivilegeService   auth.PrivilegeService
	PasswordComparator Comparator
	TokenGenerator     auth.TokenGenerator
	TokenConfig        TokenConfig
	EncryptPasswordKey string
}

func (s *DefaultAuthenticationService) Authenticate(info auth.AuthInfo) (auth.AuthResult, error) {
	result := auth.AuthResult{}
	result.Status = auth.Fail

	userName := info.UserName
	password := info.Password
	if len(strings.TrimSpace(userName)) == 0 || len(strings.TrimSpace(password)) == 0 {
		result.Status = auth.Fail
		return result, nil
	}
	/*
		if len(s.EncryptPasswordKey) > 0 {
			if passwordDecoded, err := security.DecodeRC4([]byte(password), []byte(s.EncryptPasswordKey)); err != nil {
				result.Status = auth.Fail
				return result
			} else {
				password = string(passwordDecoded)
			}
		}
	*/
	user, er1 := s.UserInfoService.GetUserInfo(info)
	if er1 != nil {
		if er1.Error() == "Out of access time." {
			result.Status = auth.AccessTimeLocked
		} else {
			result.Status = auth.SystemError
		}
		result.Message = er1.Error()
		return result, er1
	}
	if user == nil {
		result.Status = auth.Fail
		return result, nil
	}
	validPassword := s.PasswordComparator.Compare(password, user.Password)
	if !validPassword {
		result.Status = auth.WrongPassword
		s.UserInfoService.HandleWrongPassword(*user)
		return result, nil
	}
	if user.Suspended {
		result.Status = auth.Suspended
		return result, nil
	}

	locked := user.LockedUntilTime != nil && (compareDate(time.Now(), *user.LockedUntilTime) < 0)
	if locked {
		result.Status = auth.Locked
		return result, nil
	}

	var passwordExpiredTime *time.Time = nil // date.addDays(time.Now(), 10)
	if user.PasswordModifiedTime != nil && user.MaxPasswordAge != 0 {
		t := addDays(*user.PasswordModifiedTime, user.MaxPasswordAge)
		passwordExpiredTime = &t
	}
	if passwordExpiredTime != nil && compareDate(time.Now(), *passwordExpiredTime) > 0 {
		result.Status = auth.PasswordExpired
		return result, nil
	}
	tokenExpiredTime, jwtTokenExpires := s.setTokenExpiredTime(*user)
	storedUser := auth.StoredUser{UserId: user.UserId, UserName: user.UserName, Email: user.Email, UserType: user.UserType, Roles: user.Roles}
	token, er2 := s.TokenGenerator.GenerateToken(storedUser, s.TokenConfig.Secret, uint64(jwtTokenExpires))
	if er2 != nil {
		result.Status = auth.SystemError
		result.Message = er2.Error()
		return result, nil
	}

	if user.Deactivated == true {
		result.Status = auth.SuccessAndReactivated
	} else {
		result.Status = auth.Success
	}
	account := mapUserInfoToUserAccount(*user)
	account.NewUser = false
	account.Token = token
	account.TokenExpiredTime = &tokenExpiredTime
	account.PasswordExpiredTime = passwordExpiredTime
	if s.PrivilegeService != nil {
		privileges, er3 := s.PrivilegeService.GetPrivileges(user.UserId)
		if er3 != nil {
			result.Status = auth.SystemError
			result.Message = er3.Error()
			return result, er2
		}
		account.Privileges = &privileges
	}
	result.User = &account
	_, er3 := s.UserInfoService.PassAuthentication(*user)
	if er3 != nil {
		return result, er3
	}
	result.Status = auth.Success
	return result, nil
}

func mapUserInfoToUserAccount(user auth.UserInfo) auth.UserAccount {
	account := auth.UserAccount{}
	account.UserId = user.UserId
	account.UserName = user.UserName
	account.Email = user.Email
	account.DisplayName = user.DisplayName
	account.UserType = user.UserType
	account.Roles = user.Roles
	return account
}

func (s *DefaultAuthenticationService) setTokenExpiredTime(user auth.UserInfo) (time.Time, uint64) {
	if user.AccessTimeTo == nil || user.AccessTimeFrom == nil || user.AccessDateFrom == nil || user.AccessDateTo == nil {
		var tokenExpiredTime = addSeconds(time.Now(), int(s.TokenConfig.Expires/1000))
		return tokenExpiredTime, s.TokenConfig.Expires
	}
	if user.AccessTimeTo.Before(*user.AccessTimeFrom) || user.AccessTimeTo.Equal(*user.AccessTimeFrom) {
		tmp := user.AccessTimeTo.Add(time.Hour * 24)
		user.AccessTimeTo = &tmp
	}
	var tokenExpiredTime time.Time
	var jwtExpiredTime uint64
	if time.Millisecond*time.Duration(s.TokenConfig.Expires) > user.AccessTimeTo.Sub(time.Now()) {
		tokenExpiredTime = time.Now().Add(user.AccessTimeTo.Sub(time.Now())).UTC()
		jwtExpiredTime = uint64(user.AccessTimeTo.Sub(time.Now()).Seconds() * 1000)
	} else {
		tokenExpiredTime = time.Now().Add(time.Millisecond * time.Duration(s.TokenConfig.Expires)).UTC()
		jwtExpiredTime = s.TokenConfig.Expires
	}
	return tokenExpiredTime, jwtExpiredTime
}
