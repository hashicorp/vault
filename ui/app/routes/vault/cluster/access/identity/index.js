/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { paginate } from 'core/utils/paginate-list';
import { fetchIdentityItemsWithCapabilities } from 'vault/utils/identity-helpers';

export default class IdentityIndexRoute extends Route {
  @service api;
  @service capabilities;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  async model(params) {
    const { pageFilter, page } = params;
    const identityType = this.modelFor('vault.cluster.access.identity');

    try {
      const itemsWithCapabilities = await fetchIdentityItemsWithCapabilities({
        identityType,
        api: this.api,
        capabilities: this.capabilities,
      });

      return paginate(itemsWithCapabilities, { page, filter: pageFilter });
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return [];
      }
      throw error;
    }
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
  reload() {
    this.refresh();
  }
}
