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
import { filterTableData, flattenMounts } from 'core/utils/client-count-utils';

import type FlagsService from 'vault/services/flags';
import type RouterService from '@ember/routing/router-service';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedMonth = '';

  @cached
  get byMonthNewClients() {
    return this.args.activity.byMonth?.map((m) => m?.new_clients) || [];
  }

  @cached
  // Supplies data passed to dropdown filters
  get activityData() {
    // Find the namespace data for the selected month
    // If no month is selected the table displays all of the activity for the queried date range
    const namespaceData = this.selectedMonth
      ? this.byMonthNewClients.find((m) => m.timestamp === this.selectedMonth)?.namespaces
      : this.args.activity.byNamespace;

    // Get the array of "mounts" data nested in each namespace object and flatten
    return flattenMounts(namespaceData || []);
  }

  @cached
  get months() {
    return this.byMonthNewClients
      .reverse()
      .map((m) => ({ timestamp: m.timestamp, display: parseAPITimestamp(m.timestamp, 'MMMM yyyy') }));
  }

  get tableData() {
    if (this.activityData?.length) {
      return filterTableData(this.activityData, this.args.filterQueryParams);
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
