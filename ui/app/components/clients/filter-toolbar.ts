/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached, tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';
import { ClientFilters, type ClientFilterTypes, filterIsSupported } from 'core/utils/client-count-utils';

import type { HTMLElementEvent } from 'vault/forms';
interface Args {
  appliedFilters: Record<ClientFilterTypes, string>;
  // the dataset objects have more keys than the client filter types, but at minimum they have ClientFilterTypes
  dataset: Record<ClientFilterTypes, string>[];
  onFilter: CallableFunction;
}

type SearchProperty = 'namespacePathSearch' | 'mountPathSearch' | 'mountTypeSearch';

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = ClientFilters;

  @tracked namespace_path: string;
  @tracked mount_path: string;
  @tracked mount_type: string;

  @tracked namespacePathSearch = '';
  @tracked mountPathSearch = '';
  @tracked mountTypeSearch = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { namespace_path, mount_path, mount_type } = this.args.appliedFilters;
    this.namespace_path = namespace_path || '';
    this.mount_path = mount_path || '';
    this.mount_type = mount_type || '';
  }

  get anyFilters() {
    return (
      Object.keys(this.args.appliedFilters).every((f) => filterIsSupported(f)) &&
      Object.values(this.args.appliedFilters).some((v) => !!v)
    );
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
      [this.filterTypes.NAMESPACE]: [...namespacePaths],
      [this.filterTypes.MOUNT_PATH]: [...mountPaths],
      [this.filterTypes.MOUNT_TYPE]: [...mountTypes],
    };
  }

  @cached
  get dropdownConfig() {
    return {
      [this.filterTypes.NAMESPACE]: {
        label: 'namespace',
        dropdownItems: this.dropdownItems[this.filterTypes.NAMESPACE],
        searchProperty: 'namespacePathSearch',
      },
      [this.filterTypes.MOUNT_PATH]: {
        label: 'mount path',
        dropdownItems: this.dropdownItems[this.filterTypes.MOUNT_PATH],
        searchProperty: 'mountPathSearch',
      },
      [this.filterTypes.MOUNT_TYPE]: {
        label: 'mount type',
        dropdownItems: this.dropdownItems[this.filterTypes.MOUNT_TYPE],
        searchProperty: 'mountTypeSearch',
      },
    };
  }

  @action
  updateFilter(filterProperty: ClientFilterTypes, value: string, close: CallableFunction) {
    this[filterProperty] = value;
    close();
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
    // Fire callback so URL query params update when filters are cleared
    this.applyFilters();
  }

  @action
  applyFilters() {
    this.args.onFilter({
      namespace_path: this.namespace_path,
      mount_path: this.mount_path,
      mount_type: this.mount_type,
    });
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
