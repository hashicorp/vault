# How to test GCP KMS

## Test Setup

The following steps are required to run the GCP KMS wrapper tests.  This setup
assumes you have access to the internal HashiCorp Doormat CLI.  If you don't
have access to this tool, you'll need to complete that step on your own (likely
via the GCP console).

* Login to HashiCorp `doormat` and create a new temporary project with a project ID of
  `hc-vault-testing` in a location of `global`. Wait for the project to be
  created, otterbot will message you in ~10 minutes when it’s ready.
* Login to the GCP console
* Enable the GCP Cloud KMS API
* The subsequent steps I found easier via their gcloud tool (you can do it on
  the console using the “cloud shell”). Create a KMS keyring and some keys:

``` shell
gcloud kms keyrings create vault-test-keyring --location global
gcloud kms keys create --keyring vault-test-keyring --location global vault-test-key --purpose encryption
```

* Create a new service account

``` shell
gcloud iam service-accounts create go-kms-access
```

* Add the required roles to your service account

```shell
gcloud projects add-iam-policy-binding <gcp-project-id> --member="serviceAccount:go-kms-access@<gcp-project-id>.iam.gserviceaccount.com" --role="roles/cloudkms.cryptoKeyEncrypterDecrypter"
gcloud projects add-iam-policy-binding <gcp-project-id> --member="serviceAccount:go-kms-access@<gcp-project-id>.iam.gserviceaccount.com" --role="roles/cloudkms.viewer"
```

* Download a credentials key for the account (I used the console for this) and
  store them in a file in this test director names `credentials.json`. 
  
## Running tests
Now that you've completed the required setup, you can run the tests via:

```shell
❯ export GOOGLE_APPLICATION_CREDENTIALS=./credentials.json
❯ export VAULT_ACC=true
❯ go test 
```



