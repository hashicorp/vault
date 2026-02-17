/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import clearModelCache from 'vault/utils/shared-model-boundary';

export default class LogoutRoute extends Route {
  @service auth;
  @service store;
  @service controlGroup;
  @service flashMessages;
  @service console;
  @service permissions;
  @service('namespace') namespaceService;
  @service router;
  @service version;
  @service customMessages;

  modelTypes = ['secret', 'secret-engine'];

  beforeModel({ to: { queryParams } }) {
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
  }
  deactivate() {
    clearModelCache(this.store, this.modelTypes);
  }
}
