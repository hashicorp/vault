/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import ClusterRouteBase from './cluster-route-base';
import config from 'vault/config/environment';
import { isEmptyValue } from 'core/helpers/is-empty-value';

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
      // log user in directly (i.e. no form interaction) via URL query param
      const authResponse = await this.unwrapToken(wrapped_token, clusterModel.id);
      return { clusterModel, unwrapResponse: authResponse };
    }

    const visibleAuthMounts = await this.fetchMounts();
    return {
      clusterModel,
      storedLoginData: this.auth.getAuthType(),
      visibleAuthMounts,
    };
  }

  resetController(controller) {
    controller.set('authMethod', 'token');
  }

  afterModel(model) {
    if (model?.unwrapResponse) {
      // handles the transition
      return this.controllerFor('vault.cluster.auth').send('authSuccess', model.unwrapResponse);
    }
    if (config.welcomeMessage) {
      this.flashMessages.info(config.welcomeMessage, {
        sticky: true,
        priority: 300,
      });
    }
  }

  // authenticates the user if the wrapped_token query param exists
  async unwrapToken(token, clusterId) {
    try {
      const { auth } = await this.api.sys.unwrap({}, this.api.buildHeaders({ token }));
      return await this.auth.authenticate({
        clusterId,
        backend: 'token',
        data: { token: auth.clientToken },
        selectedAuth: 'token',
      });
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.controllerFor('vault.cluster.auth').unwrapTokenError = message;
    }
  }

  async fetchMounts() {
    try {
      const resp = await this.api.sys.internalUiListEnabledVisibleMounts(
        this.api.buildHeaders({ token: '' })
      );
      // return a falsy value if the object is empty
      return isEmptyValue(resp.auth) ? null : resp.auth;
    } catch {
      // swallow the error if there's an error fetching mount data (i.e. invalid namespace)
      return null;
    }
  }
}
