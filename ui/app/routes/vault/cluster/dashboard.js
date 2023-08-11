/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
// eslint-disable-next-line ember/no-mixins
import ClusterRoute from 'vault/mixins/cluster-route';
import getStorage from 'vault/lib/token-storage';
import timestamp from 'core/utils/timestamp';

export default class VaultClusterDashboardRoute extends Route.extend(ClusterRoute) {
  @service store;
  @service version;
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

  async getVaultConfiguration() {
    try {
      const adapter = this.store.adapterFor('application');
      const configState = await adapter.ajax('/v1/sys/config/state/sanitized', 'GET');
      return configState.data;
    } catch (e) {
      return null;
    }
  }

  async model() {
    const vaultConfiguration = this.getVaultConfiguration();
    const licenseStart = await this.getLicenseStartTime();
    const activity = await this.getActivity(licenseStart);

    return hash({
      vaultConfiguration,
      secretsEngines: this.store.query('secret-engine', {}),
      clientCountActivity: activity,
      version: this.version,
    });
  }
}
