/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiAliasListByIdListEnum } from '@hashicorp/vault-client-typescript';
import { paginate } from 'core/utils/paginate-list';

export default class IdentityAliasesIndexRoute extends Route {
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
    try {
      const response = await this.api.identity.groupListAliasesById(IdentityApiAliasListByIdListEnum.TRUE);
      const aliases = await this.api.keyInfoToArray(response);

      // Build capability paths for all aliases
      const capabilityPaths = aliases.map((alias) =>
        this.capabilities.pathFor('identityCapabilities', {
          identityType: 'group',
          id: alias.id,
        })
      );

      // Fetch capabilities for all aliases
      const capabilitiesMap = await this.capabilities.fetch(capabilityPaths);

      // Attach capabilities to each alias
      const aliasesWithCapabilities = aliases.map((alias) => {
        const aliasCapabilityPath = this.capabilities.pathFor('identityCapabilities', {
          identityType: 'group',
          id: alias.id,
        });
        const aliasCapabilities = capabilitiesMap[aliasCapabilityPath];

        return {
          ...alias,
          canDelete: aliasCapabilities?.canDelete || false,
          canEdit: aliasCapabilities?.canUpdate || false,
        };
      });

      return paginate(aliasesWithCapabilities, { page: params.page, filter: params.pageFilter });
    } catch (err) {
      const { status } = await this.api.parseError(err);
      if (status === 404) {
        return [];
      } else {
        throw err;
      }
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { pageFilter } = this.paramsFor(this.routeName);
    controller.setProperties({
      identityType: 'group',
      filter: pageFilter || '',
      page: resolvedModel?.meta?.currentPage || 1,
    });
  }

  resetController(controller, isExiting) {
    super.resetController(controller, isExiting);
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('filter', null);
    }
  }

  reload() {
    this.refresh();
  }
}
