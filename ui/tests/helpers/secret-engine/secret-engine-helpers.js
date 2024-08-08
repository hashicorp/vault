/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

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
      backend: backend,
      region: 'us-west-2',
      access_key: '123-key',
      iam_endpoint: 'iam-endpoint',
      sts_endpoint: 'sts-endpoint',
    },
  });
  return store.peekRecord('aws/root-config', backend);
};

const createSshCaConfig = (store, backend) => {
  store.pushPayload('ssh/ca-config', {
    id: backend,
    modelName: 'ssh/ca-config',
    data: {
      backend: backend,
      public_key: 'public-key',
      generate_signing_key: true,
    },
  });
  return store.peekRecord('ssh/ca-config', backend);
};

export function configUrl(type, backend) {
  switch (type) {
    case 'aws':
      return `${backend}/config/root`;
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
    case 'ssh':
      return createSshCaConfig(store, backend);
  }
};
// Used in tests to assert the expected keys in the config details of configurable secret engines
export const expectedConfigKeys = (type) => {
  switch (type) {
    case 'aws':
      return ['Access key', 'Region', 'IAM endpoint', 'STS endpoint'];
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
