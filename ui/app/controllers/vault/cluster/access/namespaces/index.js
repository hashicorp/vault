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
 *
 * @param {object} namespaces - list of namespaces
 * @param {string} pageFilter - value of queryParam
 * @param {string} page - value of queryParam
 */
export default class ManageNamespacesController extends Controller {
  queryParams = ['pageFilter', 'page'];

  // Use namespaceService alias to avoid collision with namespaces
  // input parameter from the route.
  @service('namespace') namespaceService;
  @service router;
  @service flashMessages;

  // The `query` property is used to track the filter
  // input value seperately from updating the `pageFilter`
  // browser query param to prevent unnecessary re-renders.
  @tracked query;
  @tracked pageFilter = '';

  constructor() {
    super(...arguments);
    this.query = this.pageFilter;
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
  async deleteNamespace(nsToDelete) {
    try {
      // Attempt to destroy the record
      await nsToDelete.destroyRecord();

      // Log success and optionally update the UI
      this.flashMessages.success(`Successfully deleted namespace: ${nsToDelete.id}`);

      // Call the refresh method to update the list
      this.refreshNamespaceList();
    } catch (error) {
      this.flashMessages.danger(`There was an error deleting this namespace: ${error.message}`);
    }
  }

  @action
  async refreshNamespaceList() {
    try {
      // Await the async operation to complete
      await this.namespaceService.findNamespacesForUser.perform();
      this.send('reload'); // Trigger the reload only after the task completes
    } catch (error) {
      this.flashMessages.danger('There was an error refreshing the namespace list.');
    }
  }

  @action
  switchNamespace(targetNamespace) {
    this.router.transitionTo('vault.cluster.dashboard', {
      queryParams: { namespace: targetNamespace },
    });
  }
}
