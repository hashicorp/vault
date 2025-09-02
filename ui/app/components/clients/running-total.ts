/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

import type { ByMonthNewClients, TotalClients } from 'core/utils/client-count-utils';
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

  get chartLegend() {
    if (this.showStacked) {
      return [
        { key: 'entity_clients', label: 'entity clients' },
        { key: 'non_entity_clients', label: 'non-entity clients' },
        ...(this.flags.secretsSyncIsActivated ? [{ key: 'secret_syncs', label: 'secret sync clients' }] : []),
        { key: 'acme_clients', label: 'acme clients' },
      ];
    }
    return [{ key: 'new_clients', label: 'new clients' }];
  }
}
