/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import ClusterRouteBase from './cluster-route-base';
import config from 'vault/config/environment';
import { isEmptyValue } from 'core/helpers/is-empty-value';
import { supportedTypes } from 'vault/utils/supported-login-methods';
import { sanitizePath } from 'core/utils/sanitize-path';

export default class AuthRoute extends ClusterRouteBase {
  queryParams = {
    authMount: { replace: true, refreshModel: true },
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
    const authMount = params?.authMount;

    return {
      clusterModel,
      visibleAuthMounts,
      directLinkData: authMount ? this.getMountOrTypeData(authMount, visibleAuthMounts) : null,
    };
  }

  afterModel() {
    if (config.welcomeMessage) {
      this.flashMessages.info(config.welcomeMessage, {
        sticky: true,
        priority: 300,
      });
    }
  }

  redirect(model, transition) {
    if (model?.unwrapResponse) {
      // handles the transition
      return this.controllerFor('vault.cluster.auth').send('authSuccess', model.unwrapResponse);
    }
    const hasQueryParam = transition.to?.queryParams?.with;
    const isInvalid = !model.directLinkData;
    if (hasQueryParam && isInvalid) {
      // redirect user and clear out the query param if it's invalid
      this.router.replaceWith(this.routeName, { queryParams: { authMount: null } });
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

  /*
    In older versions of Vault, the "with" query param could refer to either the auth mount path or the type
    (which may be the same, since the default mount path *is* the type). 
    For backward compatibility, we handle both scenarios.
    → If `authMount` matches a visible auth mount, return its mount data (which includes the type).
    → If it matches a supported auth type instead, return just the type to preselect it in the dropdown.
  */
  getMountOrTypeData(authMount, visibleAuthMounts) {
    if (visibleAuthMounts?.[authMount]) {
      return { path: authMount, ...visibleAuthMounts[authMount], isVisibleMount: true };
    }
    const types = supportedTypes(this.version.isEnterprise);
    if (types.includes(sanitizePath(authMount))) {
      return { type: authMount, isVisibleMount: false };
    }
    // `type` is necessary because it determines which login fields to render.
    // If we can't safely glean it from the query param, ignore it and return null
    return null;
  }
}
