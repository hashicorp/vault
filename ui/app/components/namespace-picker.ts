/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import keys from 'core/utils/keys';
import { buildWaiter } from '@ember/test-waiters';

import type Router from 'vault/router';
import type NamespaceService from 'vault/services/namespace';
import type AuthService from 'vault/vault/services/auth';
import type Store from '@ember-data/store';
import errorMessage from 'vault/utils/error-message';

interface NamespaceOption {
  id: string;
  path: string;
  label: string;
}

const waiter = buildWaiter('namespace-picker');

/**
 * @module NamespacePicker
 * @description component is used to display a dropdown listing all namespaces that the current user has access to.
 *  The user can select a namespace from the dropdown to navigate directly to that namespace.
 *  The "Manage" button directs the user to the namespace management page.
 *  The "Refresh List" button refreshes the list of namespaces in the dropdown.
 *
 * @example
 * <NamespacePicker class="hds-side-nav-hide-when-minimized" />
 */
export default class NamespacePicker extends Component {
  @service declare auth: AuthService;
  @service declare namespace: NamespaceService;
  @service declare router: Router;
  @service declare store: Store;

  // Load 200 namespaces in the namespace picker at a time
  @tracked batchSize = 200;

  @tracked canManageNamespaces = false; // Show/hide manage namespaces button
  @tracked canRefreshNamespaces = false; // Show/hide refresh list button
  @tracked errorLoadingNamespaces = '';
  @tracked hasNamespaces = false;
  @tracked searchInput = '';
  @tracked searchInputHelpText =
    "Enter a full path in the search bar and hit the 'Enter' ↵ key to navigate faster.";

  constructor(owner: unknown, args: Record<string, never>) {
    super(owner, args);
    this.loadOptions();
  }

  get allNamespaces(): NamespaceOption[] {
    return this.getOptions(
      this.namespace?.accessibleNamespaces,
      this.namespace?.currentNamespace,
      this.namespace?.path
    );
  }

  get selectedNamespace(): NamespaceOption | null {
    return this.getSelected(this.allNamespaces, this.namespace?.path) ?? null;
  }

  private matchesPath(option: NamespaceOption, currentPath: string): boolean {
    return option?.path === currentPath;
  }

  private getSelected(options: NamespaceOption[], currentPath: string): NamespaceOption | undefined {
    return options.find((option) => this.matchesPath(option, currentPath));
  }

  private getOptions(
    accessibleNamespaces: string[],
    currentNamespace: string,
    path: string
  ): NamespaceOption[] {
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
      ...(accessibleNamespaces || []).map((ns: string) => {
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
        id: currentNamespace,
        path: path,
        label: path,
      });
    }

    return options;
  }

  get hasSearchInput(): boolean {
    return this.searchInput?.trim().length > 0;
  }

  get namespaceCount(): number {
    return this.namespaceOptions.length;
  }

  get namespaceLabel(): string {
    return this.searchInput === '' ? 'All namespaces' : 'Matching namespaces';
  }

  get namespaceOptions(): NamespaceOption[] {
    if (this.searchInput.trim() === '') {
      return this.allNamespaces || [];
    } else {
      const filtered = this.allNamespaces.filter((ns) =>
        ns.label.toLowerCase().includes(this.searchInput.toLowerCase())
      );
      return filtered || [];
    }
  }

  get noNamespacesMessage(): string {
    const noNamespacesMessage = 'No namespaces found.';
    const noMatchingNamespacesHelpText =
      'No matching namespaces found. Try searching for a different namespace.';
    return this.hasSearchInput ? noMatchingNamespacesHelpText : noNamespacesMessage;
  }

  get showNoNamespacesMessage(): boolean {
    const hasError = this.errorLoadingNamespaces !== '';
    return this.namespaceCount === 0 && !hasError;
  }

  get visibleNamespaceOptions(): NamespaceOption[] {
    return this.namespaceOptions.slice(0, this.batchSize);
  }

  @action
  adjustElementWidth(element: HTMLElement): void {
    // Hide the element so that it doesn't affect the layout
    element.style.display = 'none';

    let maxWidth = 240; // Default minimum width
    // Calculate the maximum width of the visible namespace options
    // The namespace is displayed as an HDS::checkmark button, so we need to find the width of the checkmark element
    this.visibleNamespaceOptions.forEach((namespace: NamespaceOption) => {
      const checkmarkElement = document.querySelector(`[data-test-button="${namespace.label}"]`);

      const width = (checkmarkElement as HTMLElement).offsetWidth;
      if (width > maxWidth) {
        maxWidth = width;
      }
    });

    // Set the width of the target element
    element.style.width = `${maxWidth}px`;

    // Show the element once the width is set
    element.style.display = '';
  }

  @action
  async fetchListCapability(): Promise<void> {
    const waiterToken = waiter.beginAsync();
    try {
      const namespacePermission = await this.store.findRecord('capabilities', 'sys/namespaces/');
      this.canRefreshNamespaces = namespacePermission.get('canList');
      this.canManageNamespaces = true;
    } catch (error) {
      // If the findRecord call fails, the user lacks permissions to refresh or manage namespaces.
      this.canRefreshNamespaces = this.canManageNamespaces = false;
    } finally {
      waiter.endAsync(waiterToken);
    }
  }

  @action
  focusSearchInput(element: HTMLInputElement): void {
    // On mount, cursor should default to the search input field
    element.focus();
  }

  @action
  async loadOptions(): Promise<void> {
    try {
      await this.namespace?.findNamespacesForUser?.perform();
      this.errorLoadingNamespaces = '';
    } catch (error) {
      this.errorLoadingNamespaces = errorMessage(error);
    }

    await this.fetchListCapability();
  }

  @action
  loadMore(): void {
    // Increase the batch size to load more items
    this.batchSize += 200;
  }

  @action
  setupScrollListener(element: HTMLElement): void {
    element.addEventListener('scroll', this.onScroll);
  }

  @action
  onScroll(event: Event): void {
    const element = event.target as HTMLElement;

    // Check if the user has scrolled to the bottom
    if (element.scrollTop + element.clientHeight >= element.scrollHeight) {
      this.loadMore();
    }
  }

  @action
  async onChange(selected: NamespaceOption): Promise<void> {
    this.searchInput = '';
    this.router.transitionTo('vault.cluster.dashboard', { queryParams: { namespace: selected.path } });
  }

  @action
  async onKeyDown(event: KeyboardEvent): Promise<void> {
    if (event.key === keys.ENTER && this.searchInput?.trim()) {
      const matchingNamespace = this.allNamespaces.find((ns) => ns.label === this.searchInput.trim());

      if (matchingNamespace) {
        this.searchInput = '';
        this.router.transitionTo('vault.cluster.dashboard', {
          queryParams: { namespace: matchingNamespace.path },
        });
      }
    }
  }

  @action
  onSearchInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    this.searchInput = target.value;
  }

  @action
  async refreshList(): Promise<void> {
    this.searchInput = '';
    await this.loadOptions();
  }

  @action
  toggleNamespacePicker() {
    // Reset the search input when the dropdown is toggled
    this.searchInput = '';
  }
}
