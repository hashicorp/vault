/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
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
  beforeModel() {
    return this._super().then(() => {
      return this.version.fetchFeatures();
    });
  },
  model() {
    return this._super(...arguments);
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
