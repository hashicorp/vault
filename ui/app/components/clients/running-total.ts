/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { toLabel } from 'core/helpers/to-label';

import type { ByMonthClients, ByMonthNewClients, TotalClients } from 'vault/vault/client-counts/activity-api';
import type FlagsService from 'vault/services/flags';
import type VersionService from 'vault/services/version';

interface Args {
  byMonthClients: ByMonthClients[] | ByMonthNewClients[];
  runningTotals: TotalClients;
}

export default class RunningTotal extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  @tracked showStacked = false;

  get chartContainerText() {
    const range = this.version.isEnterprise ? 'billing period' : 'date range';
    return this.flags.isHvdManaged
      ? 'Number of total unique clients in the data period by client type, and total number of unique clients per month. The monthly total is the relevant billing metric.'
      : `Number of clients in the ${range} by client type, and a breakdown of new clients per month during the ${range}. `;
  }

  get dataKey() {
    return this.flags.isHvdManaged ? 'clients' : 'new_clients';
  }

  get runningTotalData() {
    // The parent component determines whether `monthly.clients` in @byMonthClients represents "new" or "total" clients per month.
    // (We render "new" for self-managed clusters and "total" for HVD-managed.)
    // As a result, we do not use `this.dataKey` to select a property from `monthly` but to add a superficial key
    // to the data that ensures the chart tooltip and legend text render appropriately.
    return this.args.byMonthClients.map((monthly) => ({
      ...monthly,
      [this.dataKey]: monthly.clients,
    }));
  }

  get donutChartData() {
    return [
      { value: this.args.runningTotals.entity_clients, label: 'Entity clients' },
      { value: this.args.runningTotals.non_entity_clients, label: 'Non-entity clients' },
      { value: this.args.runningTotals.acme_clients, label: 'ACME clients' },
      ...(this.flags.secretsSyncIsActivated
        ? [{ value: this.args.runningTotals.secret_syncs, label: 'Secret sync clients' }]
        : []),
    ];
  }

  get chartLegend() {
    if (this.showStacked) {
      return [
        { key: 'entity_clients', label: 'Entity clients' },
        { key: 'non_entity_clients', label: 'Non-entity clients' },
        { key: 'acme_clients', label: 'ACME clients' },
        // MUST BE LAST because conditionally renders and legend color mapping for stacked bars will be off otherwise
        ...(this.flags.secretsSyncIsActivated ? [{ key: 'secret_syncs', label: 'Secret sync clients' }] : []),
      ];
    }
    return [{ key: this.dataKey, label: toLabel([this.dataKey]) }];
  }
}
