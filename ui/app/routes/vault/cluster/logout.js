/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Route.extend(ModelBoundaryRoute, {
  auth: service(),
  session: service(),
  controlGroup: service(),
  flashMessages: service(),
  console: service(),
  permissions: service(),
  namespaceService: service('namespace'),
  router: service(),

  modelTypes: computed(function () {
    return ['secret', 'secret-engine'];
  }),

  getAuthType() {
    const selectedAuth = localStorage.getItem('selectedAuth');
    if (selectedAuth) return selectedAuth;
    // fallback to authData which discerns backend type from token
    return this.auth.authData ? this.auth.authData.backend?.type : null;
  },

  beforeModel({ to: { queryParams } }) {
    const authType = this.getAuthType();
    const ns = this.namespaceService.path;
    this.controlGroup.deleteTokens();
    this.namespaceService.reset();
    this.console.set('isOpen', false);
    this.console.clearLog(true);
    this.flashMessages.clearMessages();
    this.permissions.reset();

    this.session.invalidate();

    queryParams.with = authType;
    if (ns) {
      queryParams.namespace = ns;
    }
    if (Ember.testing) {
      // Don't redirect on the test
      this.replaceWith('vault.cluster.auth', { queryParams });
    } else {
      const { cluster_name } = this.paramsFor('vault.cluster');
      location.assign(this.router.urlFor('vault.cluster.auth', cluster_name, { queryParams }));
    }
  },
});
