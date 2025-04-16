/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

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
  @service namespace;
  @service router;
  @service store;

  // Show/hide refresh & manage namespaces buttons
  @tracked hasListPermissions = false;

  @tracked selected = {};
  @tracked options = [];
  @tracked search = '';
  @tracked namespaceLabel = 'All Namespaces';

  @tracked selectedNamespace = () => this.selected?.id || '-';

  constructor() {
    super(...arguments);
    this.loadOptions();
  }

  // TODO make private when converting from js to ts
  #matchesPath(option, currentNamespace) {
    // TODO: Revisit. A hardcoded check for "path" & "/path" seems hacky, but it fixes a breaking test:
    //  "Acceptance | Enterprise | namespaces: it shows nested namespaces if you log in with a namespace starting with a /"
    //  My assumption is that namespace shouldn't start with a "/", but is this a HVD thing? or is the test outdated?
    return option?.path === currentNamespace?.path || `/${option?.path}` === currentNamespace?.path;
  }

  // TODO make private when converting from js to ts
  #getSelected(options, currentNamespace) {
    return options.find((option) => this.#matchesPath(option, currentNamespace));
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
    return [
      // TODO: Some users (including HDS Admin User) should never see the root namespace. Address this in a followup PR.
      { id: 'root', path: '', label: 'root' },
      ...(namespace?.accessibleNamespaces || []).map((ns) => {
        const parts = ns.split('/');
        return { id: parts[parts.length - 1], path: ns, label: ns };
      }),
    ];
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
  async loadOptions() {
    // TODO: namespace service's findNamespacesForUser will never throw an error.
    //  Check with design to determine if we should continue to ignore or handle an error situation here.
    await this.namespace?.findNamespacesForUser.perform();

    this.options = this.#getOptions(this.namespace);
    this.selected = this.#getSelected(this.options, this.namespace);

    await this.fetchListCapability();
  }

  @action
  async onChange(selected) {
    this.selected = selected;
    this.router.transitionTo('vault.cluster.dashboard', { queryParams: { namespace: selected.path } });
  }
}
