---
layout: "docs"
page_title: "OCI - Vault Auth Plugin"
sidebar_title: "OCI"
sidebar_current: "docs-auth-oci"
description: |-
  The OCI Auth plugin for HashiCorp (HC) Vault enables authentication and authorization using OCI Identity credentials.
---

# Onboarding to OCI HashiCorp Vault Auth Plugin

The OCI Auth plugin for HashiCorp (HC) Vault enables authentication and authorization using OCI Identity credentials. The open source OCI Auth plugin for HC Vault has been tested with version 0.11.4 of HC Vault and requires Go.

## Role-Based Authorization

The OCI Auth plugin authorizes using roles. A role is defined as a set of allowed [vault policies](https://www.vaultproject.io/docs/concepts/policies.html) for specific entities. When an entity such as a user or instance logs into HC Vault, it requests a role for that login session. The OCI Auth plugin checks whether the entity is allowed to use the role and which HC Vault policies are associated with that role. It then assigns the given vault policies to the HC Vault login request.

The goal of roles is to restrict access to only the subset of secrets that are required for that HC Vault session, even if the entity has access to many more secrets. This conforms to the least-privilege security model.

The roles fields are:

* Name: name of the role
* Description: description of the role
* Ocid_List: a list of OCIDs (Group OCID, Dynamic Group OCID). Only members of these Groups or Dynamic Groups are allowed to take this role.
* Policy_List: a list of Vault policies. These policies are allocated for the entity taking this role.
* TTL: TTL (in seconds) for which any client token assigned through this role is valid

The API takes in parameters such as:

* add_ocid_list: A comma-separated list of of group or dynamic group Oracle Cloud Identifier (OCID) that can take up a role. This list is appended to the list of existing OCIDs for that role.
* remove_ocid_list: This can be used to remove OCIDs from the list of existing OCIDs for that role.
* add_policy_list: A comma-separated list of policies granted to the entities that log in through that role. This list is appended to the list of existing policies for that role.
* remove_policy_list: This can be used to remove policies from the list of existing policies for that role.

## Architectural Diagram
![Role Based AuthZ](/img/oci/oci-role-based-authz.png)

There is a many-to-many relationship between various items seen above:

* A user can belong to many identity groups.
* An identity group can contain many users.
* A compute instance can belong to many dynamic groups.
* A dynamic group can contain many compute instances.
* A role defined in HC Vault can be mapped to many groups and dynamic groups.
* A single HC Vault role can be mapped to both groups and dynamic groups.
* A single HC Vault role can be mapped to many HC Vault policies.
* An HC Vault policy can be mapped from different HC Vault roles.
* An HC Vault policy can allow access to many secrets stored in HC Vault. Wildcard-based mapping is also possible for secret.
* A secret can be allowed access from many HC Vault policies, using both wildcards and direct mapping to a specific secret path.


## Configure the OCI Tenancy to Run HC Vault
While authorization is performed using vault policies, the vault must authenticate that the caller is using valid OCI credentials. In order to do this, the OCI Auth plugin must call the OCI Identity API. The OCI Auth plugin requires [instance principal](https://blogs.oracle.com/cloud-infrastructure/announcing-instance-principals-for-identity-and-access-management) credentials to call OCI Identity APIs, and therefore the vault server needs to run inside an OCI compute instance.

Follow the steps below to add policies to your tenancy that allow the OCI compute instance in which the vault server is running to call certain OCI Identity APIs.

1.  In your tenancy, [launch the compute instance(s)](https://docs.cloud.oracle.com/iaas/Content/Compute/Tasks/launchinginstance.htm) that will run the vault server.
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
        http://[127.0.0.1:8200/v1/sys/policy/vaultadminpolicy](127.0.0.1:8200/v1/sys/policy/vaultadminpolicy)
    ```

1.  Configure your home tenancy in the vault, so that only users or instances from your tenancy will be allowed to log into vault through the OCI Auth plugin.  
    * Create a file named hometenancyid.json with the below content, using the tenancy OCID. To find your tenancy OCID, see [https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm).       

        `{"configValue":"your tenancy ocid here"}`
        
    * Configure the homeTenancy parameter in the vault.
       
    ```     curl --header "X-Vault-Token: $roottoken" --request PUT \       
        --data @hometenancyid.json \       
        http://[127.0.0.1:8200/v1/auth/oci/config/homeTenancyId](127.0.0.1:8200/v1/auth/oci/config/homeTenancyId)
    ```
       
    * Create a vault administrator role in the OCI Auth plugin. This role allows the administrator of the vault to log into vault and grants them the permissions allowed in the policy.
    * Create a file named vaultadminrole.json with the below contents. The TTL specifies how long, in seconds, tokens for this role are valid.        
        `{"add_policy_list":"vaultadminpolicy","description":"VaultAdminRole","ttl":"1500"}`
        
    * Run the command to create the vault admin role.
            
    ```     
        curl --header "X-Vault-Token: $roottoken" --request PUT \
        --data @vaultadminrole.json \
        http://[127.0.0.1:8200/v1/auth/oci/role/vaultadminrole](127.0.0.1:8200/v1/auth/oci/role/vaultadminrole)
    ``` 
          
1.  Add the group or dynamic group OCID corresponding to users or instances to the vault admin role.
    * Create a file named vaultadminids.json with the below content. The OCID list is a comma-separated list of group and dynamic group OCIDs in your tenancy.       
        * For a test in dev mode, you can add the OCID of the dynamic group previously created.
        * In production, add only the OCID of groups and dynamic groups that can take the admin role in vault.
        
        `{"add_ocid_list":"[ocid1.group.oc1..dummy1,ocid1.dynamicgroup.oc1..dummy1](ocid1.group.oc1..dummy1,ocid1.dynamicgroup.oc1..dummy1)"}`
        
    *  Run the command to add the admin IDs to the vault admin role.        

    ```     
        curl --header "X-Vault-Token: $roottoken" --request PUT \       
        --data @vaultadminids.json \  
        http://[127.0.0.1:8200/v1/auth/oci/role/vaultadminrole](127.0.0.1:8200/v1/auth/oci/role/vaultadminrole)
    ```            

1.  Log into the vault using instance principal.
    * This assumes that the VAULT\_ADDR export has been specified, as shown earlier in this page.     
    * The compute instance that you are logging in from should be a part of a dynamic group that was added to the Vault admin role. The compute instance should also have connectivity to the endpoint specified in VAULT\_ADDR. 
    * When testing in dev mode in the same compute instance that the vault is running, this is [http://127.0.0.1:8200](http://127.0.0.1:8200/).    
    `vault login -method=oci authType=InstancePrincipal role=vaultadminrole`

    * You will see a response that includes a token with the previously added policy.
1.  Use the received token to read secrets, writer secrets, and add roles per the instructions in [https://www.vaultproject.io/docs/secrets/kv/kv-v1.html](https://www.vaultproject.io/docs/secrets/kv/kv-v1.html).
1.  Log into vault using the user API key.  
    *  [Add an API Key](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/apisigningkey.htm) for a user in the console. This user should be part of a group that has previously been added to the vault admin role.
    *  Create the config file `~/.oci/config` using the user's credentials as detailed in [https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm).  
        Ensure that the region in the config matches the region of the compute instance that is running vault.
    *  Log into vault using the user API key. `vault login -method=oci authType=ApiKey role=vaultadminrole`      
1.  Stop vault and re-start it in the production environment. See [https://www.vaultproject.io/docs/configuration](https://www.vaultproject.io/docs/configuration/) for more information.      
1.  Repeat all steps in this [Configure the OCI Vault Plugin](#OnboardingtoOCIHashiCorpVaultAuthPlugin-ConfiguretheOCIVaultPlugin) section while in the production environment.

## Manage Roles in the OCI Auth Plugin

1.  Create a file named devrole.json and replace the content below according to your requirements.    
    `{"add_policy_list":"databasepolicy","description":"VaultAdminRole","ttl":"1800","add_ocid_list":"[ocid1.group.oc1..group2,ocid1.dynamicgroup.oc1..group3](ocid1.group.oc1..group2,ocid1.dynamicgroup.oc1..group3)","remove_ocid_list":"[ocid1.group.oc1..dummy1,ocid1.dynamicgroup.oc1..group1](ocid1.group.oc1..dummy1,ocid1.dynamicgroup.oc1..group1)"}`
    
1.  Apply the changes to the role.    

    ``` 
    curl --header "X-Vault-Token: $token" --request PUT \
    --data @devrole.json \
    http://[127.0.0.1:8200/v1/auth/oci/role/devrole](127.0.0.1:8200/v1/auth/oci/role/devrole)
    ```
