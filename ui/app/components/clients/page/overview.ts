/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent, { Args } from '../activity';
import { service } from '@ember/service';
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
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedMonth = '';
  @tracked currentPage = 1;
  @tracked sortColumn: TableColumn = 'clients';
  @tracked sortDirection: SortDirection = 'desc';

  pageSize = 10;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.setTrackedFromArgs();
  }

  setTrackedFromArgs() {
    if (this.args.month) {
      this.selectedMonth = this.args.month;
    }
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

    // TODO SLW should we put a hard cap still?
    const sorted = this.sortTableData(data);
    return sorted;
  }

  // TODO SLW should these be separate or combine with a variable set for total items?
  get paginatedTableData(): TableData[] {
    const paginated = paginate(this.tableData, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });

    return paginated;
  }

  sortTableData(data: TableData[]): TableData[] {
    if (this.sortColumn) {
      // TODO SLW review sorting logic and see if it works for each column
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
    // TODO SLW confirm copy
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

// TODO SLW are all mounts disabled in mock data? double check
