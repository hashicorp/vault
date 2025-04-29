/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { KEYS } from 'core/utils/keyboard-keys';

/**
 * @module NamespacePicker
 * @description component is used to display a dropdown listing all namespaces that the current user has access to.
 *  The user can select a namespace from the dropdown to navigate directly to that namespace.
 *  The "Manage" button directs the user to the namespace management page.
 *  The "Refresh List" button refrehes the list of namespaces in the dropdown.
 *
 * @example
 * <NamespacePicker class="hds-side-nav-hide-when-minimized" />
 */

export default class NamespacePicker extends Component {
  @service auth;
  @service namespace;
  @service router;
  @service store;

  // Show/hide refresh & manage namespaces buttons
  @tracked hasListPermissions = false;

  @tracked batchSize = 200;

  @tracked allNamespaces = [];
  @tracked hasNamespaces = false;
  @tracked searchInput = '';
  @tracked searchInputHelpText =
    "Enter a full path in the search bar and hit the 'Enter' ↵ key to navigate faster.";
  @tracked selected = {};

  constructor() {
    super(...arguments);
    this.loadOptions();
  }

  // TODO make private when converting from js to ts
  #matchesPath(option, currentPath) {
    // TODO: Revisit. A hardcoded check for "path" & "/path" seems hacky, but it fixes a breaking test:
    //  "Acceptance | Enterprise | namespaces: it shows nested namespaces if you log in with a namespace starting with a /"
    //  My assumption is that namespace shouldn't start with a "/", but is this a HVD thing? or is the test outdated?
    return option?.path === currentPath || `/${option?.path}` === currentPath;
  }

  // TODO make private when converting from js to ts
  #getSelected(options, currentPath) {
    return options.find((option) => this.#matchesPath(option, currentPath));
  }

  // TODO make private when converting from js to ts
  #getOptions(namespace) {
    /* Each namespace option has 3 properties: { id, path, and label }
     *   - id: node / namespace name (displayed when the namespace picker is closed)
     *   - path: full namespace path (used to navigate to the namespace)
     *   - label: text displayed inside the namespace picker dropdown (if root, then label = id, else label = path)
     *
     *  Example:
     *   | id       | path           | label          |
     *   | ---      | ----           | -----          |
     *   | 'root'   | ''             | 'root'         |
     *   | 'parent' | 'parent'       | 'parent'       |
     *   | 'child'  | 'parent/child' | 'parent/child' |
     */
    const options = [
      ...(namespace?.accessibleNamespaces || []).map((ns) => {
        const parts = ns.split('/');
        return { id: parts[parts.length - 1] || '', path: ns, label: ns };
      }),
    ];

    // Conditionally add the root namespace
    if (this.auth?.authData?.userRootNamespace === '') {
      options.unshift({ id: 'root', path: '', label: 'root' });
    }

    // If there are no namespaces returned by the internal endpoint, add the current namespace
    // to the list of options. This is a fallback for when the user has access to a single namespace.
    if (options.length === 0) {
      options.push({
        id: namespace.currentNamespace,
        path: namespace.path,
        label: namespace.path,
      });
    }

    return options;
  }

  get hasSearchInput() {
    return this.searchInput?.trim().length > 0;
  }

  get namespaceCount() {
    return this.namespaceOptions.length;
  }

  get namespaceLabel() {
    return this.searchInput === '' ? 'All namespaces' : 'Matching namespaces';
  }

  get namespaceOptions() {
    if (this.searchInput.trim() === '') {
      return this.allNamespaces || [];
    } else {
      const filtered = this.allNamespaces.filter((ns) =>
        ns.label.toLowerCase().includes(this.searchInput.toLowerCase())
      );
      return filtered || [];
    }
  }

  get noNamespacesMessage() {
    const noNamespacesMessage = 'No namespaces found.';
    const noMatchingNamespacesHelpText =
      'No matching namespaces found. Try searching for a different namespace.';
    return this.hasSearchInput ? noMatchingNamespacesHelpText : noNamespacesMessage;
  }

  get visibleNamespaceOptions() {
    return this.namespaceOptions.slice(0, this.batchSize);
  }

  @action
  async fetchListCapability() {
    // TODO: Revist. This logic was carried over from previous component implmenetation.
    //  When the user doesn't have this capability, shouldn't we just hide the "Manage" button,
    //  instead of hiding both the "Manage" and "Refresh List" buttons?
    try {
      await this.store.findRecord('capabilities', 'sys/namespaces/');
      this.hasListPermissions = true;
    } catch (e) {
      // If error out on findRecord call it's because you don't have permissions
      // and therefore don't have permission to manage namespaces
      this.hasListPermissions = false;
    }
  }

  @action
  focusSearchInput(element) {
    // On mount, cursor should default to the search input field
    element.focus();
  }

  @action
  async loadOptions() {
    // TODO: namespace service's findNamespacesForUser will never throw an error.
    // Check with design to determine if we should continue to ignore or handle an error situation here.
    await this.namespace?.findNamespacesForUser?.perform();

    this.allNamespaces = this.#getOptions(this.namespace);
    this.selected = this.#getSelected(this.allNamespaces, this.namespace?.path);

    await this.fetchListCapability();
  }

  @action
  loadMore() {
    this.batchSize += 200; // Increase the batch size to load more items
  }

  @action
  setupScrollListener(element) {
    element.addEventListener('scroll', this.onScroll);
  }

  @action
  onScroll(event) {
    const element = event.target;

    // Check if the user has scrolled to the bottom
    if (element.scrollTop + element.clientHeight >= element.scrollHeight) {
      this.loadMore();
    }
  }

  @action
  async onChange(selected) {
    this.selected = selected;
    this.searchInput = '';
    this.router.transitionTo('vault.cluster.dashboard', { queryParams: { namespace: selected.path } });
  }

  @action
  async onKeyDown(event) {
    if (event.key === KEYS.ENTER && this.searchInput?.trim()) {
      const matchingNamespace = this.allNamespaces.find((ns) => ns.label === this.searchInput.trim());

      if (matchingNamespace) {
        this.selected = matchingNamespace;
        this.searchInput = '';
        this.router.transitionTo('vault.cluster.dashboard', {
          queryParams: { namespace: matchingNamespace.path },
        });
      }
    }
  }

  @action
  onSearchInput(event) {
    this.searchInput = event.target.value;
  }

  @action
  async refreshList() {
    this.searchInput = '';
    await this.loadOptions();
  }
}
