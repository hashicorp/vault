/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import ClusterRouteBase from './cluster-route-base';
import config from 'vault/config/environment';

export default class AuthRoute extends ClusterRouteBase {
  queryParams = {
    authMethod: { replace: true },
    wrapped_token: { refreshModel: true },
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

  async model(params) {
    const clusterModel = this.modelFor('vault.cluster');
    const wrapped_token = params?.wrapped_token;
    if (wrapped_token) {
      await this.unwrapToken(wrapped_token, clusterModel.id);
    }

    const visibleAuthMounts = await this.fetchMounts();
    return {
      clusterModel: clusterModel,
      namespaceInput: this.namespaceInput,
      storedLoginData: this.auth.getAuthType(),
      visibleAuthMounts,
    };
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

  // authenticates the user if the wrapped_token query param exists
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

  async fetchMounts() {
    try {
      const resp = await this.api.sys.internalUiListEnabledVisibleMounts(
        this.api.buildHeaders({ token: '' })
      );
      return resp.auth;
    } catch {
      // swallow the error if there's an error fetching mount data (i.e. invalid namespace)
      return null;
    }
  }
}
