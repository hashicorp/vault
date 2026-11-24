/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class LicenseRoute extends Route {
  @service api;
  @service version;
  @service router;

  beforeModel() {
    if (this.version.isCommunity) {
      this.router.transitionTo('vault.cluster');
    }
  }

  async model() {
    const resp = await this.api.sys.systemReadLicenseStatus();
    const licenseData = resp?.data?.autoloaded ?? {};
    return {
      license_id: licenseData.license_id,
      start_time: licenseData.start_time,
      expiration_time: licenseData.expiration_time,
      features: licenseData.features ?? [],
      performance_standby_count: licenseData.performance_standby_count,
      autoloaded: resp?.data?.autoloading_used,
    };
  }
}
