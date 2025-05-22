/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import type FlagsService from 'vault/services/flags';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { HTMLElementEvent } from 'vault/forms';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { MountClients } from 'core/utils/client-count-utils';
import RouterService from '@ember/routing/router-service';

interface TableData extends MountClients {
  namespace: string;
}

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

    return data;
  }

  @action
  selectMonth(e: HTMLElementEvent<HTMLInputElement>) {
    this.selectedMonth = e.target.value;
  }
}
