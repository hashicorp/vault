/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { set } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  pagination: service(),
  store: service(),

  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  },

  templateName: 'vault/cluster/access/leases/list',

  model(params) {
    const prefix = params.prefix || '';
    if (this.modelFor('vault.cluster.access.leases').canList) {
      return hash({
        leases: this.pagination
          .lazyPaginatedQuery('lease', {
            prefix,
            responsePath: 'data.keys',
            page: params.page,
            pageFilter: params.pageFilter,
          })
          .then((model) => {
            this.set('has404', false);
            return model;
          })
          .catch((err) => {
            if (err.httpStatus === 404 && prefix === '') {
              return [];
            } else {
              throw err;
            }
          }),
        capabilities: hash({
          revokePrefix: this.store.findRecord('capabilities', `sys/leases/revoke-prefix/${prefix}`),
          forceRevokePrefix: this.store.findRecord('capabilities', `sys/leases/revoke-force/${prefix}`),
        }),
      });
    }
  },

  setupController(controller, model) {
    const params = this.paramsFor(this.routeName);
    const prefix = params.prefix ? params.prefix : '';
    const has404 = this.has404;
    controller.set('hasModel', true);
    controller.setProperties({
      model: model.leases,
      capabilities: model.capabilities,
      baseKey: { id: prefix },
      has404,
    });
    if (!has404) {
      const pageFilter = params.pageFilter;
      let filter;
      if (prefix) {
        filter = prefix + (pageFilter || '');
      } else if (pageFilter) {
        filter = pageFilter;
      }
      controller.setProperties({
        filter: filter || '',
        page: model.leases?.meta?.currentPage,
      });
    }
  },

  resetController(controller, isExiting) {
    this._super(...arguments);
    if (isExiting) {
      controller.set('filter', '');
    }
  },

  actions: {
    error(error, transition) {
      const { prefix } = this.paramsFor(this.routeName);

      set(error, 'keyId', prefix);
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      const hasModel = this.controllerFor(this.routeName).hasModel;
      // only swallow the error if we have a previous model
      if (hasModel && error.httpStatus === 404) {
        this.set('has404', true);
        transition.abort();
      } else {
        return true;
      }
    },

    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.pagination.clearDataset();
      }
      return true;
    },
  },
});
