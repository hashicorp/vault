/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, fillIn, visit } from '@ember/test-helpers';
import { FORM } from './kv-selectors';
import { createPolicyCmd, createTokenCmd, mountAuthCmd, runCmd } from '../commands';

// CUSTOM COMMANDS RELEVANT TO KV-V2

export const writeSecret = async function (backend, path, key, val) {
  await visit(`vault/secrets/${backend}/kv/create`);
  await fillIn(FORM.inputByAttr('path'), path);
  await fillIn(FORM.keyInput(), key);
  await fillIn(FORM.maskedValueInput(), val);
  return click(FORM.saveBtn);
};

export const writeVersionedSecret = async function (backend, path, key, val, version = 2) {
  await writeSecret(backend, path, 'key-1', 'val-1');
  for (let currentVersion = 2; currentVersion <= version; currentVersion++) {
    await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent(path)}/details/edit`);
    if (currentVersion === version) {
      await fillIn(FORM.keyInput(), key);
      await fillIn(FORM.maskedValueInput(), val);
    } else {
      await fillIn(FORM.keyInput(), `key-${currentVersion}`);
      await fillIn(FORM.maskedValueInput(), `val-${currentVersion}`);
    }
    await click(FORM.saveBtn);
  }
  return;
};

export const addSecretMetadataCmd = (backend, secret, options = { max_versions: 10 }) => {
  const stringOptions = Object.keys(options).reduce((prev, curr) => {
    return `${prev} ${curr}=${options[curr]}`;
  }, '');
  return `write ${backend}/metadata/${secret} ${stringOptions}`;
};

export const setupControlGroup = async ({
  userPolicy,
  adminUser = 'admin',
  adminPassword = 'password',
  userpassMount = 'userpass',
}) => {
  const authorizerPolicy = `
path "sys/control-group/authorize" {
  capabilities = ["update"]
}

path "sys/control-group/request" {
  capabilities = ["update"]
}
`;
  const userpassAccessor = await runCmd([
    // write policies for control group + authorization
    createPolicyCmd('kv-control-group', userPolicy),
    createPolicyCmd('authorizer', authorizerPolicy),
    // enable userpass, create admin user
    mountAuthCmd('userpass', userpassMount),
    // read out mount to get the accessor
    `read -field=accessor sys/internal/ui/mounts/auth/${userpassMount}`,
  ]);
  const authorizerEntityId = await runCmd([
    // create admin user and entity
    `write auth/${userpassMount}/users/${adminUser} password=${adminPassword} policies=default`,
    `write identity/entity name="admin-entity" policies=default`,
    `write -field=id identity/lookup/entity name="admin-entity"`,
  ]);
  const userToken = await runCmd([
    // create alias for authorizor and add them to the managers group
    `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name="admin-entity"`,
    `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
    // create a token to request access to kv/foo
    createTokenCmd('kv-control-group'),
  ]);
  return {
    userToken,
    userPolicyName: 'kv-control-group',
  };
};
