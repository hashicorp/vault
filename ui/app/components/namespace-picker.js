/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// SHANNONTODO: add component docs see configer-wif.ts

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

// PRIVATE Helper Functions
// SHANNONTEST post about this pattern in #vault-ui-devs

function _matchesPath(option, currentNamespace) {
  // SHANNONTODO: Revisit this, this seems hacky, but fixes a breaking test
  // assumption is that namespace shouldn't start with a "/", is this a HVD thing?
  // or is the test outdated?
  return option?.path === currentNamespace?.path || `/${option?.path}` === currentNamespace?.path;
}

function _getSelected(options, currentNamespace) {
  return options.find((option) => _matchesPath(option, currentNamespace));
}

function _getOptions(namespace) {
  return [
    // SHANNONTODO: Add Comment explaining 3 properties because root is blank
    // SHANNONTODO: HDS Admin User (and others?) should never see root
    { id: 'root', path: '', label: 'root' },
    ...(namespace?.accessibleNamespaces || []).map((ns) => {
      const parts = ns.split('/');
      return { id: parts[parts.length - 1], path: ns, label: ns };
    }),
  ];
}

function _getNamespaceLink(location, namespace) {
  const origin = _getOrigin(location);
  const encodedNamespace = encodeURIComponent(namespace.path);

  let queryParams = '';
  if (namespace.path !== '') {
    queryParams = `?namespace=${encodedNamespace}`;
  }

  // The full URL/origin is required so that the page is reloaded.
  return `${origin}/ui/vault/dashboard${queryParams}`;
}

function _getOrigin(location) {
  return location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '');
}

export default class ApplicationComponent extends Component {
  @service namespace;
  @service store;

  // Show/hide refresh & manage namespaces buttons
  @tracked showAfterOptions = false;

  @tracked selected = {};
  @tracked groupedOptions = [{ options: [], groupName: '' }];

  constructor() {
    super(...arguments);
    this.loadOptions();
  }

  // SHANNONTODO this is fetching sys internal namespaces capabilities, is this still needed??
  @action
  async fetchListCapability() {
    try {
      // SHANNONTODO test w/o capabilities
      await this.store.findRecord('capabilities', 'sys/namespaces/');
      this.showAfterOptions = true;
    } catch (e) {
      // If error out on findRecord call it's because you don't have permissions
      // and therefore don't have permission to manage namespaces
      this.showAfterOptions = false;
    }
  }

  @action
  async loadOptions() {
    // SHANNONTODO Currently this never throws an error, should we handle error situation? Check w/ design
    await this.namespace?.findNamespacesForUser.perform();

    const options = _getOptions(this.namespace);
    this.selected = _getSelected(options, this.namespace);
    await this.fetchListCapability();

    // SHANNONTODO: add comment
    this.groupedOptions[0].options = options;
    this.groupedOptions[0].groupName = `All namespaces (${options.length})`;
  }

  @action
  handleChange(selected) {
    // SHANNONTODO add github comment
    window.location.href = _getNamespaceLink(window.location, selected);
  }
}
