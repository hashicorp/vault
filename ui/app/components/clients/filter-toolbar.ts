/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached, tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';
import { capitalize } from '@ember/string';

import { ClientFilters, type ClientFilterTypes } from 'core/utils/client-count-utils';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  filterQueryParams: Record<ClientFilterTypes, string>;
  // Dataset objects technically have more keys than the client filter types, but at minimum they contain ClientFilterTypes
  dataset: Record<ClientFilterTypes, string>[];
  onFilter: CallableFunction;
}

// Correspond to each search input's tracked variable in the component class
type SearchProperty = 'namespacePathSearch' | 'mountPathSearch' | 'mountTypeSearch';

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = Object.values(ClientFilters);

  // Tracked filter values
  @tracked namespace_path: string;
  @tracked mount_path: string;
  @tracked mount_type: string;

  // Tracked search inputs
  @tracked namespacePathSearch = '';
  @tracked mountPathSearch = '';
  @tracked mountTypeSearch = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { namespace_path, mount_path, mount_type } = this.args.filterQueryParams;
    this.namespace_path = namespace_path || '';
    this.mount_path = mount_path || '';
    this.mount_type = mount_type || '';
  }

  get anyFilters() {
    return Object.values(this.filterProps).some((v) => !!v);
  }

  @cached
  get dropdownItems() {
    const namespacePaths = new Set<string>();
    const mountPaths = new Set<string>();
    const mountTypes = new Set<string>();

    // iterate over dataset once to get dropdown items
    this.args.dataset.forEach((d) => {
      // namespace_path for root is technically an empty string, so convert to 'root'
      const namespace = d.namespace_path === '' ? 'root' : d.namespace_path;
      if (namespace) namespacePaths.add(namespace);
      if (d.mount_path) mountPaths.add(d.mount_path);
      if (d.mount_type) mountTypes.add(d.mount_type);
    });

    return {
      [ClientFilters.NAMESPACE]: [...namespacePaths],
      [ClientFilters.MOUNT_PATH]: [...mountPaths],
      [ClientFilters.MOUNT_TYPE]: [...mountTypes],
    };
  }

  @cached
  get dropdownConfig() {
    return {
      [ClientFilters.NAMESPACE]: {
        label: 'namespace',
        dropdownItems: this.dropdownItems[ClientFilters.NAMESPACE],
        searchProperty: 'namespacePathSearch',
      },
      [ClientFilters.MOUNT_PATH]: {
        label: 'mount path',
        dropdownItems: this.dropdownItems[ClientFilters.MOUNT_PATH],
        searchProperty: 'mountPathSearch',
      },
      [ClientFilters.MOUNT_TYPE]: {
        label: 'mount type',
        dropdownItems: this.dropdownItems[ClientFilters.MOUNT_TYPE],
        searchProperty: 'mountTypeSearch',
      },
    };
  }

  // It's possible that a query param may not exist in the dropdown, in which case show an alert
  get filterAlert() {
    const alert = (label: string, filter: string) =>
      `${capitalize(label)} "${filter}" not found in the current data.`;
    return this.filterTypes
      .flatMap((f: ClientFilters) => {
        const filterValue = this.filterProps[f];
        const inDropdown = this.dropdownItems[f].includes(filterValue);
        return !inDropdown && filterValue ? [alert(this.dropdownConfig[f].label, filterValue)] : [];
      })
      .join(' ');
  }

  // the cached decorator recomputes this getter every time the tracked properties
  // update instead of every time it is accessed
  @cached
  get filterProps() {
    return this.filterTypes.reduce(
      (obj, filterType) => {
        obj[filterType] = this[filterType];
        return obj;
      },
      {} as Record<ClientFilterTypes, string>
    );
  }

  @action
  handleFilterSelect(filterProperty: ClientFilterTypes, value: string, close: CallableFunction) {
    this[filterProperty] = value;
    close();
  }

  @action
  handleDropdownClose(searchProperty: SearchProperty) {
    // reset search input for that dropdown
    this.updateSearch(searchProperty, '');
    this.applyFilters();
  }

  @action
  clearFilters(filterProperty: ClientFilterTypes | '') {
    if (filterProperty) {
      this[filterProperty] = '';
    } else {
      this.namespace_path = '';
      this.mount_path = '';
      this.mount_type = '';
    }
    this.applyFilters();
  }

  @action
  applyFilters() {
    // Fire callback so URL query params match selected filters
    this.args.onFilter(this.filterProps);
  }

  @action
  handleSearch(event: HTMLElementEvent<HTMLInputElement>) {
    const { value, id: searchProperty } = event.target;
    debounce(this, this.updateSearch, searchProperty as SearchProperty, value, 50);
  }

  @action
  updateSearch(searchProperty: SearchProperty, searchValue: string) {
    this[searchProperty] = searchValue;
  }

  // TEMPLATE HELPERS
  searchDropdown = (dropdownItems: string[], searchProperty: SearchProperty) => {
    const searchInput = this[searchProperty];
    return searchInput
      ? dropdownItems.filter((i) => i?.toLowerCase().includes(searchInput.toLowerCase()))
      : dropdownItems;
  };

  noItemsMessage = (searchValue: string, label: string) => {
    return searchValue ? `No matching ${label}` : `No ${label} to filter`;
  };
}
