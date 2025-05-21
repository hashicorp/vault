/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent, { Args } from '../activity';
import Ember from 'ember';
import { service } from '@ember/service';
import type FlagsService from 'vault/services/flags';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'vault/forms';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { MountClients } from 'core/utils/client-count-utils';
import RouterService from '@ember/routing/router-service';
import { paginate } from 'core/utils/paginate-list';

interface TableData extends MountClients {
  namespace: string;
}

type TableColumn = 'namespace' | 'label' | 'mount_type' | 'clients';
type SortDirection = 'asc' | 'desc';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedMonth = '';
  @tracked currentPage = 1;
  @tracked sortColumn: TableColumn = 'clients';
  @tracked sortDirection: SortDirection = 'desc';

  pageSize = Ember.testing ? 3 : 10; // lower in tests to test pagination without seeding more data

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.setTrackedFromArgs();
  }

  setTrackedFromArgs() {
    if (this.args.month) {
      this.selectedMonth = this.args.month;
    }
  }

  get hasAttributionData() {
    // we hide attribution table when mountPath filter present
    // or if there's no data
    return !this.args.mountPath && this.totalUsageCounts.clients;
  }

  get months() {
    return this.byMonthNewClients.map((m) => ({
      display: parseAPITimestamp(m.timestamp, 'MMMM yyyy'),
      value: m.month,
    }));
  }

  get tableData(): TableData[] {
    if (!this.selectedMonth) return [];
    // get data from selected month
    const monthData = this.byMonthNewClients.find((m) => m.month === this.selectedMonth);
    const namespaces = monthData?.namespaces;

    let data: TableData[] = [];
    // iterate over namespaces to add "namespace" to each mount object
    namespaces?.forEach((n) => {
      const mounts: TableData[] = n.mounts.map((m) => {
        // add namespace to mount block
        return { ...m, namespace: n.label };
      });
      data = [...data, ...mounts];
    });

    const sorted = this.sortTableData(data);
    return sorted;
  }

  get paginatedTableData(): TableData[] {
    const paginated = paginate(this.tableData, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });

    return paginated;
  }

  sortTableData(data: TableData[]): TableData[] {
    if (this.sortColumn) {
      data = [...data].sort((a, b) => {
        const valA = a[this.sortColumn];
        const valB = b[this.sortColumn];

        if (valA < valB) return this.sortDirection === 'asc' ? -1 : 1;
        if (valA > valB) return this.sortDirection === 'asc' ? 1 : -1;
        return 0;
      });
    }
    return data;
  }

  get tableHeaderMessage(): string {
    return this.selectedMonth
      ? 'No data is available for the selected month'
      : 'Select a month to view client attribution';
  }

  get tableBodyMessage(): string {
    return this.selectedMonth
      ? 'View the namespace mount breakdown of clients by selecting another month.'
      : 'View the namespace mount breakdown of clients by selecting a month.';
  }

  @action
  updateSort(column: TableColumn, direction: SortDirection) {
    this.sortColumn = column;
    this.sortDirection = direction;
  }

  @action
  selectMonth(e: HTMLElementEvent<HTMLInputElement>) {
    this.selectedMonth = e.target.value;
    this.args.updateQueryParams({ month: this.selectedMonth });
  }

  @action
  onPageChange(page: number) {
    this.currentPage = page;
  }
}
