/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import Ember from 'ember';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { paginate } from 'core/utils/paginate-list';
import { MountClients } from 'core/utils/client-count-utils';

interface TableData extends MountClients {
  namespace: string;
}

interface Args {
  data: TableData[];
}

type TableColumn = 'namespace' | 'label' | 'mount_type' | 'clients';
type SortDirection = 'asc' | 'desc';

export default class Table extends Component<Args> {
  @tracked currentPage = 1;
  @tracked sortColumn: TableColumn = 'clients';
  @tracked sortDirection: SortDirection = 'desc';

  pageSize = Ember.testing ? 3 : 10; // lower in tests to test pagination without seeding more data

  get paginatedTableData(): TableData[] {
    const sorted = this.sortTableData(this.args.data);
    const paginated = paginate(sorted, {
      page: this.currentPage,
      pageSize: this.pageSize,
    });

    return paginated;
  }

  get tableHeaderMessage(): string {
    return this.args.data
      ? 'No data is available for the selected month'
      : 'Select a month to view client attribution';
  }

  get tableBodyMessage(): string {
    return this.args.data
      ? 'View the namespace mount breakdown of clients by selecting another month.'
      : 'View the namespace mount breakdown of clients by selecting a month.';
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

  @action
  onPageChange(page: number) {
    this.currentPage = page;
  }

  @action
  updateSort(column: TableColumn, direction: SortDirection) {
    this.sortColumn = column;
    this.sortDirection = direction;
  }
}
