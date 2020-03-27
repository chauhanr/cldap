package cmd

import (
	"os"

	model "github.com/chauhanr/cldap/models"
	"gopkg.in/yaml.v2"
)

/**
  1. Check of the file mentioned is present or not.
  2. Once file is present we need to unmarshal the file
  3. save the file tp ~/.clap/cldap-conf.yml
*/

func configureClient(path string) error {
	if _, err := os.Stat(path); os.IsExist(err) {
		return err
	}
	// reading the file for config details
	lc := model.LdapConfig{}
	lc.Ldap.Creds.HasCreds = false
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	lc.Ldap.Client.HasConfig = true
	d := yaml.NewDecoder(file)
	err = d.Decode(&lc)
	if err != nil {
		return err
	}

	// finally save the file
	err = lc.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}
