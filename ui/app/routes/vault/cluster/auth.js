/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import ClusterRouteBase from './cluster-route-base';
import config from 'vault/config/environment';

export default class AuthRoute extends ClusterRouteBase {
  queryParams = {
    wrapped_token: { refreshModel: true },
    authMethod: { replace: true },
  };

  @service api;
  @service auth;
  @service flashMessages;
  @service version;

  beforeModel() {
    return super.beforeModel().then(() => {
      return this.version.fetchFeatures();
    });
  }

  async model(params, transition) {
    const wrapped_token = transition?.to?.queryParams?.wrapped_token;
    if (wrapped_token) {
      const clusterModel = this.modelFor('vault.cluster');
      await this.unwrapToken(wrapped_token, clusterModel.id);
    }
    return super.model(...arguments);
  }

  resetController(controller) {
    controller.set('authMethod', 'token');
  }

  afterModel() {
    if (config.welcomeMessage) {
      this.flashMessages.info(config.welcomeMessage, {
        sticky: true,
        priority: 300,
      });
    }
  }

  async unwrapToken(token, clusterId) {
    const authController = this.controllerFor('vault.cluster.auth');
    try {
      const { auth } = await this.api.sys.unwrap({}, this.api.buildHeaders({ token }));
      const authResponse = await this.auth.authenticate({
        clusterId,
        backend: 'token',
        data: { token: auth.clientToken },
        selectedAuth: 'token',
      });
      // handles transition
      return authController.send('authSuccess', authResponse);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      authController.unwrapTokenError = message;
    }
  }
}
