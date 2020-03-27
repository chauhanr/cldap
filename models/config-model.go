package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	CLDAP_CONFIG_HOME = ".cldap"
	CLDAP_CONFIG_FILE = "cldap-conf.yaml"
)

type LdapConfig struct {
	Ldap Ldap `yaml:"ldap"`
}

type Ldap struct {
	Client Client `yaml:"client"`
	Creds  Creds  `yaml:"creds"`
}

type Client struct {
	HasConfig          bool     `yaml:"configured"`
	Base               string   `yaml:"base"`
	BindDN             string   `yaml:"bind-dn"`
	BindPassword       string   `yaml:"bind-password"`
	GroupFilter        string   `yaml:"group-filter"`
	UserFilter         string   `yaml:"user-filter"`
	Host               string   `yaml:"host"`
	Port               int      `yaml:"port"`
	UseSSL             bool     `yaml:"usessl"`
	ServerName         string   `yaml:"server-name"`
	InsecureSkipVerify bool     `yaml:"skip-verify"`
	Attributes         []string `yaml:"attributes"`
}

type Creds struct {
	HasCreds bool   `yaml:"has-creds"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (lc *LdapConfig) SaveConfig() error {
	ld := getConfigHomeDirPath()
	if _, err := os.Stat(ld); !os.IsExist(err) {
		os.MkdirAll(ld, os.ModePerm)
	}
	l := getConfigFilePath()
	file, err := os.OpenFile(l, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("Error reading file: " + l)
	}
	valid, err := validate(lc)
	if !valid {
		fmt.Printf("Config file has errors %s", err)
		return err
	}
	lc.Ldap.Client.HasConfig = true

	defer file.Close()
	yencoder := yaml.NewEncoder(file)
	err = yencoder.Encode(lc)
	if err != nil {
		return errors.New("Error encoding file: " + l + " error: " + err.Error())
	}
	return nil
}

func validate(lc *LdapConfig) (bool, error) {
	invaildFields := []string{}
	isValid := true
	if lc.Ldap.Client.Base == "" {
		isValid = false
		invaildFields = append(invaildFields, "Base")
	}
	if lc.Ldap.Client.BindDN == "" {
		isValid = false
		invaildFields = append(invaildFields, "BindDN")
	}
	if lc.Ldap.Client.BindPassword == "" {
		isValid = false
		invaildFields = append(invaildFields, "BindPassword")
	}
	if lc.Ldap.Client.GroupFilter == "" {
		isValid = false
		invaildFields = append(invaildFields, "GroupFilter")
	}
	if lc.Ldap.Client.UserFilter == "" {
		isValid = false
		invaildFields = append(invaildFields, "UserFilter")
	}

	if lc.Ldap.Client.ServerName == "" {
		if lc.Ldap.Client.Host == "" {
			isValid = false
			invaildFields = append(invaildFields, "Host")
		}
	}
	message := fmt.Sprintf("Following fields are mandatory %v\n", invaildFields)
	return isValid, errors.New(message)
}

func (lc *LdapConfig) SaveCreds() error {
	lc.Ldap.Creds.HasCreds = true
	err := lc.SaveConfig()
	return err
}

func (lc *LdapConfig) LoadConfig() error {
	l := getConfigFilePath()
	file, err := os.Open(l)
	defer file.Close()
	if err != nil {
		return err
	}
	d := yaml.NewDecoder(file)
	err = d.Decode(lc)
	if err != nil {
		return err
	}
	return nil
}

func (lc *LdapConfig) ConfigExists() bool {
	l := getConfigFilePath()
	if _, err := os.Stat(l); err == nil {
		return true
	} else {
		return false
	}
}

func (lc *LdapConfig) CredsExist() bool {
	if !lc.ConfigExists() {
		return false
	} else {
		e := lc.LoadConfig()
		if e != nil {
			return false
		}
		return lc.Ldap.Creds.HasCreds
	}
}

func (lc *LdapConfig) CleanConfig() error {
	if lc.ConfigExists() {
		p := getConfigFilePath()
		err := os.Remove(p)
		if err != nil {
			fmt.Printf("Error cleaning config file %s, Error: %s\n", p, err)
			return err
		}
	} else {
		// do nothing
	}
	return nil
}

func getConfigHomeDirPath() string {
	home, _ := os.UserHomeDir()
	conf := filepath.Join(home, CLDAP_CONFIG_HOME)
	return conf
}

func getConfigFilePath() string {
	home, _ := os.UserHomeDir()
	conf := filepath.Join(home, CLDAP_CONFIG_HOME, CLDAP_CONFIG_FILE)
	return conf
}
