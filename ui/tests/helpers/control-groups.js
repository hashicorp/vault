/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, visit } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';

import { CONTROL_GROUP_PREFIX, TOKEN_SEPARATOR } from 'vault/services/control-group';
import authPage from 'vault/tests/pages/auth';
import controlGroup from 'vault/tests/pages/components/control-group';
import { createPolicyCmd, createTokenCmd, mountAuthCmd, runCmd } from './commands';
const controlGroupComponent = create(controlGroup);

const storageKey = (accessor, path) => {
  return `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
};

// This function is used to setup a control group for testing
// It will create a userpass backend, create an authorizing user,
// and create a controlled access token. The auth mount and policy
// names will be appended with the backend
export const setupControlGroup = async ({
  userPolicy,
  backend,
  adminUser = 'authorizer',
  adminPassword = 'testing-xyz',
}) => {
  if (!backend || !userPolicy) {
    throw new Error('missing required fields for setupControlGroup');
  }
  const userpassMount = `userpass-${backend}`;
  const userPolicyName = `kv-control-group-${backend}`;
  const authorizerPolicyName = `authorizer-${backend}`;
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
    createPolicyCmd(userPolicyName, userPolicy),
    createPolicyCmd(authorizerPolicyName, authorizerPolicy),
    // enable userpass, create admin user
    mountAuthCmd('userpass', userpassMount),
    // read out mount to get the accessor
    `read -field=accessor sys/internal/ui/mounts/auth/${userpassMount}`,
  ]);
  const authorizerEntityId = await runCmd([
    // create admin user and entity
    `write auth/${userpassMount}/users/${adminUser} password=${adminPassword} policies=default`,
    `write identity/entity name=${adminUser} policies=test`,
    `write -field=id identity/lookup/entity name=${adminUser}`,
  ]);
  const userToken = await runCmd([
    // create alias for authorizor and add them to the managers group
    `write identity/alias mount_accessor=${userpassAccessor} entity_id=${authorizerEntityId} name=${adminUser}`,
    `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=${authorizerPolicyName}`,
    // create a token to request access to kv/foo
    createTokenCmd(userPolicyName),
  ]);
  return {
    userToken,
    userPolicyName,
    adminUser,
    adminPassword,
    userpassMount,
  };
};

export async function grantAccessForWrite({
  token,
  accessor,
  creation_path,
  originUrl,
  userToken,
  backend,
  authorizerUser = 'authorizer',
  authorizerPassword = 'testing-xyz',
}) {
  if (!token || !accessor || !creation_path || !originUrl || !userToken || !backend) {
    throw new Error('missing required fields for grantAccessForWrite');
  }
  const userpassMount = `userpass-${backend}`;
  await authPage.loginUsername(authorizerUser, authorizerPassword, userpassMount);
  await visit(`/vault/access/control-groups/${accessor}`);
  await controlGroupComponent.authorize();
  await authPage.login(userToken);
  localStorage.setItem(
    storageKey(accessor, creation_path),
    JSON.stringify({
      accessor,
      token,
      creation_path,
      uiParams: {
        url: originUrl,
      },
    })
  );
  await visit(originUrl);
}

/*
 * Control group grant access flow
 * Assumes start on route 'vault.cluster.access.control-group-accessor'
 * and authorizer login is via userpass
 */
export async function grantAccess({
  apiPath,
  originUrl,
  userToken,
  backend,
  authorizerUser = 'authorizer',
  authorizerPassword = 'testing-xyz',
}) {
  if (!apiPath || !originUrl || !userToken || !backend) {
    throw new Error('missing required fields for grantAccess');
  }
  const userpassMount = `userpass-${backend}`;
  const accessor = controlGroupComponent.accessor;
  const controlGroupToken = controlGroupComponent.token;
  await authPage.loginUsername(authorizerUser, authorizerPassword, userpassMount);
  await visit(`/vault/access/control-groups/${accessor}`);
  await controlGroupComponent.authorize();
  await authPage.login(userToken);
  localStorage.setItem(
    storageKey(accessor, apiPath),
    JSON.stringify({
      accessor,
      token: controlGroupToken,
      creation_path: apiPath,
      uiParams: {
        url: originUrl,
      },
    })
  );
  await visit(`/vault/access/control-groups/${accessor}`);
  await click(`[data-test-navigate-button]`);
  /* end of control group authorization flow */
}
