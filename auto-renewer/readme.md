# Auto renew token leases
- The maximum lease for a token other than a root token is 755 hours i.e. approximately 31 days. 
- In case one wants to setup the tokens to be renewed without any intervention, like say for a production env, a kubernetes cron job can be used.

# How to use it
- Build the Dockerfile
- Give appropriate env vars to the cron-job.yaml file

# How it works
- Dockerfile's entrypoint runs a script (script.sh)
- The script makes an API call to vault server with the payload information like the token whose lease is to be renewed and the duration for which the lease has to be extended.

# Info about the env vars of the Docker image
- RENEW_TOKEN : The token which has to be renewed.
- INCREMENT_VALUE : The duration of time by which the lease has to be extended. e.g. 3h, 755h.
- ROOT_TOKEN : The root token of the vault, this is required by the vault to authorize the API request.
- URL : The url at which vault can be accessed by the Docker image inside the cronjob. This is required to make the API call for lease renewal.

