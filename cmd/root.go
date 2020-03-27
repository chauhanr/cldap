package cmd

import (
	"fmt"
	"os"

	model "github.com/chauhanr/cldap/models"
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
	Use:   "cldap",
	Short: "LDAP Explorer",
	Long:  "LDAP Explorer that allows for LDAP schema exploration",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("config")
		if path == "" {
			fmt.Println(banner)
			fmt.Println("Explore LDAP Cli")
		} else {
			err := configureClient(path)
			if err != nil {
				fmt.Printf("Error saving the config file: %s\n", err)
			}
		}
	},
}

var (
	username    string
	password    string
	config      string
	searchEntry string
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

	searchCmd.Flags().StringVarP(&searchEntry, "entry", "e", "", "Username to search.")
	searchCmd.MarkFlagRequired("entry")

	rootCmd.Flags().StringVarP(&config, "config", "c", "", "specify the config yaml to load")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(searchCmd)

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

	lc, err := InitializeClient(true)
	if err != nil {
		fmt.Printf("Error initializing client %s\n", err)
		return
	}
	authenticated, _, err := lc.SearchAndBind(u, p)
	if !authenticated {
		fmt.Printf("Error in authenticating Users: %s\n", err)
	} else {
		//fmt.Printf("%v\n", user)
		lc := model.LdapConfig{}
		err = lc.LoadConfig()
		if err != nil {
			fmt.Printf("Kindly save config file before login %s\n", err)
			return
		}
		if !lc.ConfigExists() {
			fmt.Printf("You can only login once the Ldap configurations are saved.\n")
			rootCmd.Help()
		}
		creds := lc.Ldap.Creds
		creds.Username = u
		creds.Password = p
		lc.Ldap.Creds.HasCreds = true
		err = lc.SaveConfig()
		if err != nil {
			fmt.Printf("Failed to save the login creds.\n")
		}
	}
}
