package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configure ldap details",
	Long:  `Configure the LDAP client with configurations like server, port, BindDN, BindPassword`,
	Run:   configureClient,
}

func configureClient(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("config")

	fmt.Printf("path to config %s\n", path)

	// load the config and save it to ~/.cldap/.cldap.yaml
}
