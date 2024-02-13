/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),
  type: '',

  enginePathParam() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },

  async fetchConnection(queryOptions) {
    try {
      return await this.store.query('database/connection', queryOptions);
    } catch (e) {
      return e.httpStatus;
    }
  },

  async fetchAllRoles(queryOptions) {
    try {
      return await this.store.query('database/role', queryOptions);
    } catch (e) {
      return e.httpStatus;
    }
  },

  pathQuery(backend, endpoint) {
    return {
      id: `${backend}/${endpoint}/`,
    };
  },

  async fetchCapabilitiesRole(queryOptions) {
    return this.store.queryRecord('capabilities', this.pathQuery(queryOptions.backend, 'roles'));
  },

  async fetchCapabilitiesStaticRole(queryOptions) {
    return this.store.queryRecord('capabilities', this.pathQuery(queryOptions.backend, 'static-roles'));
  },

  async fetchCapabilitiesConnection(queryOptions) {
    return this.store.queryRecord('capabilities', this.pathQuery(queryOptions.backend, 'config'));
  },

  model() {
    const backend = this.enginePathParam();
    const queryOptions = { backend, id: '' };

    const connection = this.fetchConnection(queryOptions);
    const role = this.fetchAllRoles(queryOptions);
    const roleCapabilities = this.fetchCapabilitiesRole(queryOptions);
    const staticRoleCapabilities = this.fetchCapabilitiesStaticRole(queryOptions);
    const connectionCapabilities = this.fetchCapabilitiesConnection(queryOptions);

    return hash({
      backend,
      connections: connection,
      roles: role,
      engineType: 'database',
      id: backend,
      roleCapabilities,
      staticRoleCapabilities,
      connectionCapabilities,
      icon: 'database',
    });
  },

  setupController(controller, model) {
    this._super(...arguments);
    const showEmptyState = model.connections === 404 && model.roles === 404;
    const noConnectionCapabilities =
      !model.connectionCapabilities.canList &&
      !model.connectionCapabilities.canCreate &&
      !model.connectionCapabilities.canUpdate;

    const emptyStateMessage = function () {
      if (noConnectionCapabilities) {
        return 'You cannot yet generate credentials.  Ask your administrator if you think you should have access.';
      } else {
        return 'You can connect an external database to Vault.  We recommend that you create a user for Vault rather than using the database root user.';
      }
    };
    controller.set('showEmptyState', showEmptyState);
    controller.set('emptyStateMessage', emptyStateMessage());
  },
});
