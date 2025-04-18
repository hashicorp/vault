/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { sanitizePath } from 'core/utils/sanitize-path';
import type FlagsService from 'vault/services/flags';
import { HTMLElementEvent } from 'vault/forms';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @tracked selectedMonth = '';

  constructor(owner: unknown, args: any) {
    super(owner, args);
    this.selectedMonth = '';
  }
  get hasAttributionData() {
    // we hide attribution data when mountPath filter present
    // or if there's no data
    if (this.args.mountPath || !this.totalUsageCounts.clients) return false;
    return true;
  }

  // mounts attribution
  get namespaceMountAttribution() {
    const { activity } = this.args;
    const nsLabel = this.namespacePathForFilter;
    return activity?.byNamespace?.find((ns) => sanitizePath(ns.label) === nsLabel)?.mounts || [];
  }

  get months() {
    return this.byMonthActivityData.map((m) => ({
      display: parseAPITimestamp(m.timestamp, 'MMMM yyyy'),
      value: m.month,
    }));
  }

  get tableData() {
    if (!this.selectedMonth) return false;
    // get data from selected month
    const monthData = this.byMonthActivityData.find((m) => m.month === this.selectedMonth);
    const namespaces = monthData?.new_clients.namespaces;

    let data: object[] = [];
    // iterate over namespaces to add "namespace" to each mount object
    namespaces?.forEach((n) => {
      const mounts = n.mounts.map((m) => {
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
