/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { cached, tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'vault/forms';
import { parseAPITimestamp } from 'core/utils/date-formatters';

import { filterTableData, flattenMounts, type MountClients } from 'core/utils/client-count-utils';
import type FlagsService from 'vault/services/flags';
import type RouterService from '@ember/routing/router-service';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedMonth = '';

  @cached
  get clientsByMount() {
    const namespaceData = this.selectedMonth
      ? // Find the namespace data for the selected month
        this.byMonthNewClients.find((m) => m.timestamp === this.selectedMonth)?.namespaces
      : // If no month is selected the table displays all of the by_namespace activity for the queried date range
        this.args.activity.byNamespace;

    // Get the array of "mounts" data nested in each namespace object and flatten
    return flattenMounts(namespaceData || []);
  }

  // DROPDOWN GETTERS
  @cached
  get months() {
    return this.byMonthNewClients
      .reverse()
      .map((m) => ({ timestamp: m.timestamp, display: parseAPITimestamp(m.timestamp, 'MMMM yyyy') }));
  }

  @cached
  get namespaceLabels() {
    return this.args.activity.byNamespace.map((n) => n.label);
  }

  @cached
  get mountPaths() {
    return [...new Set(this.clientsByMount.map((m: MountClients) => m.label))];
  }

  @cached
  get mountTypes() {
    return [...new Set(this.clientsByMount.map((m: MountClients) => m.mount_type))];
  }
  // end dropdown getters

  get tableData() {
    if (this.clientsByMount?.length) {
      return filterTableData(this.clientsByMount, this.args.filterQueryParams);
    }
    return null;
  }

  get tableColumns() {
    return [
      { key: 'namespace_path', label: 'Namespace', isSortable: true },
      { key: 'mount_path', label: 'Mount path', isSortable: true },
      { key: 'mount_type', label: 'Mount type', isSortable: true },
      { key: 'clients', label: 'Client count', isSortable: true },
    ];
  }

  @action
  selectMonth(e: HTMLElementEvent<HTMLInputElement>) {
    this.selectedMonth = e.target.value;
    // Reset filters when no month is selected
    if (this.selectedMonth === '') {
      this.resetFilters();
    }
  }
}
