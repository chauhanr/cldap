package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"

	model "github.com/chauhanr/cldap/models"
	"github.com/go-ldap/ldap"
)

type LDAPClient struct {
	Base               string
	BindDN             string
	BindPassword       string
	GroupFilter        string
	UserFilter         string
	ServerName         string
	Host               string
	Conn               *ldap.Conn
	Port               int
	InsecureSkipVerify bool
	UseSSL             bool
	ClientCertificates []tls.Certificate
	Attributes         []string
}

func InitializeClient(login bool) (*LDAPClient, error) {
	lc := model.LdapConfig{}
	err := lc.LoadConfig()
	if err != nil {
		return nil, err
	}
	if !lc.Ldap.Client.HasConfig {
		err = errors.New("Ldap configurations missing\n")
		return nil, err
	}
	if !login {
		if !lc.Ldap.Creds.HasCreds {
			err = errors.New("Ldap creds missing, add creds to continue.\n")
			return nil, err
		}
	}

	cl := lc.Ldap.Client

	c := LDAPClient{
		Base:               cl.Base,
		BindDN:             cl.BindDN,
		BindPassword:       cl.BindPassword,
		UserFilter:         cl.UserFilter,
		GroupFilter:        cl.GroupFilter,
		Host:               cl.Host,
		Port:               cl.Port,
		UseSSL:             cl.UseSSL,
		InsecureSkipVerify: cl.InsecureSkipVerify,
		Attributes:         cl.Attributes,
	}

	return &c, nil
}

func (lc *LDAPClient) Connect() error {
	if lc.Conn == nil {
		var l *ldap.Conn
		var err error
		address := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
		if !lc.UseSSL {
			l, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}
		} else {
			config := &tls.Config{
				InsecureSkipVerify: lc.InsecureSkipVerify,
				ServerName:         lc.ServerName,
			}
			if lc.ClientCertificates != nil && len(lc.ClientCertificates) > 0 {
				config.Certificates = lc.ClientCertificates
			}
			l, err = ldap.DialTLS("tcp", address, config)
			if err != nil {
				return err
			}
		}
		lc.Conn = l
	}
	return nil
}

func (lc *LDAPClient) SearchAndBind(username, password string) (bool, map[string]string, error) {
	err := lc.Connect()
	user := map[string]string{}
	if err != nil {
		return false, nil, err
	}

	if lc.BindDN != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return false, nil, err
		}
	}
	attributes := append(lc.Attributes, "dn")

	sReq := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(sReq)
	if err != nil {
		return false, nil, err
	}
	if len(sr.Entries) < 1 {
		return false, nil, errors.New("No users found.")
	}
	if len(sr.Entries) > 1 {
		return false, nil, errors.New("Too many matches found")
	}

	userDN := sr.Entries[0].DN
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}
	// bind to authenticate
	err = lc.Conn.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}
	// Rebind the read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err = lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return false, user, err
		}
	}
	return true, user, nil
}

type User struct {
	UserName   string
	Attributes []string
}

func (lc *LDAPClient) SearchUser(username string) ([]User, error) {
	err := lc.Connect()
	users := []User{}
	if err != nil {
		return nil, err
	}

	if lc.BindDN != "" {
		err := lc.Conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return nil, err
		}
	}
	attributes := append(lc.Attributes, "dn")
	sReq := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)

	sr, err := lc.Conn.Search(sReq)
	if err != nil {
		return users, err
	}
	if len(sr.Entries) < 1 {
		return users, nil
	} else {
		for _, e := range sr.Entries {
			user := User{UserName: username}
			ua := []string{}
			for _, attr := range lc.Attributes {
				ua = append(ua, attr+":"+e.GetAttributeValue(attr))
			}
			user.Attributes = ua
			users = append(users, user)
		}
		return users, nil
	}
}

func (lc *LDAPClient) GetUserGroups(username string) ([]string, error) {
	err := lc.Connect()
	if err != nil {
		return []string{}, err
	}
	searchReq := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.GroupFilter, username),
		[]string{"cn"},
		nil,
	)
	sr, err := lc.Conn.Search(searchReq)
	if err != nil {
		return nil, err
	}

	groups := []string{}
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups, nil

}

func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}
