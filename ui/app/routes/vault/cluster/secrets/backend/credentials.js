/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import ControlGroupError from 'vault/lib/control-group-error';
import { later } from '@ember/runloop';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';

const SUPPORTED_DYNAMIC_BACKENDS = ['database', 'ssh', 'aws', 'totp'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',
  pathHelp: service('path-help'),
  router: service(),
  store: service(),

  beforeModel() {
    const { id: backendPath, type: backendType } = this.modelFor('vault.cluster.secrets.backend');
    // redirect if the backend type does not support credentials
    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendType)) {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root', backendPath);
    }
    // hydrate model if backend type is ssh
    if (backendType === 'ssh') {
      this.pathHelp.hydrateModel('ssh-otp-credential', backendPath);
    }
  },

  getDatabaseCredential(backend, secret, roleType = '') {
    return this.store.queryRecord('database/credential', { backend, secret, roleType }).catch((error) => {
      if (error instanceof ControlGroupError) {
        throw error;
      }
      // Unless it's a control group error, we want to pass back error info
      // so we can render it on the GenerateCredentialsDatabase component
      const status = error?.httpStatus;
      let title;
      let message = `We ran into a problem and could not continue: ${
        error?.errors ? error.errors[0] : 'See Vault logs for details.'
      }`;
      if (status === 403) {
        // 403 is forbidden
        title = 'You are not authorized';
        message =
          "Role wasn't found or you do not have permissions. Ask your administrator if you think you should have access.";
      }
      return {
        errorHttpStatus: status,
        errorTitle: title,
        errorMessage: message,
      };
    });
  },

  async getAwsRole(backend, id) {
    try {
      const role = await this.store.queryRecord('role-aws', { backend, id });
      return role;
    } catch (e) {
      // swallow error, non-essential data
      return;
    }
  },

  getTOTPCode(backend, keyName) {
    return this.store.adapterFor('totp').generateCode(backend, keyName);
  },

  async getTOTPKey(backend, keyName) {
    try {
      const key = await this.store.queryRecord('totp', { id: keyName, backend });
      return key;
    } catch (e) {
      // swallow error, non-essential data
      return;
    }
  },

  async model(params) {
    const role = params.secret;
    const { id: backendPath, type: backendType } = this.modelFor('vault.cluster.secrets.backend');
    const roleType = params.roleType;
    let dbCred, awsRole, totpCode;
    if (backendType === 'database') {
      dbCred = await this.getDatabaseCredential(backendPath, role, roleType);
    } else if (backendType === 'aws') {
      awsRole = await this.getAwsRole(backendPath, role);
    } else if (backendType === 'totp') {
      totpCode = await this.getTOTPCode(backendPath, role);
      totpCode.period = (await this.getTOTPKey(backendPath, role))?.period;
    }

    return resolve({
      backendPath,
      backendType,
      roleName: role,
      roleType,
      dbCred,
      awsRoleType: awsRole?.credentialType,
      totpCode,
    });
  },

  async afterModel(model, transition) {
    if (model.backendType === 'totp' && model.totpCode.period) {
      later(
        () => {
          this.refresh();
        },
        (model.totpCode.period - (getUnixTime(timestamp.now()) % model.totpCode.period)) * 1000
      );
    }
  },

  resetController(controller) {
    controller.reset();
  },

  actions: {
    willTransition() {
      // we do not want to save any of the credential information in the store.
      // once the user navigates away from this page, remove all credential info.
      this.store.unloadAll('database/credential');
    },
  },
});
