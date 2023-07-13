"""
Python script that can be used 
to interact with the HashiCorp Cloud Platform Vault Secrets.
"""

import requests


class SecretVaultDriver:
    """
    Class to interact with the HashiCorp Cloud Platform Vault Secrets.
    """

    def __init__(self, hcp_client_id, hcp_client_secret):
        self.hcp_client_id = hcp_client_id
        self.hcp_client_secret = hcp_client_secret
        self.hcp_org_id = None
        self.hcp_proj_id = None
        self.app_name = None
        self.hcp_api_token = None

    def set_org_id(self, hcp_org_id: str):
        """
        Set the organization ID.
        """
        self.hcp_org_id = hcp_org_id

    def set_proj_id(self, hcp_proj_id: str):
        """
        Set the project ID.
        """
        self.hcp_proj_id = hcp_proj_id

    def set_app_name(self, app_name: str):
        """
        Set the app name.
        """
        self.app_name = app_name

    def fetch_token(self):
        """
        Fetches the API token from the HashiCorp Cloud Platform.
        This step is mendatory if you want to fetch secrets from the Vault Secrets.
        """
        url = 'https://auth.hashicorp.com/oauth/token'
        headers = {'Content-Type': 'application/json'}

        payload = {
            "audience": "https://api.hashicorp.cloud",
            "grant_type": "client_credentials",
            "client_id": self.hcp_client_id,
            "client_secret": self.hcp_client_secret
        }

        response = requests.post(url, headers=headers, json=payload, timeout=5)
        data = response.json()
        self.hcp_api_token = data.get('access_token')
        return self.hcp_api_token

    def get_all_apps(self) -> dict:
        """
        Get all the apps from the HashiCorp Cloud Platform
        """

        url = f"https://api.cloud.hashicorp.com/secrets/2023-06-13/organizations/{self.hcp_org_id}/projects/{self.hcp_proj_id}/apps"
        headers = {
            "Authorization": f"Bearer {self.hcp_api_token}"
        }

        response = requests.get(url, headers=headers, timeout=5)
        data = response.json()

        return data

    def get_all_secrets(self) -> dict:
        """
        Get all the secrets from the Vault Secrets
        """

        url = f"https://api.cloud.hashicorp.com/secrets/2023-06-13/organizations/{self.hcp_org_id}/projects/{self.hcp_proj_id}/apps/{self.app_name}/secrets"
        headers = {
            "Authorization": f"Bearer {self.hcp_api_token}"
        }

        response = requests.get(url, headers=headers, timeout=5)
        data = response.json()

        return data

    def get_secret(self, secret_name: str) -> dict:
        """
        Get a secret from the Vault Secrets
        """

        url = f"https://api.cloud.hashicorp.com/secrets/2023-06-13/organizations/{self.hcp_org_id}/projects/{self.hcp_proj_id}/apps/{self.app_name}/open/{secret_name}"
        headers = {
            "Authorization": f"Bearer {self.hcp_api_token}"
        }

        response = requests.get(url, headers=headers, timeout=5)
        data = response.json()

        return data
