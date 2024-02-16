/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  ['aws-sm']: trait({
    type: 'aws-sm',
    name: 'destination-aws',
    // connection_details
    access_key_id: '*****',
    secret_access_key: '*****',
    region: 'us-west-1',
    // options
    granularity: 'secret-path',
    secret_name_template: 'vault-{{ .MountAccessor | replace "_" "-" }}-{{ .SecretPath }}',
    custom_tags: { foo: 'bar' },
  }),
  ['azure-kv']: trait({
    type: 'azure-kv',
    name: 'destination-azure',
    // connection_details
    key_vault_uri: 'https://keyvault-1234abcd.vault.azure.net',
    subscription_id: 'subscription-id',
    tenant_id: 'tenant-id',
    client_id: 'azure-client-id',
    client_secret: '*****',
    cloud: 'Azure Public Cloud',
    // options
    granularity: 'secret-path',
    secret_name_template: 'vault-{{ .MountAccessor | replace "_" "-" }}-{{ .SecretPath }}',
    custom_tags: { foo: 'bar' },
  }),
  ['gcp-sm']: trait({
    type: 'gcp-sm',
    name: 'destination-gcp',
    // connection_details
    credentials: '*****',
    // options
    granularity: 'secret-path',
    secret_name_template: 'vault-{{ .MountAccessor | replace "_" "-" }}-{{ .SecretPath }}',
    custom_tags: { foo: 'bar' },
  }),
  gh: trait({
    type: 'gh',
    name: 'destination-gh',
    // connection_details
    access_token: '*****',
    repository_owner: 'my-organization-or-username',
    repository_name: 'my-repository',
    // options
    granularity: 'secret-path',
    secret_name_template: 'vault-{{ .MountAccessor | replace "_" "-" }}-{{ .SecretPath }}',
  }),
  ['vercel-project']: trait({
    type: 'vercel-project',
    name: 'destination-vercel',
    // connection_details
    access_token: '*****',
    project_id: 'prj_12345',
    team_id: 'team_12345',
    deployment_environments: ['development', 'preview'], // 'production' is also an option, but left out for testing to assert form changes value
    // options
    secret_name_template: 'vault-{{ .MountAccessor | replace "_" "-" }}-{{ .SecretPath }}',
  }),
});
