/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type { ClientsCountsRouteModel } from '../counts';
import type Store from '@ember-data/store';
export default class ClientsCountsClientListRoute extends Route {
  @service declare readonly store: Store;

  async getExportData(startTimestamp: string, endTimestamp: string, namespaceQp?: string) {
    const adapter = this.store.adapterFor('clients/activity');
    return adapter.exportData({
      // the API only accepts json or csv
      format: 'json',
      start_time: startTimestamp,
      end_time: endTimestamp,
      // TODO Wire up namespace param
      namespace: namespaceQp,
    });
  }

  async model() {
    // TODO first check query params, if none fallback on parent route's start/end times
    const { startTimestamp, endTimestamp } = this.modelFor(
      'vault.cluster.clients.counts'
    ) as ClientsCountsRouteModel;
    let formattedData;
    if (startTimestamp && endTimestamp) {
      const exportData = await this.getExportData(startTimestamp, endTimestamp);
      const jsonmaybe = await exportData.text();
      const lines = jsonmaybe.trim().split('\n');
      formattedData = lines.map((line: string) => JSON.parse(line));
    }
    return { startTimestamp, endTimestamp, formattedData };
  }
}
