/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import keys from 'core/utils/keys';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';
import { DISMISSED_WIZARD_KEY } from '../wizard';

/**
 * @module PageNamespaces
 * PageNamespaces component handles the display and management of namespaces,
 * including the namespace wizard for first-time users.
 *
 * @param {object} namespaces - list of namespaces
 * @param {string} pageFilter - current page filter value
 * @param {function} onFilterChange - callback function to handle filter changes, receives filter string or null to clear
 * @param {function} onRefresh - callback function to refresh the namespace list from the route/controller
 */

interface Args {
  model: {
    namespaces: NamespaceModel[];
    pageFilter: string | null;
  };
  onFilterChange: CallableFunction;
  onRefresh: CallableFunction;
}

interface NamespaceModel {
  id: string;
  destroyRecord: () => Promise<void>;
  [key: string]: unknown;
}

export default class PageNamespacesComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare namespace: NamespaceService;

  // The `query` property is used to track the filter
  // input value separately from updating the `pageFilter`
  // browser query param to prevent unnecessary re-renders.
  @tracked query;
  @tracked nsToDelete = null;
  @tracked hasDismissedWizard = false;

  wizardId = 'namespace';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.query = this.args.model.pageFilter || '';

    // check if the wizard has already been dismissed
    const dismissedWizards = localStorage.getItem(DISMISSED_WIZARD_KEY);
    if (dismissedWizards?.includes(this.wizardId)) {
      this.hasDismissedWizard = true;
    }
  }

  get showWizard() {
    // Show when there are no existing namespaces and it is not in a dismissed state
    return !this.hasDismissedWizard && !this.args.model.namespaces?.length;
  }

  @action
  handleKeyDown(event: KeyboardEvent) {
    const isEscKeyPressed = keys.ESC.includes(event.key);
    if (isEscKeyPressed) {
      // On escape, clear the filter
      this.args.onFilterChange(null);
    }
    // ignore all other key events
  }

  @action
  handleInput(evt: HTMLElementEvent<HTMLInputElement>) {
    this.query = evt.target.value;
  }

  @action
  handleSearch(evt: HTMLElementEvent<HTMLInputElement>) {
    evt.preventDefault();
    this.args.onFilterChange(this.query);
  }

  @action
  async deleteNamespace(nsToDelete: NamespaceModel) {
    try {
      // Attempt to destroy the record
      await nsToDelete.destroyRecord();

      // Log success and optionally update the UI
      this.flashMessages.success(`Successfully deleted namespace: ${nsToDelete.id}`);

      // Call the refresh method to update the list
      this.refreshNamespaceList();
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`There was an error deleting this namespace: ${message}`);
    }
    this.nsToDelete = null;
  }

  @action
  async refreshNamespaceList() {
    try {
      // Await the async operation to complete
      await this.namespace.findNamespacesForUser.perform();
      this.args.onRefresh();
    } catch (error) {
      this.flashMessages.danger('There was an error refreshing the namespace list.');
    }
  }

  @action handlePageChange() {
    this.args.onRefresh();
  }

  @action
  switchNamespace(targetNamespace: string) {
    this.router.transitionTo('vault.cluster.dashboard', {
      queryParams: { namespace: targetNamespace },
    });
  }

  async createNamespace(path: string, header?: string) {
    const headers = header ? this.api.buildHeaders({ namespace: header }) : undefined;
    await this.api.sys.systemWriteNamespacesPath(path, {}, headers);
  }
}
