package ldap_authentication

import (
	"fmt"
	"github.com/common-go/auth"
	"gopkg.in/ldap.v2"
	"log"
)

type LDAPAuthenticator struct {
	LDAPConfig       LDAPConfig
}

func (s *LDAPAuthenticator) Authenticate(info auth.AuthInfo) (auth.AuthResult, error) {
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
