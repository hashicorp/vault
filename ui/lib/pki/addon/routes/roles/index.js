/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { getCliMessage } from 'pki/routes/overview';
import { SecretsApiPkiListRolesListEnum } from '@hashicorp/vault-client-typescript';
import { paginate } from 'core/utils/paginate-list';

@withConfig()
export default class PkiRolesIndexRoute extends Route {
  @service api;
  @service secretMountPath;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async model(params) {
    const model = {
      hasConfig: this.pkiMountHasConfig,
      parentModel: this.modelFor('roles'),
      pageFilter: params.pageFilter,
      roles: [],
    };

    try {
      const page = Number(params.page) || 1;
      const { keys: roles } = await this.api.secrets.pkiListRoles(
        this.secretMountPath.currentPath,
        SecretsApiPkiListRolesListEnum.TRUE
      );
      model.roles = paginate(roles, { page });
    } catch (e) {
      if (e.response.status !== 404) {
        throw e;
      }
    }

    return model;
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.notConfiguredMessage = resolvedModel.roles?.length ? getCliMessage('roles') : getCliMessage();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
