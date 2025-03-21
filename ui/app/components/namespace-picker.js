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

  @tracked selected = this.namespace?.currentNamespace;
  @tracked options = [];
  // @tracked groupedOptions = {groupName: '', options: []} // SHANNONTODO: All Namespaces Label

  constructor() {
    super(...arguments);
    this.loadOptions();
  }

  @action
  async loadOptions() {
    await this.namespace?.findNamespacesForUser.perform();
    this.options = ['root', ...(this.namespace?.accessibleNamespaces || [])];
    // this.groupedOptions.options = this.options; // SHANNONTODO: All Namespaces Label
    // this.groupedOptions.groupName = `All namespaces (${this.options.length})`; // SHANNONTODO: All Namespaces Label
  }

  @action
  handleChange(selectedOption) {
    window.location.href = getNamespaceLink(window.location, selectedOption);
  }
}

// PRIVATE Helper Functions

function getNamespaceLink(location, namespace) {
  const origin = getOrigin(location);
  const encodedNamespace = encodeURIComponent(namespace);

  // The full URL/origin is required so that the page is reloaded.
  return `${origin}/ui/vault/dashboard?namespace=${encodedNamespace}`;
}

function getOrigin(location) {
  return location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '');
}
