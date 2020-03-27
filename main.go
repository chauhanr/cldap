package main

import "github.com/chauhanr/cldap/cmd"

func main() {
	cmd.Execute()
}

/*
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", 389))
	if err != nil {
		log.Fatalf("Oops fatal error to dial to ldap server %s\n", err)
	}

	defer l.Close()

	err = l.Bind("uid=chauhr,ou=People,dc=example,dc=com", "power")
	if err != nil {
		log.Fatalf("Error loging into the server %s\n", err)
	}
	fmt.Printf("Successful binding to server\n")
*/
