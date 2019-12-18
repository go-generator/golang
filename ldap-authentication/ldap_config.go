package ldap_authentication

type LDAPConfig struct {
	Server        string `mapstructure:"server"`
	BindingFormat string `mapstructure:"binding_format"`
}
