/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class IdentityIndexRoute extends Route {
  @service pagination;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  model(params) {
    const itemType = this.modelFor('vault.cluster.access.identity');
    const modelType = `identity/${itemType}`;
    return this.pagination
      .lazyPaginatedQuery(modelType, {
        responsePath: 'data.keys',
        page: params.page,
        pageFilter: params.pageFilter,
        sortBy: 'name',
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { pageFilter } = this.paramsFor(this.routeName);
    controller.setProperties({
      filter: pageFilter || '',
      page: resolvedModel?.meta?.currentPage || 1,
      identityType: this.modelFor('vault.cluster.access.identity'),
    });
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('filter', null);
    }
  }

  @action
  willTransition(transition) {
    window.scrollTo(0, 0);
    if (transition.targetName !== this.routeName) {
      this.pagination.clearDataset();
    }
    return true;
  }

  @action
  reload() {
    this.pagination.clearDataset();
    this.refresh();
  }
}
