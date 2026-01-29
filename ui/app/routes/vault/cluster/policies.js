/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
const ALLOWED_TYPES = ['acl', 'egp', 'rgp'];

export default class PoliciesRoute extends Route {
  @service version;
  @service router;

  beforeModel() {
    return this.version.fetchFeatures().then(() => {
      return super.beforeModel(...arguments);
    });
  }

  model(params) {
    const policyType = params.type;
    if (!ALLOWED_TYPES.includes(policyType)) {
      return this.router.transitionTo(this.routeName, ALLOWED_TYPES[0]);
    }
    return {};
  }
}
