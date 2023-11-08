/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import ClusterRouteBase from './cluster-route-base';
import config from 'vault/config/environment';

export default ClusterRouteBase.extend({
  queryParams: {
    authMethod: {
      replace: true,
    },
  },
  flashMessages: service(),
  version: service(),
  session: service(),
  store: service(),

  beforeModel() {
    this.session.prohibitAuthentication('/vault/dashboard');
    return this._super().then(() => {
      return this.version.fetchFeatures();
    });
  },

  async model() {
    const parent = await this._super(...arguments);
    const { wrappedToken } = this.paramsFor('vault');
    if (wrappedToken) {
      // Unwrap wrapped token if present
      // If this is successful, it will be passed to AuthV2::Token component for login
      const adapter = this.store.adapterFor('tools');
      try {
        const response = await adapter.toolAction('unwrap', null, { clientToken: wrappedToken });
        return response.auth;
      } catch (e) {
        return {
          error: `Token unwrap failed: ${e.errors[0]}`,
        };
      }
    }
    return parent;
  },

  resetController(controller) {
    controller.set('wrappedToken', '');
    controller.set('authMethod', 'token');
  },

  afterModel() {
    if (config.welcomeMessage) {
      this.flashMessages.info(config.welcomeMessage, {
        sticky: true,
        priority: 300,
      });
    }
  },
});
