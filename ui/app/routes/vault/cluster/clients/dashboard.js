/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import getStorage from 'vault/lib/token-storage';
import { inject as service } from '@ember/service';
import timestamp from 'core/utils/timestamp';

export default class DashboardRoute extends Route {
  @service store;
  currentDate = timestamp.now().toISOString();

  async getActivity(start_time) {
    // on init ONLY make network request if we have a start_time
    return start_time
      ? await this.store.queryRecord('clients/activity', {
          start_time: { timestamp: start_time },
          end_time: { timestamp: this.currentDate },
        })
      : {};
  }

  async getLicenseStartTime() {
    try {
      const license = await this.store.queryRecord('license', {});
      // if license.startTime is 'undefined' return 'null' for consistency
      return license.startTime || getStorage().getItem('vault:ui-inputted-start-date') || null;
    } catch (e) {
      // return null so user can input date manually
      // if already inputted manually, will be in localStorage
      return getStorage().getItem('vault:ui-inputted-start-date') || null;
    }
  }

  async model() {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients');
    const licenseStart = await this.getLicenseStartTime();
    const activity = await this.getActivity(licenseStart);
    return {
      config,
      versionHistory,
      activity,
      licenseStartTimestamp: licenseStart,
      currentDate: this.currentDate,
    };
  }
}
