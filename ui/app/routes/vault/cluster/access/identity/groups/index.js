/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { paginate } from 'core/utils/paginate-list';
import { IdentityApiGroupListByIdListEnum } from '@hashicorp/vault-client-typescript';

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
      // Fetch group list
      const response = await this.api.identity.groupListById(IdentityApiGroupListByIdListEnum.TRUE);
      const items = this.api.keyInfoToArray(response);

      // Build capability paths for all items
      const capabilityPaths = items.flatMap((item) => [
        this.capabilities.pathFor('identityCapabilities', { identityType: 'group', id: item.id }),
        this.capabilities.pathFor('groupAlias'),
      ]);

      // Fetch group details (type and alias) and capabilities in parallel
      const [capabilitiesMap, detailsResponses] = await Promise.all([
        this.capabilities.fetch(capabilityPaths),
        Promise.all(items.map((item) => this.api.identity.groupReadById(item.id))),
      ]);

      // Build a details map for quick lookup
      const detailsMap = {};
      detailsResponses.forEach((res) => {
        if (res?.data) {
          detailsMap[res.data.id] = { type: res.data.type, alias: res.data.alias };
        }
      });

      const groupAliasPath = this.capabilities.pathFor('groupAlias');
      const groupAliasCapabilities = capabilitiesMap[groupAliasPath];

      const itemsWithCapabilities = items.map((item) => {
        const capPath = this.capabilities.pathFor('identityCapabilities', {
          identityType: 'group',
          id: item.id,
        });
        const itemCapabilities = capabilitiesMap[capPath];

        return {
          ...item,
          ...(detailsMap[item.id] || {}),
          canDelete: itemCapabilities?.canDelete || false,
          canEdit: itemCapabilities?.canUpdate || false,
          canAddAlias: groupAliasCapabilities?.canCreate || false,
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
      identityType: 'group',
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
