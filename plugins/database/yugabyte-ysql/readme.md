#   Ysql plugin for Hashicorp Vault 
##  About YugabyteDB:
YugabyteDB is a high-performance, cloud-native distributed SQL database that aims to support all PostgreSQL features. It is best to fit for cloud-native OLTP (i.e. real-time, business-critical) applications that need absolute data correctness and require at least one of the following: scalability, high tolerance to failures, or globally-distributed deployments.

### What makes YugabyteDB unique?
YugabyteDB is a transactional database that brings together 4 must-have needs of cloud native apps, namely SQL as a flexible query language, low-latency performance, continuous availability and globally-distributed scalability. Other databases do not serve all 4 of these needs simultaneously.

Monolithic SQL databases offer SQL and low-latency reads but neither have ability to tolerate failures nor can scale writes across multiple nodes, zones, regions and clouds.
Distributed NoSQL databases offer read performance, high availability and write scalability but give up on SQL features such as relational data modeling and ACID transactions.

Read more about YugabyteDB in our [Docs](https://docs.yugabyte.com/preview/faq/general/).

##  About HashiCorp Vault:
HashiCorp Vault is designed to help organizations manage access to secrets and transmit them safely within an organization. 
Secrets are defined as any form of sensitive credentials that need to be tightly controlled and monitored and can be used to unlock sensitive information. 
Secrets could be in the form of passwords, API keys, SSH keys, RSA tokens, or OTP.

### Dynamic Secrets:
A dynamic secret is generated on demand and is unique to a client, instead of a static secret, which is defined ahead of time and shared. 
Vault associates each dynamic secret with a lease and automatically destroys the credentials when the lease expires.
In this example, a client is requesting a database credential. Vault connects to the database with a private, root level credential and creates a new username and password. This new set of credentials are provided back to the client with a lease of 7 days. A week later, Vault will connect to the database with its privileged credentials and delete the newly created username.

Using Dynamic Secrets means we don’t have to be concerned about them having the shared PEM when a developer or operator leaves the organization. It also gives us a better break glass procedure should these credentials leak, as the credentials are localized to an individual resource reducing the attack vector, and the credentials are also issued with a time to live, meaning that Vault will automatically revoke them after a predetermined duration. In addition to this, by leveraging Vault Auth and Dynamic Secrets, you also gain full access logs directly tieing a SSH session to an individual user.

![ alt text for screen readers source: HashiCorp](https://www.datocms-assets.com/2885/1519774324-dynamic-secret-img-001.jpeg?fit=max&q=80&w=2500)

##  Ysql-plugin for Hashicorp Vault:
-   ysql-plugin provides APIs for using the HashiCorp Vault's Dynamic Secrets for the yugabyteDB.
-   The APIs that can be used are as follows:  
    -   Add yugabyteDB to the manage secrets i.e. enabling `write database` for yugabyteDB(ysql) while using vault.
    -   To create new users i.e. enabling `write` roles and `read` roles commands for yugabyteDB(ysql) while using vault.
    -   Mangae lease related to the yugabyteDB(ysql) i.e. enabling `lease lookup` , `lease renew` and `lease revoke` for yugabyteDB (ysql) while using vault.
-   Why seperate plugin for yugabyteDB(ysql):
   -    Yugabyte go driver can be used for connecting with the database.
        This will allow us to use the smart features, providing a high tolerance towards failures.
        

##  Before using the vault follow the below steps:
-   Make sure that the go is added to the path
```sh
$   export GOPATH=$HOME/go
$   export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
-   Clone and go to the database plugin directory
```sh 
$   git clone https://github.com/yugabyte/hashicorp-vault-ysql-plugin   && cd  hashicorp-vault-ysql-plugin  
```
-   Build the the plugin
```sh
$   go build -o <build dir>/ysql-plugin  cmd/ysql-plugin/main.go
```

-   For using the vault in the development mode add the default Vault address and Vault tocken
```sh
#   Add the VAULT_ADDR and VAULT_TOKEN
$  export VAULT_ADDR="http://localhost:8200"
$  export VAULT_TOKEN="root"
```

##  Using Vault

-   Run the server in the development mode
    -   For running the vault server in development mode `dev` flag is used.
    -   The `dev-root-tocken` informs the vault to use the default vault tocken of `root` to login.
        In case of production mode this tocken is required to be set.   
        Tocken policies are discussed [here](https://www.vaultproject.io/docs/commands/login).
    -   While running in the development mode vault will automatically register the plugin if 
        the directory of the binary of the plugin is provided as an input with the dev-plugin-dir flag as shown below.
```sh
$   vault server -dev -dev-root-token-id=root -dev-plugin-dir=<build dir> 
```

-   Enable the database's secrets:
```sh
$ vault secrets enable database
```
-   For production mode register the plugin:
```sh
$ export SHA256=$(sha256sum <build dir>/ysql-plugin  | cut -d' ' -f1)

$ vault write sys/plugins/catalog/database/ysql-plugin \
    sha256=$SHA256 \
    command="ysql-plugin"
```
-   Add the database
    -   Once can enter the credentials or use connection string:
```sh
$ vault write database/config/yugabytedb plugin_name=ysql-plugin  \
    host="127.0.0.1" \
    port=5433 \
    username="yugabyte" \
    password="yugabyte" \
    db="yugabyte" \
    allowed_roles="*"

	vault write database/config/yugabytedb \
    plugin_name=ysql-plugin \
    connection_url="postgres://{{username}}:{{password}}@localhost:5433/yugabyte?sslmode=disable" \
    allowed_roles="*" \
    username="yugabyte" \
    password="yugabyte" 
```

-   Write the role 
```sh
$ vault write database/roles/my-first-role \
    db_name=yugabytedb \
    creation_statements="CREATE ROLE \"{{username}}\" WITH PASSWORD '{{password}}' NOINHERIT LOGIN; \
       GRANT ALL ON DATABASE \"yugabyte\" TO \"{{username}}\";" \
    default_ttl="1h" \
    max_ttl="24h"
```
-   Create the user 
```sh
$   vault read database/creds/my-first-role
```

-   Lookup the details about the lease
```sh 
$  vault lease lookup  <leaseid>
```
-   Renew the lease
```sh
$  vault lease renew   <leaseid>
```    
-   Revoke the lease
```sh
$  vault lease revoke  <leaseid>
```

##  For testing:
go test can be used for testing the ysql-plugin
Use:: `go test github.com/yugabyte/hashicorp-vault-ysql-plugin`
For individual cases
-   For Initialize:
    `go test -run ^TestYsql_Initialize$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Create User:
    `go test -run ^TestYsql_NewUser$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Update User Password:
    `go test -run ^TestUpdateUser_Password$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Update User Expiration:
    `go test -run ^TestUpdateUser_Expiration$ github.com/yugabyte/hashicorp-vault-ysql-plugin`
-   For Delete User:
    `go test -run ^TestDeleteUser$ github.com/yugabyte/hashicorp-vault-ysql-plugin`

### How to use the Makefile:
-   Set the BUILD_DIR in the Makefile
-   For building the plugin, registering it and running it in development mode use `make`.
-   For enabling the plugin and creating a basic role named 'my-first-role' use `make enable`.
-   To read user  use `vault read database/creds/my-first-role`.
-   Use `make clean` to remove the build and `make test` to test the plugin. 