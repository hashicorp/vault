/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import ControlGroupConfigForm from 'vault/forms/control-group-config';

export default class ControlGroupsConfigureRoute extends Route {
  @service api;
  @service version;

  async beforeModel() {
    await this.version.fetchFeatures();
    return super.beforeModel(...arguments);
  }

  async model() {
    if (!this.version.hasControlGroups) return null;
    let maxTtl;
    try {
      const { data } = await this.api.sys.systemReadConfigControlGroup();
      maxTtl = data?.max_ttl;
    } catch (e) {
      // if you haven't saved a config, the API 404s, so start with an empty form
      const { status } = await this.api.parseError(e);
      if (status !== 404) throw e;
    }
    return new ControlGroupConfigForm({ max_ttl: maxTtl ?? 0 }, { isNew: !maxTtl });
  }
}
