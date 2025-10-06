/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

import type { ByMonthNewClients, TotalClients } from 'vault/vault/client-counts/activity-api';
import type FlagsService from 'vault/services/flags';

interface Args {
  byMonthNewClients: ByMonthNewClients[];
  runningTotals: TotalClients;
}

export default class RunningTotal extends Component<Args> {
  @service declare readonly flags: FlagsService;

  @tracked showStacked = false;

  get chartContainerText() {
    return `The total clients in the specified date range, displayed per month. This includes entity, non-entity${
      this.flags.secretsSyncIsActivated ? ', ACME and secrets sync clients' : ' and ACME clients'
    }. The total client count number is an important consideration for Vault billing.`;
  }

  get runningTotalData() {
    return this.args.byMonthNewClients.map((monthly) => ({
      ...monthly,
      new_clients: monthly.clients,
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
    return [{ key: 'new_clients', label: 'New clients' }];
  }
}
