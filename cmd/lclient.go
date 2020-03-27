package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"

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

func InitializeClient() *LDAPClient {
	c := LDAPClient{
		Base: "dc=example,dc=com",
		//BindDN: "cn=admin,dc=example,dc=com",
		BindDN: "cn=ldapadm,dc=cfee,cf=apps,dc=com",
		//BindPassword: "schumi11",
		BindPassword: "ldappassword",
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberid=%s)",
		//Host:        "localhost",
		Host:               "169.48.114.206",
		Port:               389,
		UseSSL:             false,
		InsecureSkipVerify: false,
		Attributes:         []string{"cn", "uid", "givenName", "sn"},
	}

	return &c
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

func (lc *LDAPClient) Authenticate(username, password string) (bool, map[string]string, error) {
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

func (lc *LDAPClient) Close() {
	if lc.Conn != nil {
		lc.Conn.Close()
		lc.Conn = nil
	}
}
