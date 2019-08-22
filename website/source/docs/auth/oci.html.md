---
layout: "docs"
page_title: "OCI - Vault Auth Plugin"
sidebar_title: "OCI"
sidebar_current: "docs-auth-oci"
description: |-
  The OCI Auth plugin for HashiCorp Vault enables authentication and authorization using OCI Identity credentials.
---

# Onboarding to OCI HashiCorp Vault Auth Plugin

The OCI Auth plugin for HashiCorp Vault enables authentication and authorization using OCI Identity credentials. The open source OCI Auth plugin for Vault has been tested with version 0.11.4 of Vault and requires Go.

## Role-Based Authorization

The OCI Auth plugin authorizes using roles. A role is defined as a set of allowed [vault policies](https://www.vaultproject.io/docs/concepts/policies.html) for specific entities. When an entity such as a user or instance logs into Vault, it requests a role for that login session. The OCI Auth plugin checks whether the entity is allowed to use the role and which Vault policies are associated with that role. It then assigns the given vault policies to the Vault login request.

The goal of roles is to restrict access to only the subset of secrets that are required for that Vault session, even if the entity has access to many more secrets. This conforms to the least-privilege security model.

The roles fields are:

* ocid_List: a list of OCIDs (Group OCID, Dynamic Group OCID). Only members of these Groups or Dynamic Groups are allowed to take this role.
* token_policies: a list of Vault policies. These policies are allocated for the entity taking this role.
* token_ttl: TTL (in seconds) for which any client token assigned through this role, is valid for.

## Architectural Diagram
![Role Based Authorization](/img/oci/oci-role-based-authz.png)

There is a many-to-many relationship between various items seen above:

* A user can belong to many identity groups.
* An identity group can contain many users.
* A compute instance can belong to many dynamic groups.
* A dynamic group can contain many compute instances.
* A role defined in Vault can be mapped to many groups and dynamic groups.
* A single Vault role can be mapped to both groups and dynamic groups.
* A single Vault role can be mapped to many Vault policies.
* An Vault policy can be mapped from different Vault roles.
* An Vault policy can allow access to many secrets stored in Vault. Wildcard-based mapping is also possible for secret.
* A secret can be allowed access from many Vault policies, using both wildcards and direct mapping to a specific secret path.


## Configure the OCI Tenancy to Run Vault
While authorization is performed using vault policies, the vault must authenticate that the caller is using valid OCI credentials. In order to do this, the OCI Auth plugin must call the OCI Identity API. The OCI Auth plugin requires [instance principal](https://blogs.oracle.com/cloud-infrastructure/announcing-instance-principals-for-identity-and-access-management) credentials to call OCI Identity APIs, and therefore the vault server needs to run inside an OCI compute instance.

Follow the steps below to add policies to your tenancy that allow the OCI compute instance in which the vault server is running to call certain OCI Identity APIs.

1. In your tenancy, [launch the compute instance(s)](https://docs.cloud.oracle.com/iaas/Content/Compute/Tasks/launchinginstance.htm) that will run the vault server. The VCN in which you launch the Compute Instance should have a [Service Gateway](https://docs.cloud.oracle.com/iaas/Content/Network/Tasks/servicegateway.htm) added to it .
    * If you would like high availability for the vault service, use a [vault storage backend](https://www.vaultproject.io/docs/configuration/storage/index.html) that supports high availability, such as the OCI Storage plugin.  
        * In this case, you will be running vault in multiple compute instances.    
        * Otherwise, you will be running vault in a single compute instance.  
        
1.  Make a note of the Oracle Cloud Identifier (OCID) of the compute instance(s) running vault.
1.  In your tenancy, [create a dynamic group](https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/managingdynamicgroups.htm) with the name VaultDynamicGroup to contain the computer instance(s).
1.  Add the OCID of the compute instance(s) to the dynamic group.
1.  Add the following policies to the root compartment of your tenancy that allow the dynamic group to call specific Identity APIs.

```
    allow dynamic-group VaultDynamicGroup to {AUTHENTICATION_INSPECT} in tenancy   
    allow dynamic-group VaultDynamicGroup to {GROUP_MEMBERSHIP_INSPECT} in tenancy
```   
    
## Configure the OCI Vault Plugin
Note: Aside from the initial test, never run vault in dev mode other than in your development environment. For starting vault in your production environment, refer to [https://www.vaultproject.io/docs/configuration](https://www.vaultproject.io/docs/configuration/).

1.  Start vault in dev mode to test that it's working correctly.
    
    `vault server -dev`
        
1.  Note the root token. Never give out the root token or store in plain text form. This token should only be used for the initial configuration of the vault. For more information on tokens, see [https://www.vaultproject.io/docs/concepts/tokens.html](https://www.vaultproject.io/docs/concepts/tokens.html).
1.  Export the vault address. The code below assumes you're testing vault in the same compute instance in which the vault server is running. In production environments, always configure vault for Transport Layer Security (TLS), replace the endpoint with the actual values, and use HTTPS.
    
    `export VAULT_ADDR='[http://127.0.0.1:8200'](http://127.0.0.1:8200')`
      
1.  Enable the OCI Auth plugin.

    * Create a file named authenable.json.

    ```
    {"config":{"audit_non_hmac_request_keys":null,"audit_non_hmac_response_keys":null,"default_lease_ttl":"0s","max_lease_ttl":"0s"},"description":"","local":false,"seal_wrap":false,"type":"oci"}
    ```
    
    * Run the below command, replacing $roottoken with the root token you copied earlier.

    ```     
        curl --header "X-Vault-Token: $roottoken" --request POST \       
        --data @authenable.json \        
        http://[127.0.0.1:8200/v1/sys/auth/oci](127.0.0.1:8200/v1/sys/auth/oci)
    ```

1.  Create the vault admin policy.
    * Determine what [set of policies](https://learn.hashicorp.com/vault/identity-access-management/iam-policies) best suits your requirements.
    * Create a file named vaultadminpolicy.json. Configure the policy below according to your requirements.

    ```       
        {       
        "policy": "path \"auth/oci/*\"       
        {       
        capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]        
        }        
        path \"sys/auth/*\"       
        {       
        capabilities = [\"read\", \"list\"]       
        }       
        path \"sys/policy\"       
        {       
        capabilities = [\"read\"]        
        }
        path \"sys/policy/*\"
        {        
        capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]        
        }        
        path \"secret/*\"        
        {        
        capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]        
        }        
        path \"sys/health\"        
        {        
        capabilities = [\"read\"]        
        }"        
        }
    ```                          

    * Run the command to create the policy.

    ```        
        curl \       
        --header "X-Vault-Token: $roottoken" \       
        --request PUT \       
        --data @vaultadminpolicy.json \       
        http://127.0.0.1:8200/v1/sys/policy/vaultadminpolicy (127.0.0.1:8200/v1/sys/policy/vaultadminpolicy)
    ```

1.  Configure your home tenancy in the vault, so that only users or instances from your tenancy will be allowed to log into vault through the OCI Auth plugin.  
    * Create a file named hometenancyid.json with the below content, using the tenancy OCID. To find your tenancy OCID, see [https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm).       

        `{"home_tenancy_id":"your tenancy ocid here"}`
        
    * Configure the home_tenancy_id parameter in the vault.

    ```     
        curl --header "X-Vault-Token: $roottoken" --request PUT \       
        --data @hometenancyid.json \       
        http://127.0.0.1:8200/v1/auth/oci/config (127.0.0.1:8200/v1/auth/oci/config)
    ```
       
1.  Create a vault administrator role in the OCI Auth plugin. 
    * The vaultadminrole allows the administrator of Vault to log into Vault and grants them the permissions allowed in the policy.
    * Create a file named vaultadminrole.json with the below contents. Replace the ocid_list with the Group or Dynamic Group OCIDs in your tenancy that has users or instances that you want to take the vault admin role. 
        * For testing in dev mode, you can add the OCID of the dynamic group previously created.
        * In production, add only the OCID of groups and dynamic groups that can take the admin role in vault.
        
        `{"token_policies":"vaultadminpolicy","token_ttl":"1800","ocid_list":"ocid1.group.oc1..aaaaaaaaiqnblimpvmegkqh3bxilrdvjobr7qd223g275idcqhexamplefq,ocid1.dynamicgroup.oc1..aaaaaaaa5hmfyrdaxvmt52ekju5n7ffamn2pdvxaq6esb2vzzoduexamplea"}`
        
    * Run the following command to create the vault admin role.
            
    ```     
        curl --header "X-Vault-Token: $roottoken" --request PUT \
        --data @vaultadminrole.json \
        http://127.0.0.1:8200/v1/auth/oci/role/vaultadminrole (127.0.0.1:8200/v1/auth/oci/role/vaultadminrole)
    ```           

1.  Log into the vault using instance principal.
    * This assumes that the VAULT\_ADDR export has been specified, as shown earlier in this page.     
    * The compute instance that you are logging in from should be a part of a dynamic group that was added to the Vault admin role. The compute instance should also have connectivity to the endpoint specified in VAULT\_ADDR. 
    * When testing in dev mode in the same compute instance that the vault is running, this is [http://127.0.0.1:8200](http://127.0.0.1:8200/).    
    `vault login -method=oci auth_type=instance role=vaultadminrole`

    * You will see a response that includes a token with the previously added policy.
1.  Use the received token to read secrets, writer secrets, and add roles per the instructions in [https://www.vaultproject.io/docs/secrets/kv/kv-v1.html](https://www.vaultproject.io/docs/secrets/kv/kv-v1.html).
1.  Log into vault using the user API key.  
    *  [Add an API Key](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/apisigningkey.htm) for a user in the console. This user should be part of a group that has previously been added to the vault admin role.
    *  Create the config file `~/.oci/config` using the user's credentials as detailed in [https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm).  
        Ensure that the region in the config matches the region of the compute instance that is running vault.
    *  Log into vault using the user API key. 
    
       `vault login -method=oci auth_type=apikey role=vaultadminrole`      
1.  Stop vault and re-start it in the production environment. See [https://www.vaultproject.io/docs/configuration](https://www.vaultproject.io/docs/configuration/) for more information.      
1.  Repeat all steps in this [Configure the OCI Vault Plugin](#OnboardingtoOCIHashiCorpVaultAuthPlugin-ConfiguretheOCIVaultPlugin) section while in the production environment.

## Manage Roles in the OCI Auth Plugin

1.  Similar to creating the vault administrator role, create other roles mapped to other policies. Create a file named devrole.json with the following contents. Replace ocid_list with Groups or Dynamic Groups in your tenancy.

        `{"token_policies":"devpolicy","token_ttl":"1500","ocid_list":"ocid1.group.oc1..aaaaaaaaiqnblimpvmgrouplrdvjobr7qd223g275idcqhexamplefq,ocid1.dynamicgroup.oc1..aaaaaaaa5hmfyrdaxvmdg2u5n7ffamn2pdvxaq6esb2vzzoduexamplea"}`
    
1.  Add the role.    

    ``` 
    curl --header "X-Vault-Token: $token" --request PUT \
    --data @devrole.json \
    http://127.0.0.1:8200/v1/auth/oci/role/devrole (127.0.0.1:8200/v1/auth/oci/role/devrole)
    ```

1.  Login to vault assuming the devrole.

    `vault login -method=oci auth_type=instance role=vaultadminrole`
    
## Authentication

When authenticating, users can use vault cli.

### Via the CLI

   * With Compute Instance credentials: 
```
$ vault login -method=oci auth_type=instance role=devrole
```

   * With User credentials: [SDK Config](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm)
```
$ vault login -method=oci auth_type=apikey role=devrole
```

### Via the API

1.  First, sign the following request with your OCI credentials and obtain the signing string and the authorization header. Replace the endpoint, scheme (http or https) & role of the URL corresponding to your vault configuration. For more information on signing, see [signing the request](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/signingrequests.htm).

    `http://127.0.0.1/v1/auth/oci/role/devrole`

1.  On signing the above request, you would get headers similar to:

    ```
    The signing string would look like (line breaks inserted into the (request-target) header for easier reading):
                
    date: Thu, 05 Jan 2014 21:31:40 GMT
    (request-target): get /v1/auth/oci/role/devrole
    host: 127.0.0.1
    
    The Authorization header would look like:
    
    Signature version="1",headers="date (request-target) host",keyId="ocid1.t
    enancy.oc1..aaaaaaaaba3pv6wkcr4jqae5f15p2b2m2yt2j6rx32uzr4h25vqstifsfdsq/
    ocid1.user.oc1..aaaaaaaat5nvwcna5j6aqzjcaty5eqbb6qt2jvpkanghtgdaqedqw3ryn
    jq/73:61:a2:21:67:e0:df:be:7e:4b:93:1e:15:98:a5:b7",algorithm="rsa-sha256
    ",signature="GBas7grhyrhSKHP6AVIj/h5/Vp8bd/peM79H9Wv8kjoaCivujVXlpbKLjMPe
    DUhxkFIWtTtLBj3sUzaFj34XE6YZAHc9r2DmE4pMwOAy/kiITcZxa1oHPOeRheC0jP2dqbTll
    8fmTZVwKZOKHYPtrLJIJQHJjNvxFWeHQjMaR7M="
    ```

1.  Add the signed headers to the "request_headers" field and make the actual request to vault. An exampe is given below:

    ```
    POST http://127.0.0.1/v1/auth/oci/role/devrole
       "request_headers": {
           "date": ["Fri, 22 Aug 2019 21:02:19 GMT"],
           "(request-target)": ["get /v1/auth/oci/role/devrole"],
           "host": ["127.0.0.1"],
           "content-type": ["application/json"],
           "authorization": ["Signature algorithm=\"rsa-sha256\",headers=\"date (request-target) host\",keyId=\"dummy dummy dummy\",signature=\"dummy dummy dummy\",version=\"1\""]
       }
    ```