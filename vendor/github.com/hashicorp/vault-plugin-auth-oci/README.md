# Vault Plugin Auth OCI
Vault auth plugin for Oracle Cloud Infrastructure.

## Acceptance tests

The acceptance tests can only be run from an OCI instance.

If you are running this code on an OCI instance, you can run them directly with `make testacc`.
You will need to set the following environtment variables:
* `HOME_TENANCY_ID` to the tenancy you are running under (or your root tenancy ID)
* `ROLE_OCID_LIST` to a comma-separated list of group OCIDs to at least two groups. At least one should be a dynamic group that contains the instance, and another should be an identity group that contains your user.

For example:

```sh
make testacc HOME_TENANCY_ID=ocid1.tenancy.oc1..aaaaaaaasomecharacter ROLE_OCID_LIST=ocid1.group.oc1..aaaaaaaasomecharacters OCI_GO_SDK_DEBUG=info VAULT_LOG_LEVEL=debug
```

### Terraform

You can run the acceptance tests with terraform as well.

You will need an [OCI](https://signup.cloud.oracle.com) account.

You need to generate and download a private key in your account settings.
This should give you a private key file, the fingerprint, your tenancy OCID, and your user OCID.

Using those, you can run the acceptance tests via:

```sh
cd tests/terraform
# download your private key to this directory
terraform init
terraform apply \
  -var "fingerprint=YOURFINGERPRINT" \
  -var "tenancy_ocid=YOUR_TENANCY_OCID" \
  -var "user_ocid=YOUR_USER_OCID" \
  -var "private_key_path=YOUR_PRIVATE_KEY" \
  -var "region=YOUR_REGION"
```

This downloads the current `main` branch from GitHub and runs the tests on an OCI instance.
It takes about 5 minutes.

Don't forget to destroy the resources when you are done:

```sh
terraform destroy \
  -var "fingerprint=YOURFINGERPRINT" \
  -var "tenancy_ocid=YOUR_TENANCY_OCID" \
  -var "user_ocid=YOUR_USER_OCID" \
  -var "private_key_path=YOUR_PRIVATE_KEY" \
  -var "region=YOUR_REGION"
```
