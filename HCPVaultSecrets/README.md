# Introduction

This project is a driver that fetches secrets from HashiCorp Vault Secrets.

# Setup

### Create a service principal

- Go to your project or organization -> Access control (IAM) -> Service principals
- Create a service principal and save the client id and secret

### Create an env file

Create `.env` file with the following variables:

- HCP_CLIENT_ID : service principal's client id
- HCP_CLIENT_SECRET : service principal's client secret
- HCP_ORG_ID : organization's id
- HCP_PROJ_ID : project's id

## Run

Run main.py
