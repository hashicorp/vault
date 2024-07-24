/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// import { visit, currentURL } from '@ember/test-helpers';

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

export const createAwsRootConfig = (store, backend) => {
  store.pushPayload('aws/root-config', {
    modelName: 'aws/root-config',
    data: {
      backend: backend,
      region: 'us-west-2',
      accessKey: '123-key',
      iamEndpoint: 'iam-endpoint',
      stsEndpoint: 'sts-endpoint',
    },
  });
  return store.peekRecord('aws/root-config', backend);
};

export const createSshCaConfig = (store, backend) => {
  store.pushPayload('aws/root-config', {
    modelName: 'aws/root-config',
    data: {
      backend: backend,
      publicKey: 'public-key',
      generateSigningKey: true,
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
