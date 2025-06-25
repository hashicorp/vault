/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { action } from '@ember/object';
import { buildWaiter } from '@ember/test-waiters';
import { hash } from 'rsvp';

const waiter = buildWaiter('namespace-list-route');

export default class NamespaceListRoute extends Route {
  @service pagination;
  @service store;
  @service version;

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  beforeModel() {
    this.store.unloadAll('namespace');
    return this.version.fetchFeatures().then(() => {
      return super.beforeModel(...arguments);
    });
  }

  async fetchNamespaces(params) {
    const waiterToken = waiter.beginAsync();
    try {
      const model = await this.pagination.lazyPaginatedQuery('namespace', {
        responsePath: 'data.keys',
        page: Number(params?.page) || 1,
        pageFilter: params?.pageFilter,
      });
      return model;
    } catch (err) {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    } finally {
      waiter.endAsync(waiterToken);
    }
  }

  model(params) {
    const { pageFilter } = params;
    return hash({
      namespaces: this.fetchNamespaces(params),
      pageFilter,
    });
  }

  setupController(controller, model) {
    const has404 = this.has404;
    controller.setProperties({
      model: model,
      has404,
      hasModel: true,
    });
    if (!has404) {
      controller.setProperties({
        page: Number(model?.meta?.currentPage) || 1,
      });
    }
  }

  @action
  error(error, transition) {
    /* eslint-disable-next-line ember/no-controller-access-in-routes */
    const hasModel = this.controllerFor(this.routeName).hasModel;
    if (hasModel && error.httpStatus === 404) {
      this.has404 = true;
      transition.abort();
    } else {
      return true;
    }
  }

  @action
  willTransition(transition) {
    window.scrollTo(0, 0);
    if (!transition || transition.targetName !== this.routeName) {
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
