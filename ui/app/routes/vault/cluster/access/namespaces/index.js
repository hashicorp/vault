/**
 * Copyright IBM Corp. 2016, 2025
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
