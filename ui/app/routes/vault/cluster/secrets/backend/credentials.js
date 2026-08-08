/**
 * Copyright IBM Corp. 2016, 2026
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

  async getDatabaseCredential(backend, secret, roleType = '') {
    try {
      if (roleType === 'static') {
        const { last_vault_rotation, lease_duration, data } =
          await this.api.secrets.databaseReadStaticRoleCredentials(secret, backend);
        return {
          last_vault_rotation,
          lease_duration,
          ...data,
        };
      } else {
        const { data, lease_id, lease_duration } = await this.api.secrets.databaseGenerateCredentials(
          secret,
          backend
        );
        return {
          ...data,
          lease_id,
          lease_duration,
        };
      }
    } catch (error) {
      const { response } = await this.api.parseError(error);
      if (response.isControlGroupError) {
        throw response;
      }
      // Unless it's a control group error, we want to pass back error info
      // so we can render it on the GenerateCredentialsDatabase component
      return response;
    }
  },

  async getAwsRole(backend, id) {
    try {
      const { data } = await this.api.secrets.awsReadRole(id, backend);
      return data;
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
