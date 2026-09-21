/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class PoliciesRoute extends Route {
  @service version;

  beforeModel() {
    return this.version.fetchFeatures().then(() => {
      return super.beforeModel(...arguments);
    });
  }
}
