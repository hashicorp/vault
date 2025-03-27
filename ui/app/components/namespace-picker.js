/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

export default class ApplicationComponent extends Component {
  @service namespace;
  @service store;
  @service router;

  @tracked showAfterOptions = false;
  @tracked selected = {};
  @tracked groupedOptions = [{ options: [], groupName: '' }];

  constructor() {
    super(...arguments);
    this.loadOptions();
  }

  @action
  async fetchListCapability() {
    try {
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
    await this.namespace?.findNamespacesForUser.perform();
    const options = getOptions(this.namespace);
    this.selected = getSelected(options, this.namespace);
    this.fetchListCapability();
    this.groupedOptions[0].options = options;
    this.groupedOptions[0].groupName = `All namespaces (${options.length})`;
  }

  @action
  handleChange(selected) {
    window.location.href = getNamespaceLink(window.location, selected);
  }
}

// PRIVATE Helper Functions

function getSelected(options, currentNamespace) {
  return options.find((option) => matchesPath(option, currentNamespace));
}

function matchesPath(option, currentNamespace) {
  return option?.path === currentNamespace?.path || `/${option?.path}` === currentNamespace?.path;
}

function getOptions(namespace) {
  return [
    { id: 'root', path: '', label: 'root' },
    ...(namespace?.accessibleNamespaces || []).map((ns) => {
      const parts = ns.split('/');
      return { id: parts[parts.length - 1], path: ns, label: ns };
    }),
  ];
}

function getNamespaceLink(location, namespace) {
  const origin = getOrigin(location);
  const encodedNamespace = encodeURIComponent(namespace.path);

  let queryParams = '';
  if (namespace.path !== '') {
    queryParams = `?namespace=${encodedNamespace}`;
  }

  // The full URL/origin is required so that the page is reloaded.
  return `${origin}/ui/vault/dashboard${queryParams}`;
}

function getOrigin(location) {
  return location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '');
}
