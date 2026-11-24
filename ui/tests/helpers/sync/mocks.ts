/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SYNC_DESTINATION_AWS_WIF_RESPONSE = {
  request_id: 'c2f069c7-bee5-ce3c-992c-59422c55f45a',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    connection_details: {
      identity_token_audience: '*****',
      identity_token_ttl: 3600,
      region: 'us-east-1',
      role_arn: 'arn:aws:iam::111111111111:role/wif_test',
    },
    name: 'test-aws',
    options: {
      custom_tags: {},
      granularity_level: 'secret-path',
      secret_name_template: 'vault/{{ .MountAccessor }}/{{ .SecretPath }}',
    },
    type: 'aws-sm',
    uses_wif: true,
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: 'system',
};

export const SYNC_DESTINATION_AZURE_WIF_RESPONSE = {
  request_id: 'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    connection_details: {
      identity_token_audience: '*****',
      identity_token_ttl: 3600,
      key_vault_uri: 'https://test-keyvault.vault.azure.net',
      tenant_id: '11111111-1111-1111-1111-111111111111',
      client_id: 'test-client-id',
    },
    name: 'test-azure',
    options: {
      custom_tags: {},
      granularity_level: 'secret-path',
      secret_name_template: 'vault/{{ .MountAccessor }}/{{ .SecretPath }}',
    },
    type: 'azure-kv',
    uses_wif: true,
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: 'system',
};

export const SYNC_DESTINATION_GCP_WIF_RESPONSE = {
  request_id: 'b2c3d4e5-f6a7-8901-bcde-f12345678901',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    connection_details: {
      identity_token_audience: '*****',
      identity_token_ttl: 3600,
      project_id: 'test-gcp-project',
      service_account_email: 'test-sa@test-gcp-project.iam.gserviceaccount.com',
    },
    name: 'test-gcp',
    options: {
      custom_tags: {},
      granularity_level: 'secret-path',
      secret_name_template: 'vault/{{ .MountAccessor }}/{{ .SecretPath }}',
    },
    type: 'gcp-sm',
    uses_wif: true,
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: 'system',
};
