/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached, tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { debounce } from '@ember/runloop';
import { capitalize } from '@ember/string';
import { buildISOTimestamp, parseAPITimestamp } from 'core/utils/date-formatters';
import { ClientFilters } from 'core/utils/client-counts/helpers';

import type {
  ActivityExportData,
  ClientFilterTypes,
  MountClients,
} from 'vault/vault/client-counts/activity-api';
import type { HTMLElementEvent } from 'vault/forms';

interface Args {
  dataset: ActivityExportData[] | MountClients[];
  dropdownMonths?: string[];
  filterQueryParams: Record<ClientFilterTypes, string>;
  isExportData?: boolean;
  onFilter: CallableFunction;
}

// Correspond to each search input's tracked variable in the component class
type SearchProperty = 'namespacePathSearch' | 'mountPathSearch' | 'mountTypeSearch' | 'monthSearch';

export default class ClientsFilterToolbar extends Component<Args> {
  filterTypes = Object.values(ClientFilters);

  // Tracked filter values
  @tracked namespace_path: string;
  @tracked mount_path: string;
  @tracked mount_type: string;
  @tracked month: string;

  // Tracked search inputs
  @tracked namespacePathSearch = '';
  @tracked mountPathSearch = '';
  @tracked mountTypeSearch = '';
  @tracked monthSearch = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { namespace_path, mount_path, mount_type, month } = this.args.filterQueryParams;
    this.namespace_path = namespace_path || '';
    this.mount_path = mount_path || '';
    this.mount_type = mount_type || '';
    this.month = month || '';
  }

  get anyFilters() {
    return Object.values(this.filterProps).some((v) => !!v);
  }

  @cached
  get dropdownItems() {
    const namespacePaths = new Set<string>();
    const mountPaths = new Set<string>();
    const mountTypes = new Set<string>();
    const months = new Set<string>();

    if (this.args.dataset) {
      // iterate over dataset once to get dropdown items
      this.args.dataset.forEach((d) => {
        // namespace_path for root is technically an empty string, so convert to 'root'
        const namespace = d.namespace_path === '' ? 'root' : d.namespace_path;
        if (namespace) namespacePaths.add(namespace);
        if (d.mount_path) mountPaths.add(d.mount_path);
        if (d.mount_type) mountTypes.add(d.mount_type);
        // `client_first_used_time` only exists for the dataset rendered in the "Client list" tab (ActivityExportData),
        // and if the client ID was initially used in version 1.21 or later.
        if ('client_first_used_time' in d && d.client_first_used_time) {
          // for now, we're only concerned with month granularity so we want the dropdown filter to contain an ISO timestamp
          // of the first of the month for each client_first_used_time
          const date = parseAPITimestamp(d.client_first_used_time) as Date;
          const year = date.getUTCFullYear();
          const monthIdx = date.getUTCMonth();
          const timestamp = buildISOTimestamp({ year, monthIdx, isEndDate: false });
          months.add(timestamp);
        }
      });
    }

    return {
      [ClientFilters.NAMESPACE]: [...namespacePaths],
      [ClientFilters.MOUNT_PATH]: [...mountPaths],
      [ClientFilters.MOUNT_TYPE]: [...mountTypes],
      // The "Overview tab" manually passes an array of months
      [ClientFilters.MONTH]: this.args.dropdownMonths || [...months],
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
      [ClientFilters.MONTH]: {
        label: 'month',
        dropdownItems: this.dropdownItems[ClientFilters.MONTH],
        searchProperty: 'monthSearch',
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
        // Don't show an alert for the "Month" filter because it doesn't match dataset values one to one
        if (ClientFilters.MONTH === f) return [];
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
      this.month = '';
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
  formatTimestamp = (isoTimestamp: string) => parseAPITimestamp(isoTimestamp, 'MMMM yyyy') as string;

  searchDropdown = (dropdownItems: string[], searchProperty: SearchProperty) => {
    const searchInput = this[searchProperty];

    if (searchInput) {
      return dropdownItems.filter((i) => {
        const isMatch = (item: string) => item?.toLowerCase().includes(searchInput.toLowerCase());
        // For months, search both the ISO timestamp and formatted display value (e.g., "January 2024")
        return searchProperty === 'monthSearch' ? isMatch(i) || isMatch(this.formatTimestamp(i)) : isMatch(i);
      });
    }

    return dropdownItems;
  };

  noItemsMessage = (searchValue: string, label: string) => {
    if (searchValue) return `No matching ${label}`;

    // The version upgrade message is only relevant if the toolbar filtering activity export data
    // because that is when the months dropdown is populated by the `client_first_used_time` key.
    return label === 'months' && this.args.isExportData
      ? 'Filtering by month is only available for clients initially used after upgrading to version 1.21.'
      : `No ${label} to filter`;
  };
}
