/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Controller from '@ember/controller';
import { service } from '@ember/service';

/**
 * @module ManageNamespaces
 * ManageNamespacesController is the controller for the
 * vault.cluster.access.namespaces.index route.
 *
 * @param {object} namespaces - list of namespaces
 * @param {string} pageFilter - value of queryParam
 * @param {string} page - value of queryParam
 */

export default class ManageNamespacesController extends Controller {
  @service router;

  constructor() {
    super(...arguments);
  }

  @action
  navigate(pageFilter) {
    const route = 'vault.cluster.access.namespaces.index';
    const args = [route, { queryParams: { page: 1, pageFilter: pageFilter || null } }];
    this.router.transitionTo(...args);
  }

  @action
  refreshRoute() {
    this.send('reload');
  }
}
