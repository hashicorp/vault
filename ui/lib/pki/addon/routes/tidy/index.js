/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from '../overview';
import { hash } from 'rsvp';
import { service } from '@ember/service';
import timestamp from 'core/utils/timestamp';

export default class PkiTidyIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  async fetchTidyStatus() {
    const adapter = this.store.adapterFor('application');
    const tidyStatusResponse = await adapter.ajax(
      `/v1/${this.secretMountPath.currentPath}/tidy-status`,
      'GET'
    );
    const responseTimestamp = timestamp.now();
    tidyStatusResponse.data.responseTimestamp = responseTimestamp;
    return tidyStatusResponse.data;
  }

  model() {
    const { hasConfig, autoTidyConfig, engine } = this.modelFor('tidy');

    return hash({
      tidyStatus: this.fetchTidyStatus(),
      hasConfig,
      autoTidyConfig,
      engine,
    });
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
