/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class ControlGroupsRoute extends Route {
  @service api;
  @service capabilities;
  @service version;

  async beforeModel() {
    await this.version.fetchFeatures();
    return super.beforeModel(...arguments);
  }

  async model() {
    if (!this.version.hasControlGroups) return null;
    const { canUpdate } = await this.capabilities.for('controlGroupConfig');
    return { canConfigure: canUpdate };
  }

  resetController(controller, isExiting) {
    if (isExiting) controller.accessor = '';
  }
}
