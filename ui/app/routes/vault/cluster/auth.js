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
  @service namespace;
  @service store;
  @service version;

  get adapter() {
    return this.store.adapterFor('application');
  }

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

    const loginSettings = this.version.isEnterprise ? await this.fetchLoginSettings() : null;
    const visibleAuthMounts = await this.fetchMounts();
    const authMount = params?.authMount;

    return {
      clusterModel,
      visibleAuthMounts,
      directLinkData: this.getDirectLinkData(authMount, visibleAuthMounts),
      loginSettings,
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

  async fetchLoginSettings() {
    try {
      // TODO update with api service when api-client is updated
      const response = await this.adapter.ajax(
        '/v1/sys/internal/ui/default-auth-methods',
        'GET',
        this.api.buildHeaders({ token: '' })
      );

      if (response?.data) {
        const { default_auth_type, backup_auth_types } = response.data;
        return {
          defaultType: default_auth_type,
          // TODO WIP backend PR consistently return empty array when no backup_auth_types
          backupTypes: backup_auth_types?.length ? backup_auth_types : null,
        };
      }
    } catch {
      // swallow if there's an error and fallback to default login form configuration
      return null;
    }
  }

  async fetchMounts() {
    try {
      const { data } = await this.adapter.ajax(
        '/v1/sys/internal/ui/mounts',
        'GET',
        this.api.buildHeaders({ token: '' })
      );
      // return a falsy value if the object is empty
      return isEmptyValue(data.auth) ? null : data.auth;
    } catch {
      // catch error if there's a problem fetching mount data (i.e. invalid namespace)
      return null;
    }
  }

  /*
    In older versions of Vault, the "with" query param could refer to either the auth mount path or the type
    (which may be the same, since the default mount path *is* the type). 
    For backward compatibility, we handle both scenarios.
    → If `authMount` matches a visible auth mount the method will assume that mount path to login and render as the default in the login form.
    → If `authMount` matches a supported auth type (and the mount does not have `listing_visibility="unauth"`), that type is preselected in the login form.
  */
  getDirectLinkData(authMount, visibleAuthMounts) {
    if (!authMount) return null;

    const sanitizedParam = sanitizePath(authMount); // strip leading/trailing slashes
    // mount paths in visibleAuthMounts always end in a slash, so format for consistency
    const formattedPath = `${sanitizedParam}/`;
    const mountData = visibleAuthMounts?.[formattedPath];
    if (mountData) {
      return { path: formattedPath, type: mountData.type };
    }

    const types = supportedTypes(this.version.isEnterprise);
    if (types.includes(sanitizedParam)) {
      return { type: sanitizedParam };
    }
    // `type` is necessary because it determines which login fields to render.
    // If we can't safely glean it from the query param, ignore it and return null.
    return null;
  }
}
