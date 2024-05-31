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
  pathHelp: service(),
  scope() {
    return this.paramsFor('scope').scope_name;
  },
  beforeModel() {
    return this.pathHelp.getNewModel('kmip/role', this.secretMountPath.currentPath);
  },
  model(params) {
    return this.store
      .lazyPaginatedQuery('kmip/role', {
        backend: this.secretMountPath.currentPath,
        scope: this.scope(),
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
    this._super(...arguments);
    controller.set('scope', this.scope());
  },
});
