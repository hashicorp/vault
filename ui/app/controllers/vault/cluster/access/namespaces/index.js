/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import Controller from '@ember/controller';
import keys from 'core/utils/keys';

/**
 * @module ManageNamespaces
 * ManageNamespacesController is the controller for the
 * vault.cluster.access.namespaces.index route.
 */
export default class ManageNamespacesController extends Controller {
  queryParams = ['pageFilter', 'page'];

  @service namespace;
  @service router;

  @tracked query;
  @tracked pageFilter = '';

  constructor() {
    super(...arguments);
    this.query = this.pageFilter;
  }

  get accessibleNamespaces() {
    return this.namespace.accessibleNamespaces;
  }

  get currentNamespace() {
    return this.namespace.path;
  }

  navigate(pageFilter) {
    const route = 'vault.cluster.access.namespaces.index';
    const args = [route, { queryParams: { page: 1, pageFilter: pageFilter || null } }];
    this.router.transitionTo(...args);
  }

  @action
  handleKeyDown(event) {
    const isEscKeyPressed = keys.ESC.includes(event.key);
    if (isEscKeyPressed) {
      // On escape, transition to roles index route.
      this.navigate();
    }
    // ignore all other key events
  }

  @action handleInput(evt) {
    this.query = evt.target.value;
  }

  @action
  handleSearch(evt) {
    evt.preventDefault();
    this.navigate(this.query);
  }

  @action
  refreshNamespaceList() {
    // fetch new namespaces for the namespace picker
    this.namespace.findNamespacesForUser.perform();
    this.send('reload');
  }
}
