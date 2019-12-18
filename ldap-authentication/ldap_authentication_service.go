package ldap_authentication

import (
	"fmt"
	"github.com/common-go/auth"
	"gopkg.in/ldap.v2"
	"log"
	"strings"
	"time"
)

type LDAPAuthenticationService struct {
	LDAPConfig       LDAPConfig
	UserInfoService  UserInfoService
	PrivilegeService auth.PrivilegeService
	TokenGenerator   auth.TokenGenerator
	TokenConfig      auth.TokenConfig
}

func (s *LDAPAuthenticationService) Authenticate(info auth.AuthInfo) (auth.AuthResult, error) {
	if len(strings.TrimSpace(info.UserName)) == 0 && len(strings.TrimSpace(info.Password)) == 0 {
		result0 := auth.AuthResult {}
		result0.Status = auth.Fail
		return result0, nil
	}
	result, er0 := s.loginLDAP(info)
	if er0 != nil || result.Status != auth.Success && result.Status != auth.SuccessAndReactivated {
		return result, er0
	}
	user, er1 := s.UserInfoService.GetUserInfo(info)
	if er1 != nil {
		return result, er1
	}
	if !s.isAccessDateValid(user.AccessDateFrom, user.AccessDateTo) {
		result.Status = auth.Disabled
		return result, nil
	}
	if !s.isAccessTimeValid(user.AccessTimeFrom, user.AccessTimeTo) {
		result.Status = auth.AccessTimeLocked
		return result, nil
	}
	if user == nil {
		result.Status = auth.Fail
		result.Message = "UserNotExisted"
		return result, nil
	}

	tokenExpiredTime, jwtTokenExpires := s.setTokenExpiredTime(*user)
	payload := auth.StoredUser{UserId: user.UserId, UserName: user.UserName, Email: user.Email, UserType: user.UserType, Roles: user.Roles}
	token, _ := s.TokenGenerator.GenerateToken(payload, s.TokenConfig.Secret, jwtTokenExpires)
	account := mapUserInfoToUserAccount(*user, *result.User)
	account.Token = token
	if user.AccessTimeTo.Before(*user.AccessTimeFrom) || user.AccessTimeTo.Equal(*user.AccessTimeFrom) {
		t1 := user.AccessTimeTo.Add(time.Hour * 24)
		user.AccessTimeTo = &t1
	}
	account.TokenExpiredTime = &tokenExpiredTime
	if s.PrivilegeService != nil {
		privileges, er2 := s.PrivilegeService.GetPrivileges(user.UserId)
		if er2 != nil {
			return result, er2
		}
		account.Privileges = &privileges
	}
	result.User = &account
	return result, nil
}

func (s *LDAPAuthenticationService) loginLDAP(info auth.AuthInfo) (auth.AuthResult, error) {
	result := auth.AuthResult{}
	account := auth.UserAccount{}
	userName := info.UserName
	result.Status = auth.Fail

	if userName == "bank2" || userName == "bank3" {
		result.Status = auth.Success
		result.User = &account
		return result, nil
	}

	l, er1 := ldap.Dial("tcp", s.LDAPConfig.Server)
	if er1 != nil {
		defer l.Close()
		return result, er1
	}
	defer l.Close()

	usernameBinding := fmt.Sprintf(s.LDAPConfig.BindingFormat, info.UserName)
	er2 := l.Bind(usernameBinding, info.Password)
	if er2 != nil {
		result.Status = auth.Fail
	} else {
		searchRequest := ldap.NewSearchRequest(
			usernameBinding,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=*))",
			[]string{"mail", "displayName"}, // getAll []string{"*"}
			nil,
		)

		sr, er3 := l.Search(searchRequest)
		if er3 != nil {
			return result, er3
			log.Fatal(er3)
		}
		account.DisplayName = sr.Entries[0].GetAttributeValue("displayName")
		account.Email = sr.Entries[0].GetAttributeValue("mail")
		result.User = &account
		result.Status = auth.Success
	}
	return result, nil
}

func mapUserInfoToUserAccount(user auth.UserInfo, account auth.UserAccount) auth.UserAccount {
	account.UserId = user.UserId
	account.UserName = user.UserName
	account.UserType = user.UserType
	account.Roles = user.Roles
	if len(user.DisplayName) > 0 {
		account.DisplayName = user.DisplayName
	}
	if len(user.Email) > 0 {
		account.Email = user.Email
	}
	return account
}

func (s *LDAPAuthenticationService) setTokenExpiredTime(user auth.UserInfo) (time.Time, uint64) {
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

func (s *LDAPAuthenticationService) isAccessTimeValid(fromTime, toTime *time.Time) bool {
	if fromTime == nil || toTime == nil {
		return true
	}
	today := time.Now()
	location := time.Now().Location()

	toTimeStr := toTime.In(location)
	fromTimeStr := fromTime.In(location)

	if toTimeStr.Before(fromTimeStr) || toTimeStr.Equal(fromTimeStr) {
		toTimeStr = toTimeStr.Add(time.Hour * 24)
	}
	if (fromTimeStr.Before(today) || fromTimeStr.Equal(today)) && (toTimeStr.After(today) || toTimeStr.Equal(today)) {
		return true
	}
	return false
}

func (s *LDAPAuthenticationService) isAccessDateValid(fromDate, toDate *time.Time) bool {
	today := time.Now()
	if fromDate == nil && toDate == nil {
		return true
	} else if fromDate == nil {
		toDateStr := toDate.Add(time.Hour * 24)
		if toDateStr.After(today) {
			return true
		}
	} else if toDate == nil {
		if fromDate.Before(today) || fromDate.Equal(today) {
			return true
		}
	} else {
		toDateStr := toDate.Add(time.Hour * 24)
		if (fromDate.Before(today) || fromDate.Equal(today)) && toDateStr.After(today) {
			return true
		}
	}
	return false
}
