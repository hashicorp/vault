/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { withConfig } from 'pki/decorators/check-issuers';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

@withConfig()
export default class ConfigurationIndexRoute extends Route {
  async model() {
    return {
      hasConfig: this.pkiMountHasConfig,
      mountConfig: this.modelFor('application'),
      ...this.modelFor('configuration'),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
  }
}
