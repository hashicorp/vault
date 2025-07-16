/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'vault/forms';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { formatTableData, TableData } from 'core/utils/client-count-utils';
import type FlagsService from 'vault/services/flags';
import type RouterService from '@ember/routing/router-service';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @service('app-router') declare readonly router: RouterService;

  @tracked selectedMonth = '';

  get hasAttributionData() {
    // we hide attribution table when mountPath filter present
    // or if there's no data
    return !this.args.mountPath && this.totalUsageCounts.clients;
  }

  get months() {
    return this.byMonthNewClients.reverse().map((m) => ({
      display: parseAPITimestamp(m.timestamp, 'MMMM yyyy'),
      value: m.month,
    }));
  }

  get tableData(): TableData[] | undefined {
    if (!this.selectedMonth) return undefined;
    return formatTableData(this.byMonthNewClients, this.selectedMonth);
  }

  @action
  selectMonth(e: HTMLElementEvent<HTMLInputElement>) {
    this.selectedMonth = e.target.value;
  }
}
