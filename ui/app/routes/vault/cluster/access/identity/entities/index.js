/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { paginate } from 'core/utils/paginate-list';
import { IdentityApiEntityListByIdListEnum } from '@hashicorp/vault-client-typescript';

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

    try {
      // Fetch entity list
      const response = await this.api.identity.entityListById(IdentityApiEntityListByIdListEnum.TRUE);
      const items = this.api.keyInfoToArray(response);

      // Build capability paths for all items plus the alias create path
      const capabilityPaths = [
        ...items.map((item) =>
          this.capabilities.pathFor('identityCapabilities', { identityType: 'entity', id: item.id })
        ),
        this.capabilities.pathFor('groupAlias'),
      ];

      // Fetch capabilities for all items
      const capabilitiesMap = await this.capabilities.fetch(capabilityPaths);

      const aliasPath = this.capabilities.pathFor('groupAlias');
      const aliasCapabilities = capabilitiesMap[aliasPath];

      const itemsWithCapabilities = items.map((item) => {
        const capPath = this.capabilities.pathFor('identityCapabilities', {
          identityType: 'entity',
          id: item.id,
        });
        const itemCapabilities = capabilitiesMap[capPath];

        return {
          ...item,
          canDelete: itemCapabilities?.canDelete || false,
          canEdit: itemCapabilities?.canUpdate || false,
          canAddAlias: aliasCapabilities?.canCreate || false,
        };
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
      identityType: 'entity',
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
