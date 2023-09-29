/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import getStorage from 'vault/lib/token-storage';
import timestamp from 'core/utils/timestamp';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

/**
 * @module DashboardClientCountCard
 * DashboardClientCountCard component are used to display total and new client count information
 *
 * @example
 * ```js
 * <Dashboard::ClientCountCard @license={{@model.license}}  />
 * ```
 *  @param {object} license - license object passed from the parent
 */

export default class DashboardClientCountCard extends Component {
  @service store;

  @tracked activityData = null;
  @tracked clientConfig = null;
  @tracked updatedAt = timestamp.now().toISOString();

  constructor() {
    super(...arguments);
    this.fetchClientActivity.perform();
    this.clientConfig = this.store.queryRecord('clients/config', {}).catch(() => {});
  }

  get currentMonthActivityTotalCount() {
    return this.activityData?.byMonth?.lastObject?.new_clients.clients;
  }

  get licenseStartTime() {
    return this.args.license.startTime || getStorage().getItem('vault:ui-inputted-start-date') || null;
  }

  @task
  @waitFor
  *fetchClientActivity(e) {
    if (e) e.preventDefault();
    this.updatedAt = timestamp.now().toISOString();
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
