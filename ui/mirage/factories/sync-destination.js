/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  ['aws-sm']: trait({
    access_key_id: '*****',
    secret_access_key: '*****',
    region: 'us-west-1',
    type: 'aws-sm',
    name: 'destination-aws',
  }),
  ['azure-kv']: trait({
    key_vault_uri: 'https://keyvault-1234abcd.vault.azure.net',
    subscription_id: 'subscription-id',
    tenant_id: 'tenant-id',
    client_id: 'client-id',
    client_secret: 'my-secret',
    type: 'azure-kv',
    name: 'destination-azure',
  }),
  ['gcp-sm']: trait({
    credentials: '{"username":"foo","password":"bar"}',
    type: 'gcp-sm',
    name: 'destination-gcp',
  }),
  gh: trait({
    access_token: 'github_pat_12345',
    repository_owner: 'my-organization-or-username',
    repository_name: 'my-repository',
    type: 'gh',
    name: 'destination-gh',
  }),
  ['vercel-project']: trait({
    access_token: 'my-access-token',
    project_id: 'prj_12345',
    deployment_environments: ['development', 'preview', 'production'],
    type: 'vercel-project',
    name: 'destination-vercel',
  }),
});
