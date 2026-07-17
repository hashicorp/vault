/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { cached } from '@glimmer/tracking';
import { debounce } from '@ember/runloop';
import { duration } from 'core/helpers/format-duration';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import mapOrderStatus from 'pki/helpers/map-order-status';

import type { OrdersIndexRouteParams, RecentOrderListItem } from 'pki/routes/external/orders/index';
import type RouterService from '@ember/routing/router-service';
import type { HTMLElementEvent } from 'vault/forms';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface Args {
  model: {
    engine: SecretsEngineResource;
    recentOrders: RecentOrderListItem[];
    query: OrdersIndexRouteParams;
  };
}

export default class ExternalPkiPageRecentOrdersComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  // Dropdown filters
  @tracked roleFilter = '';
  @tracked roleNameSearch = '';
  @tracked orderIdSearch = '';
  @tracked statusFilter = '';
  // Lookup for empty state
  @tracked orderIdLookup = '';

  // Options for ?within= query parameter
  // Min allowed is 1hr, max is 1 week (168h)
  timeQueryOptions = [
    { label: '1 hour', value: '1h' },
    { label: '1 day', value: '24h' },
    { label: '3 days', value: '72h' },
    { label: '5 days', value: '120h' },
    { label: '1 week', value: '168h' },
  ];

  tableColumns = [
    {
      key: 'order_id',
      label: 'Order ID',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'role_name',
      label: 'Role',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'order_status',
      label: 'Status',
      isSortable: true,
      customTableItem: true,
      width: '150px',
    },
    {
      key: 'creation_date',
      label: 'Created',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'last_update',
      label: 'Last updated',
      isSortable: true,
      customTableItem: true,
    },
  ];

  get filterCount(): number {
    return [this.roleFilter, this.statusFilter].filter(Boolean).length;
  }

  get timeQuery() {
    // `within` comes from the route model and is always a duration string,
    // either from the URL query param or the default value of '1h'.
    const { within } = this.args.model.query;
    const durationString = duration([within]) as string;
    // If duration is just a singular hour/day/week remove the "1 " prefix (e.g., "1 day" → "day")
    // Will not format multi-unit values, e.g. 1 day 2 hours
    const formatted = durationString.replace(/^1 (\w+)$/, '$1');
    return { formatted, raw: within };
  }

  @cached
  get roleNameOptions(): string[] {
    const roleNames = new Set<string>();
    this.args.model.recentOrders.forEach((order) => {
      if (order.role_name) roleNames.add(order.role_name);
    });
    return [...roleNames];
  }

  @cached
  get statusOptions(): string[] {
    const statuses = new Set<string>();
    this.args.model.recentOrders.forEach((order) => {
      if (order.order_status) {
        const mappedState = mapOrderStatus(order.order_status);
        statuses.add(mappedState.text);
      }
    });
    return [...statuses];
  }

  get filteredRoleNames(): string[] {
    if (this.roleNameSearch) {
      return this.roleNameOptions.filter((name) =>
        name.toLowerCase().includes(this.roleNameSearch.toLowerCase())
      );
    }
    return this.roleNameOptions;
  }

  get filteredOrders() {
    return this.args.model.recentOrders.filter((order) => {
      // Filter by order_id search input
      // Order IDs ARE case sensitive, but for filtering purposes we don't need to be that strict
      if (this.orderIdSearch && !order.order_id.includes(this.orderIdSearch.toLowerCase())) {
        return false;
      }
      // Filter by role_name
      if (this.roleFilter && order.role_name !== this.roleFilter) {
        return false;
      }
      // Filter by status
      if (this.statusFilter && mapOrderStatus(order.order_status).text !== this.statusFilter) {
        return false;
      }
      return true;
    });
  }

  get emptyState() {
    switch (true) {
      case !!this.filterCount:
        return {
          title: 'No matching orders',
          message: 'Clear or update filters to view recent orders.',
        };
      case this.orderIdSearch !== '':
        return {
          title: `No recent orders matching: ${this.orderIdSearch}`,
          message: '',
        };
      default:
        return {
          title: 'No recent orders',
          message: `No orders have been created in the last ${this.timeQuery.formatted} (${this.timeQuery.raw}). Select a different time period or lookup an archived order by its ID.`,
          yieldAction: true,
        };
    }
  }

  // Filter handlers
  @action
  dropdownClick(filterKey: 'roleFilter' | 'statusFilter', value: string, close: CallableFunction) {
    this[filterKey] = value;
    close();
  }

  @action
  handleSearch(searchKey: 'orderIdSearch' | 'roleNameSearch', event: HTMLElementEvent<HTMLInputElement>) {
    const { value } = event.target;
    debounce(this, this.updateSearchValue, searchKey, value, 50);
  }

  @action
  updateSearchValue(searchKey: 'orderIdSearch' | 'roleNameSearch', value: string) {
    this[searchKey] = value;
  }

  @action
  clearFilter(filterKey?: 'roleFilter' | 'statusFilter') {
    if (filterKey) {
      this[filterKey] = '';
    } else {
      // Clear all
      this.roleFilter = '';
      this.statusFilter = '';
    }
  }
  // End filter handlers

  @action
  lookupOrder() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.external.orders.order',
      this.args.model.engine.id,
      this.orderIdLookup
    );
  }

  @action
  refresh() {
    this.router.refresh('vault.cluster.secrets.backend.pki.external.orders');
  }

  // TEMPLATE HELPERS
  formatDate = (isoString: string) => parseAPITimestamp(isoString, "MM/dd/yyyy, HH:mm 'UTC'");
}
