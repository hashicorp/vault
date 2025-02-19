/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

@withConfig()
export default class PkiKeysIndexRoute extends Route {
  @service pagination;
  @service secretMountPath;
  @service store; // used by @withConfig decorator

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  model(params) {
    const page = Number(params.page) || 1;
    return hash({
      hasConfig: this.pkiMountHasConfig,
      parentModel: this.modelFor('keys'),
      keyModels: this.pagination
        .lazyPaginatedQuery('pki/key', {
          backend: this.secretMountPath.currentPath,
          responsePath: 'data.keys',
          page,
          skipCache: page === 1,
        })
        .catch((err) => {
          if (err.httpStatus === 404) {
            return [];
          } else {
            throw err;
          }
        }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: resolvedModel.parentModel.id },
      { label: 'Keys', route: 'keys.index', model: resolvedModel.parentModel.id },
    ];
    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
