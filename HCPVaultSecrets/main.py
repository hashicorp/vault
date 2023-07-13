"""
Main script to test the Secret Vault Driver
"""
import os
from dotenv import load_dotenv
from secret_vault_driver_script import SecretVaultDriver

# LOAD ENV VARIABLES

load_dotenv()

# Mandatory variables
HCP_CLIENT_ID = os.getenv("HCP_CLIENT_ID")
HCP_CLIENT_SECRET = os.getenv("HCP_CLIENT_SECRET")
HCP_ORG_ID = os.getenv("HCP_ORG_ID")
HCP_PROJ_ID = os.getenv("HCP_PROJ_ID")

# Optional variables
APP_NAME = os.getenv("APP_NAME")


# Create an instance of SecretVaultDriver
driver = SecretVaultDriver(
    HCP_CLIENT_ID, HCP_CLIENT_SECRET)

# Fetch the API token, this step is mandatory
driver.fetch_token()

# Get all the apps of a specific project

# You need to set the org_id and proj_id before calling the method get_all_apps()
driver.set_org_id(HCP_ORG_ID)
driver.set_proj_id(HCP_PROJ_ID)

print(driver.get_all_apps())

# Get all the secrets of the app

# we can get APP_NAME from get_all_apps() method too.
driver.set_app_name(APP_NAME)

print(driver.get_all_secrets())

# Get a secret
secret = driver.get_secret("username")
print(secret['secret']['version']['value'])
