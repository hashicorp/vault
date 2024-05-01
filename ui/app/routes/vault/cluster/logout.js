/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { computed } from '@ember/object';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Route.extend(ModelBoundaryRoute, {
  auth: service(),
  controlGroup: service(),
  flashMessages: service(),
  console: service(),
  permissions: service(),
  namespaceService: service('namespace'),
  router: service(),
  version: service(),
  customMessages: service(),

  modelTypes: computed(function () {
    return ['secret', 'secret-engine'];
  }),

  beforeModel({ to: { queryParams } }) {
    const authType = this.auth.getAuthType();
    const ns = this.namespaceService.path;
    this.auth.deleteCurrentToken();
    this.controlGroup.deleteTokens();
    this.namespaceService.reset();
    this.console.set('isOpen', false);
    this.console.clearLog(true);
    this.flashMessages.clearMessages();
    this.permissions.reset();
    this.version.version = null;

    if (this.version.isEnterprise) {
      this.customMessages.clearCustomMessages();
    }

    queryParams.with = authType;
    if (ns) {
      queryParams.namespace = ns;
    }
    if (Ember.testing) {
      // Don't redirect on the test
      this.router.replaceWith('vault.cluster.auth', { queryParams });
    } else {
      const { cluster_name } = this.paramsFor('vault.cluster');
      location.assign(this.router.urlFor('vault.cluster.auth', cluster_name, { queryParams }));
    }
  },
});
