/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

const SUPPORTED_DYNAMIC_BACKENDS = ['database', 'ssh', 'aws', 'totp'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',
  pathHelp: service('path-help'),
  router: service(),
  api: service(),

  beforeModel(transition) {
    const { id: backendPath, type: backendType } = this.modelFor('vault.cluster.secrets.backend');
    // redirect if the backend type does not support credentials
    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendType)) {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root', backendPath);
    }

    // assign back button route
    if (backendType === 'totp') {
      const previousRoute = transition.from?.name ?? 'vault.cluster.secrets.backend.list-root';
      this.set('backRoute', previousRoute);
    }
  },

  // Fetches credentials with their resolved type without swallowing request failures.
  async getDatabaseCredentialByType(backend, secret, roleType) {
    if (roleType === 'static') {
      const { last_vault_rotation, lease_duration, data } =
        await this.api.secrets.databaseReadStaticRoleCredentials(secret, backend);
      const { rsa_private_key, ...rest } = data;
      return {
        last_vault_rotation,
        lease_duration,
        ...rest,
        ...(rsa_private_key !== undefined ? { rsaPrivateKey: rsa_private_key } : {}),
        roleType: 'static',
      };
    } else {
      const { data, lease_id, lease_duration } = await this.api.secrets.databaseGenerateCredentials(
        secret,
        backend
      );
      const { rsa_private_key: rsaPrivateKey, ...rest } = data;
      return {
        ...rest,
        ...(rsaPrivateKey !== undefined ? { rsaPrivateKey } : {}),
        lease_id,
        lease_duration,
        roleType: 'dynamic',
      };
    }
  },

  async getDatabaseCredential(backend, secret, roleType = '') {
    try {
      if (roleType === 'static' || roleType === 'dynamic') {
        return await this.getDatabaseCredentialByType(backend, secret, roleType);
      }

      // Overview permits entering a role name without knowing its type or having permission to read it.
      const [staticResp, dynamicResp] = await Promise.allSettled([
        this.getDatabaseCredentialByType(backend, secret, 'static'),
        this.getDatabaseCredentialByType(backend, secret, 'dynamic'),
      ]);
      if (staticResp.status === 'fulfilled') return staticResp.value;
      if (dynamicResp.status === 'fulfilled') return dynamicResp.value;

      const { status: staticStatus, response: staticError } = await this.api.parseError(staticResp.reason);
      const { status: dynamicStatus, response: dynamicError } = await this.api.parseError(dynamicResp.reason);
      if (dynamicError?.isControlGroupError) throw dynamicResp.reason;
      if (staticError?.isControlGroupError) throw staticResp.reason;
      throw staticStatus < dynamicStatus ? dynamicResp.reason : staticResp.reason;
    } catch (error) {
      const { response, status, message, path } = await this.api.parseError(error);
      if (response?.isControlGroupError) {
        throw response;
      }
      // Unless it's a control group error, we want to pass back error info
      // so we can render it on the GenerateCredentialsDatabase component
      return { ...response, status, message, path };
    }
  },

  async getAwsRole(backend, id) {
    try {
      return await this.api.secrets.awsReadRole(id, backend);
    } catch (e) {
      // swallow error, non-essential data
      return;
    }
  },

  async getTotpKey(backend, keyName) {
    try {
      const { data } = await this.api.secrets.totpReadKey(keyName, backend);
      return data || {};
    } catch (e) {
      // swallow error, non-essential data
      return {};
    }
  },

  async model(params) {
    const role = params.secret;
    const { id: backendPath, type: backendType } = this.modelFor('vault.cluster.secrets.backend');
    const backendData = { backendPath, backendType };
    const roleType = params.roleType;
    let dbCred, awsRole, totpCodePeriod, backRoute;
    if (backendType === 'database') {
      dbCred = await this.getDatabaseCredential(backendPath, role, roleType);
    } else if (backendType === 'aws') {
      awsRole = await this.getAwsRole(backendPath, role);
    } else if (backendType === 'totp') {
      totpCodePeriod = (await this.getTotpKey(backendPath, role))?.period ?? 30;
      backRoute = this.backRoute;
      return { ...backendData, keyName: role, totpCodePeriod, backRoute };
    }

    return {
      ...backendData,
      roleName: role,
      roleType,
      dbCred,
      awsRoleType: awsRole?.credential_type,
    };
  },

  resetController(controller) {
    controller.reset();
  },
});
