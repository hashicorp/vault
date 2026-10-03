/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';

@withConfig()
export default class PkiTidyRoute extends Route {
  @service api;

  async model() {
    const engine = this.modelFor('application');
    const autoTidyConfig = await this.api.secrets.pkiReadAutoTidyConfiguration(engine.id);

    return {
      hasConfig: this.pkiMountHasConfig,
      engine,
      autoTidyConfig,
    };
  }
}
