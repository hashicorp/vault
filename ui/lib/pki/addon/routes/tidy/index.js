/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from '../overview';
import timestamp from 'core/utils/timestamp';

export default class PkiTidyIndexRoute extends Route {
  @service api;
  @service secretMountPath;

  async fetchTidyStatus() {
    const status = await this.api.secrets.pkiTidyStatus(this.secretMountPath.currentPath);
    return { ...status, responseTimestamp: timestamp.now() };
  }
  async model() {
    const { hasConfig, autoTidyConfig, engine } = this.modelFor('tidy');
    const tidyStatus = await this.fetchTidyStatus();

    return {
      tidyStatus,
      hasConfig,
      autoTidyConfig,
      engine,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;

    controller.tidyStatus = resolvedModel.tidyStatus;
    controller.fetchTidyStatus = this.fetchTidyStatus;
    controller.pollTidyStatus.perform();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.pollTidyStatus.cancelAll();
    }
  }
}
