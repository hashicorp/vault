/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import ListRoute from 'core/mixins/list-route';
import { service } from '@ember/service';

export default Route.extend(ListRoute, {
  store: service(),
  secretMountPath: service(),
  credParams() {
    const { role_name: role, scope_name: scope } = this.paramsFor('credentials');
    return {
      role,
      scope,
    };
  },
  model(params) {
    const { role, scope } = this.credParams();
    return this.store
      .lazyPaginatedQuery('kmip/credential', {
        role,
        scope,
        backend: this.secretMountPath.currentPath,
        responsePath: 'data.keys',
        page: params.page,
        pageFilter: params.pageFilter,
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  },

  setupController(controller) {
    const { role, scope } = this.credParams();
    this._super(...arguments);
    controller.setProperties({ role, scope });
  },
});
