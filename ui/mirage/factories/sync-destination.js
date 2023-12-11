/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  ['aws-sm']: trait({
    type: 'aws-sm',
    name: 'destination-aws',
    access_key_id: '*****',
    secret_access_key: '*****',
    region: 'us-west-1',
  }),
  ['azure-kv']: trait({
    type: 'azure-kv',
    name: 'destination-azure',
    key_vault_uri: 'https://keyvault-1234abcd.vault.azure.net',
    subscription_id: 'subscription-id',
    tenant_id: 'tenant-id',
    client_id: 'azure-client-id',
    client_secret: '*****',
  }),
  ['gcp-sm']: trait({
    type: 'gcp-sm',
    name: 'destination-gcp',
    credentials: '*****',
    project_id: 'gcp-project-id', // TODO backend will add, doesn't exist yet
  }),
  gh: trait({
    type: 'gh',
    name: 'destination-gh',
    access_token: '*****',
    repository_owner: 'my-organization-or-username',
    repository_name: 'my-repository',
  }),
  ['vercel-project']: trait({
    type: 'vercel-project',
    name: 'destination-vercel',
    access_token: '*****',
    project_id: 'prj_12345',
    team_id: 'team_12345',
    deployment_environments: 'development,preview', // 'production' is also an option, but left out for testing to assert form changes value
  }),
});
