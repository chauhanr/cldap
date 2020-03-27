# cldap

A cldap is a command line utility that helps create a client for interacting with LDAP server. There
are several commands that can be used ot interact with the cldap CLI. 

1. [Config](docs/config.md)  
2. [Login](docs/login.md)


## Config File sample to save preferenaces. 

```
ldap: 
  client: 
    base: 
    bind-dn:
    bind-password:
    group-filter:
    user-filter:
    server-name: 
    host: 
    port: 
    skip-verify:
    attributes: 
     - cn
     - surname
     - givenName 
  creds: 
    username:
    password: 

```
