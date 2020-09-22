# Auto renew token leases
The maximum lease for a token other than a root token is 755 hours i.e. approximately 31 days. In case one wants to setup the tokens to be renewed without any intervention, a kubernetes cron job can be used.

# How to use it
- Build the Dockerfile
- Give appropriate env vars to the cron-job.yaml file

# How it works
- Dockerfile's entrypoint runs a script (script.sh)
- The script makes an API call to vault server with the payload information like the token whose lease is to be renewed and the duration for which the lease has to be extended.

