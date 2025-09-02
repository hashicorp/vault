/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, find } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { v4 as uuidv4 } from 'uuid';
import SecretsEngineResource from 'vault/resources/secrets/engine';

/* Secret Create/Edit methods */
// ARG TODO unsure if should be moved to another file
export async function createSecret(path, key, value) {
  await fillIn(SES.secretPath('create'), path);
  await fillIn('[data-test-secret-key]', key);
  await fillIn('[data-test-secret-value] textarea', value);
  await click(GENERAL.submitButton);
  return;
}

export const createSecretsEngine = (store, type, path) => {
  if (store) {
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
  }

  return new SecretsEngineResource({
    path: `${path}/`,
    type,
  });
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

const createAwsRootConfig = (accessType = 'iam') => {
  if (accessType === 'wif') {
    return {
      role_arn: '123-role',
      identity_token_audience: '123-audience',
      identity_token_ttl: 7200,
    };
  } else if (accessType === 'no-access') {
    // set root config options that are not associated with accessType 'wif' or 'iam'
    return {
      region: 'ap-northeast-1',
    };
  } else {
    return {
      region: 'us-west-2',
      access_key: '123-key',
      iam_endpoint: 'iam-endpoint',
      sts_endpoint: 'sts-endpoint',
      max_retries: 1,
    };
  }
};

const createAwsLeaseConfig = () => {
  return {
    lease: '50s',
    lease_max: '55s',
  };
};

const createSshCaConfig = () => {
  return {
    public_key: 'public-key',
    generate_signing_key: true,
  };
};

const createAzureConfig = (accessType = 'generic') => {
  // note: allowed "environment" params for testing https://github.com/hashicorp/vault-plugin-secrets-azure/blob/main/client.go#L35-L37
  if (accessType === 'azure') {
    return {
      client_secret: 'client-secret',
      subscription_id: 'subscription-id',
      tenant_id: 'tenant-id',
      client_id: 'client-id',
      root_password_ttl: '1800000s',
      environment: 'AZUREPUBLICCLOUD',
    };
  } else if (accessType === 'wif') {
    return {
      subscription_id: 'subscription-id',
      tenant_id: 'tenant-id',
      client_id: 'client-id',
      identity_token_audience: 'audience',
      identity_token_ttl: 7200,
      root_password_ttl: '1800000s',
      environment: 'AZUREPUBLICCLOUD',
    };
  } else {
    return {
      subscription_id: 'subscription-id-2',
      tenant_id: 'tenant-id-2',
      client_id: 'client-id-2',
      environment: 'AZUREPUBLICCLOUD',
      root_password_ttl: '1800000s',
    };
  }
};

const createGcpConfig = (accessType = 'gcp') => {
  if (accessType === 'wif') {
    return {
      service_account_email: 'service-email',
      identity_token_audience: 'audience',
      identity_token_ttl: 7200,
    };
  } else {
    return {
      credentials: '{"some-key":"some-value"}',
      ttl: '100s',
      max_ttl: '101s',
    };
  }
};

export const createConfig = (type) => {
  switch (type) {
    case 'aws':
    case 'aws-generic':
      return createAwsRootConfig();
    case 'aws-wif':
      return createAwsRootConfig('wif');
    case 'aws-no-access':
      return createAwsRootConfig('no-access');
    case 'aws-lease':
      return createAwsLeaseConfig();
    case 'ssh':
      return createSshCaConfig();
    case 'azure':
      return createAzureConfig('azure');
    case 'azure-wif':
      return createAzureConfig('wif');
    case 'azure-generic':
      return createAzureConfig('generic');
    case 'gcp':
    case 'gcp-generic':
      return createGcpConfig();
    case 'gcp-wif':
      return createGcpConfig('wif');
  }
};
/* Manually create the configuration by filling in the configuration form */
export const fillInAwsConfig = async (situation = 'withAccess') => {
  if (situation === 'withAccess') {
    await fillIn(GENERAL.inputByAttr('access_key'), 'foo');
    await fillIn(GENERAL.inputByAttr('secret_key'), 'bar');
  }
  if (situation === 'withAccessOptions') {
    await click(GENERAL.button('Root config options'));
    await fillIn(GENERAL.inputByAttr('region'), 'ca-central-1');
    await fillIn(GENERAL.inputByAttr('iam_endpoint'), 'iam-endpoint');
    await fillIn(GENERAL.inputByAttr('sts_endpoint'), 'sts-endpoint');
    await fillIn(GENERAL.inputByAttr('max_retries'), '3');
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
    await fillIn(GENERAL.inputByAttr('role_arn'), 'foo-role');
    await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'foo-audience');
    await click(GENERAL.ttl.toggle('Identity token TTL'));
    await fillIn(GENERAL.ttl.input('Identity token TTL'), '7200');
  }
};

export const fillInAzureConfig = async (withWif = false) => {
  await fillIn(GENERAL.inputByAttr('subscription_id'), 'subscription-id');
  await fillIn(GENERAL.inputByAttr('tenant_id'), 'tenant-id');
  await fillIn(GENERAL.inputByAttr('client_id'), 'client-id');
  // options may already be toggled so check before clicking
  if (!find(GENERAL.inputByAttr('environment'))) {
    await click(GENERAL.button('More options'));
  }
  await fillIn(GENERAL.inputByAttr('environment'), 'AZUREPUBLICCLOUD');
  // similarly, the root password TTL may already be toggled
  if (!find(GENERAL.ttl.input('Root password TTL'))) {
    await click(GENERAL.ttl.toggle('Root password TTL'));
  }
  await fillIn(GENERAL.ttl.input('Root password TTL'), '200');

  if (withWif) {
    await click(SES.wif.accessType('wif')); // toggle to wif
    await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'azure-audience');
    await click(GENERAL.ttl.toggle('Identity token TTL'));
    await fillIn(GENERAL.ttl.input('Identity token TTL'), '7200');
  } else {
    await fillIn(GENERAL.inputByAttr('client_secret'), 'client-secret');
  }
};

export const fillInGcpConfig = async (withWif = false) => {
  if (withWif) {
    await click(SES.wif.accessType('wif')); // toggle to wif
    await fillIn(GENERAL.inputByAttr('identity_token_audience'), 'azure-audience');
    await click(GENERAL.ttl.toggle('Identity token TTL'));
    await fillIn(GENERAL.ttl.input('Identity token TTL'), '7200');
    await fillIn(GENERAL.inputByAttr('service_account_email'), 'some@email.com');
  } else {
    await click(GENERAL.button('More options'));
    await click(GENERAL.ttl.toggle('Config TTL'));
    await fillIn(GENERAL.ttl.input('Config TTL'), '7200');
    await click(GENERAL.ttl.toggle('Max TTL'));
    await fillIn(GENERAL.ttl.input('Max TTL'), '8200');
    await click(GENERAL.textToggle);
    await fillIn(GENERAL.maskedInput, '{"some-key":"some-value"}');
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
const genericAzureKeys = ['Subscription ID', 'Tenant ID', 'Client ID', 'Environment', 'Root password TTL'];
const azureKeys = [...genericAzureKeys, 'Client secret'];
const azureWifKeys = [...genericAzureKeys, ...genericWifKeys];
// GCP specific keys
const genericGcpKeys = ['Config TTL', 'Max TTL'];
const gcpKeys = [...genericGcpKeys, 'Credentials'];
const gcpWifKeys = [...genericWifKeys, 'Service account email'];
// SSH specific keys
const sshKeys = ['Private key', 'Public key', 'Generate signing key'];

export const expectedConfigKeys = (type, snake_case = false) => {
  const getKeys = (keys) => (snake_case ? keys.map((str) => str.replace(/\s+/g, '_').toLowerCase()) : keys);

  switch (type) {
    case 'aws':
      return getKeys(awsKeys);
    case 'aws-wif':
      return getKeys(awsWifKeys);
    case 'aws-lease':
      return getKeys(awsLeaseKeys);
    case 'azure':
      return getKeys(azureKeys);
    case 'azure-wif':
      return getKeys(azureWifKeys);
    case 'gcp':
      return getKeys(gcpKeys);
    case 'gcp-wif':
      return getKeys(gcpWifKeys);
    case 'ssh':
      return getKeys(sshKeys);
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
      return '1 minute 40 seconds';
    case 'Max TTL':
      return '1 minute 41 seconds';
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
