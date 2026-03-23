/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type FlagsService from 'vault/services/flags';
import type RouterService from '@ember/routing/router-service';
import type VersionService from 'vault/services/version';

export default class ClientsCountsClientListRoute extends Route {
  @service declare readonly flags: FlagsService;
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;

  // The "Client list" tab is only available on enterprise versions
  // The "Client list" tab is hidden on HVD managed clusters (for now) because the "Month" filter for that page
  // uses the `client_first_used_time` timestamp. This timestamp tracks when a client is FIRST seen in the queried date range (i.e. billing period).
  // This is useful for self-managed customers who are billed on monthly NEW clients, but not for HVD users who are billed on TOTAL clients per
  // month regardless of whether the client was seen in a previous month.
  redirect() {
    if (this.version.isCommunity || this.flags.isHvdManaged) {
      this.router.transitionTo('vault.cluster.clients.counts');
    }
  }
}
