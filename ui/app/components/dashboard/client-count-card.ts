/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import timestamp from 'core/utils/timestamp';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import {
  destructureClientCounts,
  formatByMonths,
  formatByNamespace,
} from 'core/utils/client-counts/serializers';

import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import Owner from '@ember/owner';
import type { InternalClientActivityReadConfigurationResponse } from '@hashicorp/vault-client-typescript';
import { HTMLElementEvent } from 'vault/forms';
import type {
  Activity,
  ByNamespaceClients,
  NamespaceObject,
  Counts,
  ActivityMonthBlock,
} from 'vault/client-counts/activity-api';

/**
 * @module DashboardClientCountCard
 * DashboardClientCountCard component are used to display total and new client count information
 *
 * @example
 * <Dashboard::ClientCountCard />
 */

export default class DashboardClientCountCard extends Component<object> {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  @tracked activityData: Activity | null = null;
  @tracked activityConfig: InternalClientActivityReadConfigurationResponse | null = null;
  @tracked canUpdateActivityConfig = true;
  @tracked updatedAt = '';

  constructor(owner: Owner, args: object) {
    super(owner, args);
    this.fetchClientActivity.perform();
  }

  get currentMonthActivityTotalCount() {
    const byMonth = this.activityData?.by_month;
    return byMonth?.[byMonth.length - 1]?.new_clients.clients;
  }

  get statSubText() {
    let formattedStart, formattedEnd;
    if (this.activityData) {
      const { start_time, end_time } = this.activityData;
      formattedStart = start_time ? parseAPITimestamp(start_time, 'MMM yyyy') : null;
      formattedEnd = end_time ? parseAPITimestamp(end_time, 'MMM yyyy') : null;
    }
    return formattedStart && formattedEnd
      ? {
          total: `The number of clients in this billing period (${formattedStart} - ${formattedEnd}).`,
          new: 'The number of clients new to Vault in the current month.',
        }
      : { total: 'No total client data available.', new: 'No new client data available.' };
  }

  fetchClientActivity = task(
    waitFor(async (e?: HTMLElementEvent<HTMLInputElement>) => {
      if (e) e.preventDefault();
      this.updatedAt = timestamp.now().toISOString();
      this.activityData = null;
      this.activityConfig = null;

      try {
        const response = await this.api.sys.internalClientActivityReportCounts();
        if (response) {
          this.activityData = {
            ...response,
            by_namespace: formatByNamespace(response.by_namespace as NamespaceObject[] | null),
            by_month: formatByMonths(response.months as ActivityMonthBlock[]),
            total: destructureClientCounts(response.total as ByNamespaceClients | Counts),
          };
        }
      } catch (error) {
        // used for rendering the "No data" empty state, swallow any errors requesting config data
        this.activityConfig = await this.api.sys.internalClientActivityReadConfiguration().catch(() => null);
        // Clients::NoData needs to know if the user can update the config
        const { canUpdate } = await this.capabilities.for('clientsConfig');
        this.canUpdateActivityConfig = canUpdate;
      }
    })
  );
}
