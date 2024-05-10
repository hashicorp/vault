/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import timestamp from 'core/utils/timestamp';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { parseAPITimestamp } from 'core/utils/date-formatters';

/**
 * @module DashboardClientCountCard
 * DashboardClientCountCard component are used to display total and new client count information
 *
 * @example
 * <Dashboard::ClientCountCard />
 */

export default class DashboardClientCountCard extends Component {
  @service store;

  @tracked activityData = null;
  @tracked hasActivity = false;
  @tracked updatedAt = null;

  constructor() {
    super(...arguments);
    this.fetchClientActivity.perform();
  }

  get currentMonthActivityTotalCount() {
    return this.activityData?.byMonth?.lastObject?.new_clients.clients;
  }

  get statSubText() {
    const format = (date) => parseAPITimestamp(date, 'MMM yyyy');
    const { startTime, endTime } = this.activityData;
    return startTime && endTime
      ? {
          total: `The number of clients in this billing period (${format(startTime)} - ${format(endTime)}).`,
          new: 'The number of clients new to Vault in the current month.',
        }
      : { total: 'No total client data available.', new: 'No new client data available.' };
  }

  @task
  @waitFor
  *fetchClientActivity(e) {
    if (e) e.preventDefault();
    this.updatedAt = timestamp.now().toISOString();

    try {
      this.activityData = yield this.store.queryRecord('clients/activity', {
        current_billing_period: true,
      });
      this.hasActivity = this.activityData.id === 'no-data' ? false : true;
    } catch (error) {
      this.error = error;
    }
  }
}
