package cmd

import (
	"fmt"

	model "github.com/chauhanr/cldap/models"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search the user to a LDAP instance",
	Long: `The search command will allow for user to serach the user
	        This command assumes that you are already logged into the system.
		if not then it will ask you to login.`,
	Run: searchEntities,
}

func searchEntities(cmd *cobra.Command, args []string) {
	u, _ := cmd.Flags().GetString("entry")

	lc := model.LdapConfig{}
	err := lc.LoadConfig()
	if err != nil {
		fmt.Printf("Ldap client not configured. Kindly load configurationi\n")
		rootCmd.Help()
		return
	}
	if !lc.Ldap.Client.HasConfig {
		fmt.Printf("Configuraitons for Ldap are missing. Run the --config option to load config\n")
		rootCmd.Help()
		return
	}
	if !lc.Ldap.Creds.HasCreds {
		fmt.Printf("You have not saved your user creds. Use the Login command to save your creds\n")
		loginCmd.Help()
		return
	}

	c, err := InitializeClient(false)
	if err != nil {
		fmt.Printf("Unable to initialize ldap Client %s\n", err)
		return
	}
	users, err := c.SearchUser(u)

	fmt.Printf("%v\n", users)
}
