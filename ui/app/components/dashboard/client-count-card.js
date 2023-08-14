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
 * <Dashboard::ClientCountCard  />
 * ```
 */

export default class DashboardClientCountCard extends Component {
  currentDate = timestamp.now().toISOString();

  @service store;

  @tracked startDate = null;
  @tracked activityData = null;

  constructor() {
    super(...arguments);
    this.fetchClientActivity.perform();
  }

  @task
  @waitFor
  *getActivity(start_time) {
    // on init ONLY make network request if we have a start_time
    return start_time
      ? yield this.store.queryRecord('clients/activity', {
          start_time: { timestamp: start_time },
          end_time: { timestamp: this.currentDate },
        })
      : {};
  }

  @task
  @waitFor
  *getLicenseStartTime() {
    try {
      const license = yield this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || getStorage().getItem('vault:ui-inputted-start-date') || null;
    } catch (e) {
      // return null so user can input date manually
      // if already inputted manually, will be in localStorage
      return getStorage().getItem('vault:ui-inputted-start-date') || null;
    }
  }

  @task
  @waitFor
  *fetchClientActivity() {
    try {
      this.startDate = yield this.getLicenseStartTime.perform();
      this.activityData = yield this.getActivity.perform(this.startDate);
    } catch (error) {
      this.error = error;
    }
  }
}
