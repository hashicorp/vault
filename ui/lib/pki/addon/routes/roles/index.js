/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import { getCliMessage } from 'pki/routes/overview';
@withConfig()
export default class PkiRolesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async fetchRoles(params) {
    try {
      const page = Number(params.page) || 1;
      return await this.store.lazyPaginatedQuery('pki/role', {
        backend: this.secretMountPath.currentPath,
        responsePath: 'data.keys',
        page,
        skipCache: page === 1,
      });
    } catch (e) {
      if (e.httpStatus === 404) {
        return { parentModel: this.modelFor('roles') };
      }
      throw e;
    }
  }

  model(params) {
    return hash({
      hasConfig: this.pkiMountHasConfig,
      roles: this.fetchRoles(params),
      parentModel: this.modelFor('roles'),
      pageFilter: params.pageFilter,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const roles = resolvedModel.roles;

    if (roles?.length) controller.notConfiguredMessage = getCliMessage('roles');
    else controller.notConfiguredMessage = getCliMessage();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
