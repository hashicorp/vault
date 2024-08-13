/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn } from '@ember/test-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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

const createAwsRootConfig = (store, backend) => {
  store.pushPayload('aws/root-config', {
    id: backend,
    modelName: 'aws/root-config',
    data: {
      backend,
      region: 'us-west-2',
      access_key: '123-key',
      iam_endpoint: 'iam-endpoint',
      sts_endpoint: 'sts-endpoint',
    },
  });
  return store.peekRecord('aws/root-config', backend);
};

const createAwsLeaseConfig = (store, backend) => {
  store.pushPayload('aws/lease-config', {
    id: backend,
    modelName: 'aws/lease-config',
    data: {
      backend: backend,
      lease: '50s',
      leaseMax: '55s',
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
      generate_signing_key: true,
    },
  });
  return store.peekRecord('ssh/ca-config', backend);
};

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
// send the type of config you want and the name of the backend path to push the config to the store.
export const createConfig = (store, backend, type) => {
  switch (type) {
    case 'aws':
      return createAwsRootConfig(store, backend);
    case 'aws-lease':
      return createAwsLeaseConfig(store, backend);
    case 'ssh':
      return createSshCaConfig(store, backend);
  }
};
// Used in tests to assert the expected keys in the config details of configurable secret engines
export const expectedConfigKeys = (type) => {
  switch (type) {
    case 'aws':
      return ['Access key', 'Region', 'IAM endpoint', 'STS endpoint', 'Maximum retries'];
    case 'ssh':
      return ['Public key', 'Generate signing key'];
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
    case 'Maximum retries':
      return '-1';
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
    case 'ssh':
      return valueOfSshKeys(string);
  }
};

export const fillInAwsConfig = async (withAccess = true, withAccessOptions = false, withLease = false) => {
  if (withAccess) {
    await fillIn(GENERAL.inputByAttr('accessKey'), 'foo');
    await fillIn(GENERAL.inputByAttr('secretKey'), 'bar');
  }
  if (withAccessOptions) {
    await click(GENERAL.toggleGroup('Root config options'));
    await fillIn(GENERAL.selectByAttr('region'), 'ca-central-1');
    await fillIn(GENERAL.inputByAttr('iamEndpoint'), 'iam-endpoint');
    await fillIn(GENERAL.inputByAttr('stsEndpoint'), 'sts-endpoint');
    await fillIn(GENERAL.inputByAttr('maxRetries'), '3');
  }
  await click(GENERAL.hdsTab('lease'));
  if (withLease) {
    await click(GENERAL.ttl.toggle('Lease'));
    await fillIn(GENERAL.ttl.input('Lease'), '33');
    await click(GENERAL.ttl.toggle('Maximum Lease'));
    await fillIn(GENERAL.ttl.input('Maximum Lease'), '44');
  }
};
