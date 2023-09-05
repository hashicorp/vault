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

export const setupControlGroup = async ({
  userPolicy,
  adminUser = 'authorizer',
  adminPassword = 'password',
  userpassMount = 'userpass',
}) => {
  const userPolicyName = 'kv-control-group';
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
    createPolicyCmd('authorizer', authorizerPolicy),
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
    `write identity/group name=managers member_entity_ids=${authorizerEntityId} policies=authorizer`,
    // create a token to request access to kv/foo
    createTokenCmd(userPolicyName),
  ]);
  return {
    userToken,
    userPolicyName,
    userPolicy,
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
  authorizerUser = 'authorizer',
  authorizerPassword = 'password',
}) {
  await authPage.loginUsername(authorizerUser, authorizerPassword);
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

export async function grantAccess({
  apiPath,
  originUrl,
  userToken,
  authorizerUser = 'authorizer',
  authorizerPassword = 'password',
}) {
  /*
   * Control group grant access flow
   * Assumes start on route 'vault.cluster.access.control-group-accessor'
   * and authorizer login is via userpass
   */
  const accessor = controlGroupComponent.accessor;
  const controlGroupToken = controlGroupComponent.token;
  await authPage.loginUsername(authorizerUser, authorizerPassword);
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
