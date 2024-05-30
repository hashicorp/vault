/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { resolve } from 'rsvp';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import ControlGroupError from 'vault/lib/control-group-error';

const SUPPORTED_DYNAMIC_BACKENDS = ['database', 'ssh', 'aws'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',
  pathHelp: service('path-help'),
  router: service(),
  store: service(),

  beforeModel() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    if (backend != 'ssh') {
      return;
    }
    const modelType = 'ssh-otp-credential';
    return this.pathHelp.getNewModel(modelType, backend);
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

  async model(params) {
    const role = params.secret;
    const { id: backendPath, type: backendType } = this.modelFor('vault.cluster.secrets.backend');
    const roleType = params.roleType;
    let dbCred;
    if (backendType === 'database') {
      dbCred = await this.getDatabaseCredential(backendPath, role, roleType);
    }
    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendType)) {
      return this.router.transitionTo('vault.cluster.secrets.backend.list-root', backendPath);
    }
    return resolve({
      backendPath,
      backendType,
      roleName: role,
      roleType,
      dbCred,
    });
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
