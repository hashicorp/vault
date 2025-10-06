/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { cached, tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { filterTableData, flattenMounts } from 'core/utils/client-count-utils';

import type ClientsActivityModel from 'vault/vault/models/clients/activity';
import type { ClientFilterTypes } from 'vault/vault/client-counts/activity-api';

export interface Args {
  activity: ClientsActivityModel;
  onFilterChange: CallableFunction;
  filterQueryParams: Record<ClientFilterTypes, string>;
}

export default class ClientsOverviewPageComponent extends Component<Args> {
  @tracked selectedMonth = '';

  @cached
  get byMonthNewClients() {
    return this.args.activity.byMonth?.map((m) => m?.new_clients) || [];
  }

  // Supplies data passed to dropdown filters (except months which is computed below )
  get activityData() {
    // If no month is selected the table displays all of the activity for the queried date range.
    const selectedMonth = this.args.filterQueryParams.month;
    const namespaceData = selectedMonth
      ? this.byMonthNewClients.find((m) => m.timestamp === selectedMonth)?.namespaces
      : this.args.activity.byNamespace;

    // Get the array of "mounts" data nested in each namespace object and flatten
    return flattenMounts(namespaceData || []);
  }

  @cached
  get months() {
    return this.byMonthNewClients.reverse().map((m) => m.timestamp);
  }

  get tableData() {
    if (this.activityData?.length) {
      // Reset the `month` query param because it determines which dataset (see this.activityData)
      // is passed to the table and is does not filter for key/value pairs within this dataset.
      const filters = { ...this.args.filterQueryParams };
      filters.month = '';
      return filterTableData(this.activityData, filters);
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
  handleFilter(filters: Record<ClientFilterTypes, string>) {
    this.args.onFilterChange(filters);
  }
}
