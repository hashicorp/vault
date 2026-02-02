/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import keys from 'core/utils/keys';
import errorMessage from 'vault/utils/error-message';

import type CapabilitiesService from 'vault/services/capabilities';
import type Router from 'vault/router';
import type NamespaceService from 'vault/services/namespace';
import type { NamespaceOption } from 'vault/services/namespace';
import type AuthService from 'vault/services/auth';
import type Store from '@ember-data/store';

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
  @service declare capabilities: CapabilitiesService;
  @service declare namespace: NamespaceService;
  @service declare router: Router;
  @service declare store: Store;

  // Load 200 namespaces in the namespace picker at a time
  @tracked batchSize = 200;

  @tracked canManageNamespaces = false; // Show/hide manage namespaces button
  @tracked errorLoadingNamespaces = '';
  @tracked hasNamespaces = false;
  @tracked searchInput = '';
  @tracked searchInputHelpText =
    "Enter a full path in the search bar and hit the 'Enter' ↵ key to navigate faster.";

  constructor(owner: unknown, args: Record<string, never>) {
    super(owner, args);
    this.loadOptions();
    this.fetchManageCapability();
  }

  get allNamespaces(): NamespaceOption[] {
    return this.namespace.getOptions();
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
  async fetchManageCapability(): Promise<void> {
    // The namespace picker options are from `sys/internal/ui/namespaces` which all users have permissions to request.
    // The UI view for managing namespaces (i.e. CRUD actions) calls `sys/namespaces` and DOES require LIST permissions.
    // This is the capability check to hide/show the button that navigates to that route.
    const { canList } = await this.capabilities.fetchPathCapabilities('sys/namespaces');
    this.canManageNamespaces = canList;
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
