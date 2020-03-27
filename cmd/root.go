package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var banner = `
    _____ _      _____          _____  
  / ____| |    |  __ \   /\   |  __ \ 
 | |    | |    | |  | | /  \  | |__) |
 | |    | |    | |  | |/ /\ \ |  ___/ 
 | |____| |____| |__| / ____ \| |     
  \_____|______|_____/_/    \_\_|     `

var rootCmd = &cobra.Command{
	Use:   "ej",
	Short: "LDAP Explorer",
	Long:  "LDAP Explorer that allows for LDAP schema exploration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner)
		fmt.Println("Explore LDAP Cli")
	},
}

var (
	username string
	password string
	config   string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "username for login")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "user password login")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")

	configCmd.Flags().StringVarP(&config, "config", "c", "", "specify the config yaml to load")
	configCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(loginCmd)

}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login the user to a LDAP instance",
	Long:  `Allow the user to login to the LDAP instance by asking the credentials.`,
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	u, _ := cmd.Flags().GetString("username")
	p, _ := cmd.Flags().GetString("password")

	lc := InitializeClient()
	authenticated, user, err := lc.Authenticate(u, p)
	if !authenticated {
		fmt.Printf("Error in authenticating Users: %s\n", err)
	} else {
		fmt.Printf("%v\n", user)
	}
}
