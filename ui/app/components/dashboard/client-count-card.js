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
import { setStartTimeQuery } from 'core/utils/client-count-utils';
import { dateFormat } from 'core/helpers/date-format';

/**
 * @module DashboardClientCountCard
 * DashboardClientCountCard component are used to display total and new client count information
 *
 * @example
 *
 * <Dashboard::ClientCountCard @isEnterprise={{@version.isEnterprise}} />
 *
 * @param {boolean} isEnterprise - used for setting the start time for the activity log query
 */

export default class DashboardClientCountCard extends Component {
  @service store;

  clientConfig = null;
  licenseStartTime = null;
  @tracked activityData = null;
  @tracked updatedAt = timestamp.now().toISOString();

  constructor() {
    super(...arguments);
    this.fetchClientActivity.perform();
  }

  get currentMonthActivityTotalCount() {
    return this.activityData?.byMonth?.lastObject?.new_clients.clients;
  }

  get statSubText() {
    const format = (date) => dateFormat([date, 'MMM yyyy'], {});
    return this.licenseStartTime
      ? {
          total: `The number of clients in this billing period (${format(this.licenseStartTime)} - ${format(
            this.updatedAt
          )}).`,
          new: 'The number of clients new to Vault in the current month.',
        }
      : { total: 'No total client data available.', new: 'No new client data available.' };
  }

  @task
  @waitFor
  *fetchClientActivity(e) {
    if (e) e.preventDefault();
    this.updatedAt = timestamp.now().toISOString();

    if (!this.clientConfig) {
      // set config and license start time when component initializes
      this.clientConfig = yield this.store.queryRecord('clients/config', {}).catch(() => {});
      this.licenseStartTime = setStartTimeQuery(this.args.isEnterprise, this.clientConfig);
    }

    // only make the network request if we have a start_time
    if (!this.licenseStartTime) return {};
    try {
      this.activityData = yield this.store.queryRecord('clients/activity', {
        start_time: { timestamp: this.licenseStartTime },
        end_time: { timestamp: this.updatedAt },
      });
      this.noActivityData = this.activityData.activity.id === 'no-data' ? true : false;
    } catch (error) {
      this.error = error;
    }
  }
}
