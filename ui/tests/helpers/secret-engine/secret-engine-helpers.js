/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { stringArrayToCamelCase } from 'vault/helpers/string-array-to-camel';
import { v4 as uuidv4 } from 'uuid';

export const createSecretsEngine = (store, type, path) => {
  store.pushPayload('secret-engine', {
    modelName: 'secret-engine',
    id: path,
    path: `${path}/`,
    type: type,
    data: {
      type: type,
    },
  });
  return store.peekRecord('secret-engine', path);
};
/* Create configurations methods
 * for each configuration we create the record and then push it to the store.
 */
export function configUrl(type, backend) {
  switch (type) {
    case 'aws':
      return `/${backend}/config/root`;
    case 'aws-lease':
      return `/${backend}/config/lease`;
    case 'ssh':
      return `/${backend}/config/ca`;
    default:
      return `/${backend}/config`;
  }
}

const createIssuerConfig = (store) => {
  store.pushPayload('identity/oidc/config', {
    id: 'identity-oidc-config',
    modelName: 'identity/oidc/config',
    data: {
      issuer: ``,
    },
  });
  return store.peekRecord('identity/oidc/config', 'identity-oidc-config');
};

const createAwsRootConfig = (store, backend, accessType = 'iam') => {
  // clear any records first
  store.unloadAll('aws/root-config');
  if (accessType === 'wif') {
    store.pushPayload('aws/root-config', {
      id: backend,
      modelName: 'aws/root-config',
      data: {
        backend,
        role_arn: '123-role',
        identity_token_audience: '123-audience',
        identity_token_ttl: 7200,
      },
    });
  } else if (accessType === 'no-access') {
    // set root config options that are not associated with accessType 'wif' or 'iam'
    store.pushPayload('aws/root-config', {
      id: backend,
      modelName: 'aws/root-config',
      data: {
        backend,
        region: 'ap-northeast-1',
      },
    });
  } else {
    store.pushPayload('aws/root-config', {
      id: backend,
      modelName: 'aws/root-config',
      data: {
        backend,
        region: 'us-west-2',
        access_key: '123-key',
        iam_endpoint: 'iam-endpoint',
        sts_endpoint: 'sts-endpoint',
        max_retries: 1,
      },
    });
  }
  return store.peekRecord('aws/root-config', backend);
};

const createAwsLeaseConfig = (store, backend) => {
  store.pushPayload('aws/lease-config', {
    id: backend,
    modelName: 'aws/lease-config',
    data: {
      backend,
      lease: '50s',
      lease_max: '55s',
    },
  });
  return store.peekRecord('aws/lease-config', backend);
};

const createSshCaConfig = (store, backend) => {
  store.pushPayload('ssh/ca-config', {
    id: backend,
    modelName: 'ssh/ca-config',
    data: {
      backend,
      public_key: 'public-key',
      generate_signing_key: true,
    },
  });
  return store.peekRecord('ssh/ca-config', backend);
};

const createAzureConfig = (store, backend, accessType = 'generic') => {
  // clear any records first
  // note: allowed "environment" params for testing https://github.com/hashicorp/vault-plugin-secrets-azure/blob/main/client.go#L35-L37
  store.unloadAll('azure/config');
  if (accessType === 'azure') {
    store.pushPayload('azure/config', {
      id: backend,
      modelName: 'azure/config',
      data: {
        backend,
        client_secret: 'client-secret',
        subscription_id: 'subscription-id',
        tenant_id: 'tenant-id',
        client_id: 'client-id',
        root_password_ttl: '20 days 20 hours',
        environment: 'AZUREPUBLICCLOUD',
      },
    });
  } else if (accessType === 'wif') {
    store.pushPayload('azure/config', {
      id: backend,
      modelName: 'azure/config',
      data: {
        backend,
        subscription_id: 'subscription-id',
        tenant_id: 'tenant-id',
        client_id: 'client-id',
        identity_token_audience: 'audience',
        identity_token_ttl: 7200,
        root_password_ttl: '20 days 20 hours',
        environment: 'AZUREPUBLICCLOUD',
      },
    });
  } else {
    store.pushPayload('azure/config', {
      id: backend,
      modelName: 'azure/config',
      data: {
        backend,
        subscription_id: 'subscription-id-2',
        tenant_id: 'tenant-id-2',
        client_id: 'client-id-2',
        environment: 'AZUREPUBLICCLOUD',
      },
    });
  }
  return store.peekRecord('azure/config', backend);
};

const createGcpConfig = (store, backend, accessType = 'gcp') => {
  // clear any records first
  store.unloadAll('gcp/config');
  if (accessType === 'wif') {
    store.pushPayload('gcp/config', {
      id: backend,
      modelName: 'gcp/config',
      data: {
        backend,
        service_account_email: 'service-email',
        identity_token_audience: 'audience',
        identity_token_ttl: 7200,
      },
    });
  } else {
    store.pushPayload('gcp/config', {
      id: backend,
      modelName: 'gcp/config',
      data: {
        backend,
        credentials: '{"some-key":"some-value"}',
        ttl: '1 hour',
        max_ttl: '4 hours',
      },
    });
  }
  return store.peekRecord('gcp/config', backend);
};

export const createConfig = (store, backend, type) => {
  switch (type) {
    case 'aws':
    case 'aws-generic':
      return createAwsRootConfig(store, backend);
    case 'aws-wif':
      return createAwsRootConfig(store, backend, 'wif');
    case 'aws-no-access':
      return createAwsRootConfig(store, backend, 'no-access');
    case 'issuer':
      return createIssuerConfig(store);
    case 'aws-lease':
      return createAwsLeaseConfig(store, backend);
    case 'ssh':
      return createSshCaConfig(store, backend);
    case 'azure':
      return createAzureConfig(store, backend, 'azure');
    case 'azure-wif':
      return createAzureConfig(store, backend, 'wif');
    case 'azure-generic':
      return createAzureConfig(store, backend, 'generic');
    case 'gcp':
      return createGcpConfig(store, backend);
  }
};
/* Manually create the configuration by filling in the configuration form */
export const fillInAwsConfig = async (situation = 'withAccess') => {
  if (situation === 'withAccess') {
    await fillIn(GENERAL.inputByAttr('accessKey'), 'foo');
    await fillIn(GENERAL.inputByAttr('secretKey'), 'bar');
  }
  if (situation === 'withAccessOptions') {
    await click(GENERAL.toggleGroup('Root config options'));
    await fillIn(GENERAL.inputByAttr('region'), 'ca-central-1');
    await fillIn(GENERAL.inputByAttr('iamEndpoint'), 'iam-endpoint');
    await fillIn(GENERAL.inputByAttr('stsEndpoint'), 'sts-endpoint');
    await fillIn(GENERAL.inputByAttr('maxRetries'), '3');
  }
  if (situation === 'withLease') {
    await click(GENERAL.ttl.toggle('Default Lease TTL'));
    await fillIn(GENERAL.ttl.input('Default Lease TTL'), '33');
    await click(GENERAL.ttl.toggle('Max Lease TTL'));
    await fillIn(GENERAL.ttl.input('Max Lease TTL'), '44');
  }
  if (situation === 'withWif') {
    await click(SES.wif.accessType('wif')); // toggle to wif
    await fillIn(GENERAL.inputByAttr('issuer'), `http://bar.${uuidv4()}`); // make random because global setting
    await fillIn(GENERAL.inputByAttr('roleArn'), 'foo-role');
    await fillIn(GENERAL.inputByAttr('identityTokenAudience'), 'foo-audience');
    await click(GENERAL.ttl.toggle('Identity token TTL'));
    await fillIn(GENERAL.ttl.input('Identity token TTL'), '7200');
  }
};

export const fillInAzureConfig = async (situation = 'azure') => {
  await fillIn(GENERAL.inputByAttr('subscriptionId'), 'subscription-id');
  await fillIn(GENERAL.inputByAttr('tenantId'), 'tenant-id');
  await fillIn(GENERAL.inputByAttr('clientId'), 'client-id');
  await fillIn(GENERAL.inputByAttr('environment'), 'AZUREPUBLICCLOUD');

  if (situation === 'azure') {
    await fillIn(GENERAL.inputByAttr('clientSecret'), 'client-secret');
    await click(GENERAL.ttl.toggle('Root password TTL'));
    await fillIn(GENERAL.ttl.input('Root password TTL'), '5200');
  }
  if (situation === 'withWif') {
    await click(SES.wif.accessType('wif')); // toggle to wif
    await fillIn(GENERAL.inputByAttr('identityTokenAudience'), 'azure-audience');
    await click(GENERAL.ttl.toggle('Identity token TTL'));
    await fillIn(GENERAL.ttl.input('Identity token TTL'), '7200');
  }
};

/* Generate arrays of keys to iterate over.
 * used to check details of the secret engine configuration
 * and used to check the form to configure the secret engine
 */
// WIF specific keys
const genericWifKeys = ['Identity token audience', 'Identity token TTL'];
// AWS specific keys
const awsLeaseKeys = ['Default Lease TTL', 'Max Lease TTL'];
const awsKeys = ['Access key', 'Secret key', 'Region', 'IAM endpoint', 'STS endpoint', 'Max retries'];
const awsWifKeys = ['Issuer', 'Role ARN', ...genericWifKeys];
// Azure specific keys
const genericAzureKeys = ['Subscription ID', 'Tenant ID', 'Client ID', 'Environment'];
const azureKeys = [...genericAzureKeys, 'Client secret', 'Root password TTL'];
const azureWifKeys = [...genericAzureKeys, ...genericWifKeys];
// GCP specific keys
const genericGcpKeys = ['Config TTL', 'Max TTL'];
const gcpKeys = [...genericGcpKeys, 'Credentials'];
const gcpWifKeys = [...genericGcpKeys, ...genericWifKeys, 'Service account email'];
// SSH specific keys
const sshKeys = ['Private key', 'Public key', 'Generate signing key'];

export const expectedConfigKeys = (type, camelCase = false) => {
  switch (type) {
    case 'aws':
      return camelCase ? stringArrayToCamelCase(awsKeys) : awsKeys;
    case 'aws-wif':
      return camelCase ? stringArrayToCamelCase(awsWifKeys) : awsWifKeys;
    case 'aws-lease':
      return camelCase ? stringArrayToCamelCase(awsLeaseKeys) : awsLeaseKeys;
    case 'azure':
      return camelCase ? stringArrayToCamelCase(azureKeys) : azureKeys;
    case 'azure-wif':
      return camelCase ? stringArrayToCamelCase(azureWifKeys) : azureWifKeys;
    case 'gcp':
      return camelCase ? stringArrayToCamelCase(gcpKeys) : gcpKeys;
    case 'gcp-wif':
      return camelCase ? stringArrayToCamelCase(gcpWifKeys) : gcpWifKeys;
    case 'ssh':
      return camelCase ? stringArrayToCamelCase(sshKeys) : sshKeys;
  }
};

const valueOfAwsKeys = (string) => {
  switch (string) {
    case 'Access key':
      return '123-key';
    case 'Region':
      return 'us-west-2';
    case 'IAM endpoint':
      return 'iam-endpoint';
    case 'STS endpoint':
      return 'sts-endpoint';
    case 'Max retries':
      return '1';
  }
};

const valueOfAzureKeys = (string) => {
  switch (string) {
    case 'Subscription ID':
      return 'subscription-id';
    case 'Tenant ID':
      return 'tenant-id';
    case 'Client ID':
      return 'client-id';
    case 'Environment':
      return 'AZUREPUBLICCLOUD';
    case 'Root password TTL':
      return '20 days 20 hours';
    case 'Identity token audience':
      return 'audience';
    case 'Identity token TTL':
      return '8 days 8 hours';
  }
};

const valueOfGcpKeys = (string) => {
  switch (string) {
    case 'Credentials':
      return '"{"some-key":"some-value"}",';
    case 'Service account email':
      return 'service-email';
    case 'Config TTL':
      return '1 hour';
    case 'Max TTL':
      return '4 hours';
    case 'Identity token audience':
      return 'audience';
    case 'Identity token TTL':
      return '8 days 8 hours';
  }
};

const valueOfSshKeys = (string) => {
  switch (string) {
    case 'Public key':
      return '***********';
    case 'Generate signing key':
      return 'Yes';
  }
};
// Used in tests to assert the expected values in the config details of configurable secret engines
export const expectedValueOfConfigKeys = (type, string) => {
  switch (type) {
    case 'aws':
      return valueOfAwsKeys(string);
    case 'azure':
      return valueOfAzureKeys(string);
    case 'gcp':
      return valueOfGcpKeys(string);
    case 'ssh':
      return valueOfSshKeys(string);
  }
};

// Example usage
// createLongJson (2, 3) will create a json object with 2 original keys, each with 3 nested keys
// {
// 	"key-0": {
// 		"nested-key-0": {
// 			"nested-key-1": {
// 				"nested-key-2": "nested-value"
// 			}
// 		}
// 	},
// 	"key-1": {
// 		"nested-key-0": {
// 			"nested-key-1": {
// 				"nested-key-2": "nested-value"
// 			}
// 		}
// 	}
// }

export function createLongJson(lines = 10, nestLevel = 3) {
  const keys = Array.from({ length: nestLevel }, (_, i) => `nested-key-${i}`);
  const jsonObject = {};

  for (let i = 0; i < lines; i++) {
    nestLevel > 0
      ? (jsonObject[`key-${i}`] = createNestedObject({}, keys, 'nested-value'))
      : (jsonObject[`key-${i}`] = 'non-nested-value');
  }
  return jsonObject;
}

function createNestedObject(obj = {}, keys, value) {
  let current = obj;

  for (let i = 0; i < keys.length - 1; i++) {
    const key = keys[i];
    if (!current[key]) {
      current[key] = {};
    }
    current = current[key];
  }

  current[keys[keys.length - 1]] = value;
  return obj;
}
