/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

const ALLOWED_TYPES = ['acl', 'egp', 'rgp'];

export default Route.extend(ClusterRoute, {
  version: service(),

  beforeModel() {
    return this.version.fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model(params) {
    const policyType = params.type;
    if (!ALLOWED_TYPES.includes(policyType)) {
      return this.router.transitionTo(this.routeName, ALLOWED_TYPES[0]);
    }
    return {};
  },
});
