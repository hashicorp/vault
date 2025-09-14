/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type AdapterError from 'vault/@ember-data/adapter/error';
import type { ClientsCountsRouteModel } from '../counts';
import type Store from '@ember-data/store';

export default class ClientsCountsClientListRoute extends Route {
  @service declare readonly store: Store;

  async fetchAndFormatExportData(startTimestamp: string | undefined, endTimestamp: string | undefined) {
    const adapter = this.store.adapterFor('clients/activity');
    let exportData, exportError;
    try {
      const resp = await adapter.exportData({
        // the API only accepts json or csv
        format: 'json',
        start_time: startTimestamp,
        end_time: endTimestamp,
      });
      const jsonLines = await resp.text();
      const lines = jsonLines.trim().split('\n');
      exportData = lines.map((line: string) => JSON.parse(line));
    } catch (error) {
      // Ideally we would not handle errors manually but this is the pattern the other client.counts
      // route follow since the sys/internal/counters API doesn't always return helpful error messages.
      // When these routes are migrated away from ember data we should revisit the error handling.
      exportError = error as AdapterError;
    }
    return { exportData, exportError };
  }

  async model() {
    const { startTimestamp, endTimestamp } = this.modelFor(
      'vault.cluster.clients.counts'
    ) as ClientsCountsRouteModel;

    const { exportData, exportError } = await this.fetchAndFormatExportData(startTimestamp, endTimestamp);
    return { exportData, exportError };
  }
}
