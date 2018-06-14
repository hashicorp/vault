---
layout: "guides"
page_title: "Java Application Demo - Guides"
sidebar_current: "guides-encryption-spring-demo"
description: |-
  This guide discusses the concepts necessary to help users
  understand Vault's AppRole authentication pattern and how to use it to
  securely introduce a Vault authentication token to a target server,
  application, container, etc. in a Java environment.
---

# Java Sample App using Spring Cloud Vault

Once you have learned the fundamentals of Vault, the next step is to start
integrating your system with Vault to secure your organization's secrets.  

The purpose of this guide is to go through the working implementation demo
introduced in the [Manage secrets, access, and encryption in the public cloud
with
Vault](https://www.hashicorp.com/resources/solutions-engineering-webinar-series-episode-2-vault)
webinar.

[![YouTube](/assets/images/vault-java-demo-1.png)](https://youtu.be/NxL2-XuZ3kc)

The Java application in this demo leverages the [_Spring Cloud
Vault_](https://cloud.spring.io/spring-cloud-vault/) library which provides
lightweight client-side support for connecting to Vault in a distributed
environment.


## Reference Material

- [Manage secrets, access, and encryption in
the public cloud with
Vault](https://www.hashicorp.com/resources/solutions-engineering-webinar-series-episode-2-vault)
- [Spring Cloud Vault](https://cloud.spring.io/spring-cloud-vault/)
- [Transit Secrets Engine](/docs/secrets/transit/index.html)
- [Secrets as a Service: Dynamic Secrets](/guides/secret-mgmt/dynamic-secrets.html)


## Estimated Time to Complete

15 minutes


## Challenge

Incidents of data breaches which expose sensitive information make headlines more
often than we like to hear. It becomes more and more important to protect data
by encrypting it whether the data is in-transit or at-rest. However, creating
a highly secure and sophisticated solution by yourself requires time and resources
which are in demand when an organization is facing a constant threat.


## Solution

Vault centralizes management of cryptographic services used to protect your
data.  Your system can communicate with Vault easily through the Vault API to
encrypt and decrypt your data, and the encryption keys never have to leave the
Vault.

![Encryption as a Service](/assets/images/vault-eaas.png)


## Prerequisites

To perform the tasks described in this guide:

- Install [HashiCorp Vagrant](https://www.vagrantup.com/intro/getting-started/install.html)
- Clone or download the demo assets from the [hashicorp/vault-guides](https://github.com/hashicorp/vault-guides/tree/master/secrets/spring-cloud-vault)
GitHub repository


## Steps

-> For the purposes of this guide, you are going to provision a Linux machine
locally using Vagrant.  However, the GitHub repository provides supporting files
to provision the environment demonstrated in the webinar.

After downloading the [demo assets](#prerequisites) from the GitHub repository,
you should find the following folders:

| Folder           | Description                                               |
|------------------|-----------------------------------------------------------|
| `aws`            | Supporting files to deploy the demo app to AWS            |
| `kubernetes`     | Supporting files to deploy the demo app to Kubernetes     |
| `nomad`          | Supporting files to deploy the demo app to Nomad          |
| `scripts`        | Scripts to setup PostgreSQL and Vault                     |
| `src/main`       | Sample app source code                                    |
| `vagrant-local`  | Vagrant file to deploy the demo locally                   |


<br>

In this guide, you will perform the following:

1. [Review the demo application implementation](#step1)
1. [Deploy and review the demo environment](#step2)
1. [Run the demo application](#step3)
1. [Reload the Static Secrets](#step4)

![Encryption as a Service](/assets/images/vault-java-demo-10.png)

### <a name="step1"></a>Step 1: Review the demo application implementation

The source code can be found under the `src/main` directory.

```plaintext
├── java
│   └── com
│       └── hashicorp
│           └── vault
│               └── spring
│                   └── demo
│                       ├── BeanUtil.java
│                       ├── Order.java
│                       ├── OrderAPIController.java
│                       ├── OrderRepository.java
│                       ├── Secret.java
│                       ├── SecretController.java
│                       ├── TransitConverter.java
│                       └── VaultDemoOrderServiceApplication.java
└── resources
    └── application.yaml
```

The demo Java application leverages the Spring Cloud Vault library to
communicate with Vault.

In the `TransitConverter` class, the `convertToDatabaseColumn` method invokes a
Vault operation to encrypt the `order`. Similarly, the
`convertToEntityAttribute` method decrypts the `order` data.

```plaintext
...
	@Override
	public String convertToDatabaseColumn(String customer) {
		VaultOperations vaultOps = BeanUtil.getBean(VaultOperations.class);
		Plaintext plaintext = Plaintext.of(customer);
		String cipherText = vaultOps.opsForTransit().encrypt("order", plaintext).getCiphertext();
		return cipherText;
	}

	@Override
	public String convertToEntityAttribute(String customer) {
		VaultOperations vaultOps = BeanUtil.getBean(VaultOperations.class);
		Ciphertext ciphertext = Ciphertext.of(customer);
        String plaintext = vaultOps.opsForTransit().decrypt("order", ciphertext).asString();
		return plaintext;
...
```

The `VaultDemoOrderServiceApplication` class defines the `main` method.

```plaintext
public class VaultDemoOrderServiceApplication  {

	private static final Logger logger = LoggerFactory.getLogger(VaultDemoOrderServiceApplication.class);

	@Autowired
    	private SessionManager sessionManager;

	@Value("${spring.datasource.username}")
	private String dbUser;

	@Value("${spring.datasource.password}")
	private String dbPass;

	public static void main(String[] args) {
		SpringApplication.run(VaultDemoOrderServiceApplication.class, args);
	}

	@PostConstruct
	public void initIt() throws Exception {
		logger.info("Got Vault Token: " + sessionManager.getSessionToken().getToken());
		logger.info("Got DB User: " + dbUser);
	}
}
```

The `OrderAPIController` class defines the API endpoint (`api/orders`).


### <a name="step2"></a>Step 2: Deploy and review the demo environment

Now let's run the demo app and examine how it behaves.

~> To keep it simple and lightweight, you are going to run a Linux virtual
machine locally using Vagrant.

#### Task 1: Run Vagrant

In the **`vault-guides/secrets/spring-cloud-vault/vagrant-local`** folder,
a `Vagrantfile` is provided which spins up a Linux machine where the demo
components are installed and configured.

```shell
# Change the working directory to vagrant-local
$ cd /vault-guides/secrets/spring-cloud-vault/vagrant-local

# Create and configure a Linux machine. This takes about 3 minutes
$ vagrant up
...
demo: Success! Data written to: database/roles/order
demo: Success! Enabled the transit secrets engine at: transit/
demo: Success! Data written to: transit/keys/order
demo: Success! Data written to: secret/spring-vault-demo

# Verify that the virtual machine was successfully created and running
$ vagrant status
Current machine states:
demo                      running (virtualbox)
...

# Connect to the demo machine
$ vagrant ssh demo
```


There are 3 Docker containers running on the machine: `spring`, `vault`, and `postgres`.

```plaintext
[vagrant@demo ~]$ docker ps
CONTAINER ID     IMAGE            COMMAND                  CREATED           STATUS           NAMES
684d8fb23ae5     spring           "java -Djava.secur..."   7 minutes ago     Up 7 minutes     spring
dc6a3454b323     vault:0.10.0     "docker-entrypoint..."   7 minutes ago     Up 7 minutes     vault
4093a45c209f     postgres         "docker-entrypoint..."   7 minutes ago     Up 7 minutes     postgres
```

#### Task 2: Examine the Vault environment

During the demo machine provisioning, the `/scripts/vault.sh` script was
executed to perform the following:

- Created a policy named, **`order`**
- Enabled the `transit` secret engine and created an encryption key named, **`order`**
- Enabled the `database` secret engine and created a role named, **`order`**

View the `vault` log:

```plaintext
[vagrant@demo ~]$  docker logs vault
...
==> Vault shutdown triggered
==> Vault server configuration:
             Api Address: http://0.0.0.0:8200
                     Cgo: disabled
         Cluster Address: https://0.0.0.0:8201
              Listener 1: tcp (addr: "0.0.0.0:8200", cluster address: "0.0.0.0:8201", tls: "disabled")
               Log Level: info
                   Mlock: supported: true, enabled: false
                 Storage: inmem
                 Version: Vault v0.10.0
             Version Sha: 5dd7f25f5c4b541f2da62d70075b6f82771a650d
WARNING! dev mode is enabled! In this mode, Vault runs entirely in-memory
and starts unsealed with a single unseal key. The root token is already
authenticated to the CLI, so you can immediately begin using Vault.
You may need to set the following environment variable:
    $ export VAULT_ADDR='http://0.0.0.0:8200'
The unseal key and root token are displayed below in case you want to
seal/unseal the Vault or re-authenticate.
Unseal Key: 2QIPWPDykRG/xWWl0quSHiXq8u+pFg3yEq0sgJPhMbA=
Root Token: root
...
```

Notice that the log indicates that the Vault server is running in the `dev`
mode, and the root token is `root`.  

You can visit the Vault UI at http://localhost:8200/ui.  Enter **`root`** and
click **Sign In**.


Select the **`transit/`** secrets engine, and you should find an encryption key
named, "`order`".

![Vault UI](/assets/images/vault-java-demo-2.png)

Under the **Policies**, verify that the `order` policy exists.

![Vault UI](/assets/images/vault-java-demo-3.png)

This `order` policy is for the application.  It permits `read` on the
`database/creds/order` path so that the demo app can get a dynamically generated
database credential from Vault. Therefore, the PostgreSQL credentials are not
hard-coded anywhere.

```plaintext
path "database/creds/order"
{
  capabilities = ["read"]
}
```

An `update` permission allows the app to request data encryption and decryption
using the `order` encryption key in Vault.

```plaintext
...
path "transit/decrypt/order" {
  capabilities = ["update"]
}

path "transit/encrypt/order" {
  capabilities = ["update"]
}
...
```



#### Task 3: Examine the Spring container

Remember that the `VaultDemoOrderServiceApplication` class logs messages during
the successful execution of `initIt()`:

```plaintext
...
  @PostConstruct
  	public void initIt() throws Exception {
  		logger.info("Got Vault Token: " + sessionManager.getSessionToken().getToken());
  		logger.info("Got DB User: " + dbUser);
...
```

Verify that the log indicates that the demo app obtained a database
credentials from Vault successfully:

```plaintext
[vagrant@demo ~]$  docker logs spring | grep Got
...VaultDemoOrderServiceApplication : Got Vault Token: root
...VaultDemoOrderServiceApplication : Got DB User: v-token-order-rywqz61432yyx2x27w8r-1524067226
```

Create a new shell session in the `spring` container.

```plaintext
[vagrant@demo ~]$  docker exec -it spring sh
/ #
```

Find the `bootstrap.yaml` file:

```plaintext
/ #  ls -al
total 36720
drwxr-xr-x    1 root     root            51 Apr 18 16:00 .
drwxr-xr-x    1 root     root            51 Apr 18 16:00 ..
-rwxr-xr-x    1 root     root             0 Apr 18 16:00 .dockerenv
-rwxr--r--    1 root     root      37587245 Apr 18 15:59 app.jar
drwxr-xr-x    2 root     root          4096 Jan  9 19:37 bin
-rw-r--r--    1 1000     1000           426 Apr 17 17:58 bootstrap.yaml
...

/ # cat bootstrap.yaml
spring.application.name: spring-vault-demo
spring.cloud.vault:
    authentication: TOKEN
    token: ${VAULT_TOKEN}
    host: localhost
    port: 8200
    scheme: http
    fail-fast: true
    config.lifecycle.enabled: true
    generic:
      enabled: true
      backend: secret
    database:
        enabled: true
        role: order
        backend: database
spring.datasource:
  url: jdbc:postgresql://localhost:5432/postgres
```

The client token was injected into the `spring` container as an environment
variable (`VAULT_TOKEN`) by Vagrant.

Enter `exit` to close the shell session in the `spring` container.


#### Task 4: Examine the PostgreSQL database

Connect to the PostgreSQL database running in the `postgres` container:

```plaintext
[vagrant@demo ~]$ docker exec -it postgres psql -U postgres -d postgres
psql (10.3 (Debian 10.3-1.pgdg90+1))
Type "help" for help.


postgres=# \d orders
                                          Table "public.orders"
    Column     |            Type             | Collation | Nullable |              Default
---------------+-----------------------------+-----------+----------+------------------------------------
 id            | bigint                      |           | not null | nextval('orders_id_seq'::regclass)
 customer_name | character varying(60)       |           | not null |
 product_name  | character varying(20)       |           | not null |
 order_date    | timestamp without time zone |           | not null |
Indexes:
    "orders_pkey" PRIMARY KEY, btree (id)
```

Let's list the existing database roles.

```plaintext
postgres-# \du
                                                     List of roles
                   Role name                   |                         Attributes                         | Member of
-----------------------------------------------+------------------------------------------------------------+-----------
 postgres                                      | Superuser, Create role, Create DB, Replication, Bypass RLS | {}
 v-token-order-rywqz61432yyx2x27w8r-1524067226 | Password valid until 2018-04-18 20:56:31+00                | {}
```

Notice that there is a role name starting with `v-token-order-` which was
dynamically created by the database secret engine.

~> **NOTE:** To learn more
about the database secret engine, read the [Secrets as a Service: Dynamic
Secrets](/guides/secret-mgmt/dynamic-secrets.html) guide.

Enter `\q` to exit out of the `psql` session, or you can open another terminal
and SSH into the demo virtual machine.



### <a name="step3"></a>Step 3: Run the demo application

If everything looked fine in [Step 2](#step2), you are ready to write some data.

![Vault UI](/assets/images/vault-java-demo-9.png)

You have [verified in the `spring` log](#task-3-examine-the-sprig-container)
that the demo app successfully retrieved a database credential from the Vault
server during its initialization.

The next step is to send a new order request via the demo app's _orders_ API
(http://localhost:8080/api/orders).  

```shell
# Create a new order data
[vagrant@demo ~]$ tee payload.json<<EOF
{
  "customerName": "John",
  "productName": "Nomad"
}
EOF

# Send a request using cURL
[vagrant@demo ~]$ curl --request POST --header "Content-Type: application/json" \
                       --data @payload.json http://localhost:8080/api/orders | jq
{
  "id": 2,
  "customerName": "John",
  "productName": "Nomad",
  "orderDate": "2018-04-18T22:07:42.916+0000"
}
```

**NOTE:** Alternatively, you can use tool such as
[Postman](https://www.getpostman.com/apps) instead of cURL to invoke the API if
you prefer.

![Postman](/assets/images/vault-java-demo-4.png)


The order data you sent gets encrypted by Vault. The database only sees the
ciphertext. Let's verify that the order information stored in the database
is encrypted.

```plaintext
[vagrant@demo ~]$ docker exec -it postgres psql -U postgres -d postgres

postgres=# select * from orders;
 id |                     customer_name                     | product_name |       order_date
----+-------------------------------------------------------+--------------+-------------------------
  1 | vault:v1:Qj0lx5DSZvwcHeMOX/5UX/ErHTaDPA3mVlSSpaXd1tbM | VE           | 2018-04-18 21:56:37.924
  2 | vault:v1:UwL3HnyqTUac5ElS5WYAuNg3NdIMFtd6vvwukL+FaKun | Nomad        | 2018-04-18 22:07:42.916
(2 rows)

postgres=# \q
```

>In this demo, Vault encrypts the customer names; therefore, the values in the
**`customer_name`** column do not display the names in a human readable
manner (e.g. "James" and "John").  

Now, retrieve the order data via the orders API:

```plaintext
[vagrant@demo ~]$ curl --header "Content-Type: application/json" \
                       http://localhost:8080/api/orders | jq
[
 {
   "id": 1,
   "customerName": "James",
   "productName": "VE",
   "orderDate": "2018-04-18T21:56:37.924+0000"
 },
 {
   "id": 2,
   "customerName": "John",
   "productName": "Nomad",
   "orderDate": "2018-04-18T22:07:42.916+0000"
 }
]
```

The customer names should be readable. Remember that the **`order`** policy
permits the demo app to encrypt and decrypt data using the `order` encryption
key in Vault.

#### Web UI

Vault UI makes it easy to decrypt the data.

In the **Secrets** tab, select **`transit/` > `orders`**, and select **Key
actions**.

![Web UI](/assets/images/vault-java-demo-5.png)

Select **Decrypt** from the transit actions.  Now, copy the ciphertext from the
**`orders`** table and paste it in.

![Web UI](/assets/images/vault-java-demo-6.png)

Click **Decrypt**.

![Web UI](/assets/images/vault-java-demo-7.png)

Finally, click **Decode from base64** to reveal the customer name.

![Web UI](/assets/images/vault-java-demo-8.png)


### <a name="step4"></a>Step 4: Reloading the Static Secrets

Now, let's test another API endpoint, **`api/secret`** provided by the demo app.
A plain old Java object, `Secret` defines a get method for `key` and `value`.
The `SecretController.java` defines an API endpoint, **`api/secret`**.  

```plaintext
package com.hashicorp.vault.spring.demo;
...

@RefreshScope
@RestController
public class SecretController {

	@Value("${secret:n/a}")
	String secret;

	@RequestMapping("/api/secret")
	public Secret secret() {
		return new Secret("secret", secret);
	}
}
```

Remember from [Step 2](#task-2-examine-the-vault-environment) that the
**`order`** policy granted permissions on the `secret/spring-vault-demo` path.

```plaintext
path "secret/spring-vault-demo" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
...
```
<br>

The demo app retrieved the secret from `secret/spring-vault-demo` and has a
local copy.  If someone (or perhaps another app) updates the secret, it makes the
secret held by the demo app to be obsolete.

![Static Secret](/assets/images/vault-java-demo-11.png)

Spring offers [Spring Boot
Actuator](https://docs.spring.io/spring-boot/docs/current/reference/htmlsingle/#production-ready)
which can be used to facilitate the reloading of the static secret.

#### Task 1: Read the secret

The initial key-value was set by Vagrant during the provisioning. (See the
`Vagrantfile` at line 48.)

Let's invoke the demo app's secret API (**`api/secret`**):

```plaintext
$ curl -s http://localhost:8080/api/secret | jq
{
  "key": "secret",
  "value": "hello-vault"
}
```

This is the secret that the demo app knows about.


#### Task 2: Update the Secrets

Now, update the secret stored in Vault using API:

```shell
# Update the value via API
$ curl --header "X-Vault-Token: root" \
       --request POST \
       --data '{ "secret": "my-api-key" }' \
       http://127.0.0.1:8200/v1/secret/spring-vault-demo

# Verify that the secret value was updated
$ curl --header "X-Vault-Token: root" \
       http://127.0.0.1:8200/v1/secret/spring-vault-demo | jq
{
 "request_id": "514601e4-a790-3dc6-14b0-537d6982a6c6",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 2764800,
 "data": {
   "secret": "my-api-key"
 },
...
}
```


#### Task 3: Refresh the secret on demo app

Run the demo app's secret API again:

```plaintext
$ curl -s http://localhost:8080/api/secret | jq
{
  "key": "secret",
  "value": "hello-vault"
}
```

The current value stored in Vault is now `my-api-key`; however, the demo app
still holds `hello-vault`.  


Spring provides an
[actuator](https://docs.spring.io/spring-boot/docs/current/reference/htmlsingle/#production-ready)
which can be leveraged to refresh the secret value.  At line 54 of the
`vault-guides/secrets/spring-cloud-vault/pom.xml`, you see that the actuator was
added to the project.

```plaintext
...
<dependency>
  <groupId>org.springframework.boot</groupId>
  <artifactId>spring-boot-starter-actuator</artifactId>
</dependency>
...
```

Let's refresh the secret using the actuator:

```plaintext
$ curl -s --request POST http://localhost:8080/actuator/refresh | jq
[
  "secret"
]
```

Read back the secret from the demo app again:

```plaintext
$ curl -s http://localhost:8080/api/secret | jq
{
  "key": "secret",
  "value": "my-api-key"
}
```

It should display the correct value.

---

When you are done exploring the demo implementation, you can destroy the virtual
machine:

```plaintext
$ vagrant destroy
demo: Are you sure you want to destroy the 'demo' VM? [y/N] y
==> demo: Forcing shutdown of VM...
==> demo: Destroying VM and associated drives...
```

~> In the webinar, the demo environment was running in a public cloud, and Nomad
and Consul were also installed and configured.  If you wish to build a similar
environment using Kubernetes, the assets in the `vault-guides/secrets/spring-cloud-vault/kubernetes`
folder provides you with some guidance.

## Next steps

[AppRole](/docs/auth/approle.html) is an authentication mechanism within Vault
to allow machines or apps to acquire a token to interact with Vault. Read the
[AppRole Pull Authentication](/guides/identity/authentication.html) guide
which introduces the steps to generate tokens for machines or apps by enabling
AppRole auth method.
