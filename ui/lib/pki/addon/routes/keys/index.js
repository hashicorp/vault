/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

@withConfig()
export default class PkiKeysIndexRoute extends Route {
  @service secretMountPath;
  @service store;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  model(params) {
    return hash({
      hasConfig: this.shouldPromptConfig,
      parentModel: this.modelFor('keys'),
      keyModels: this.store
        .lazyPaginatedQuery('pki/key', {
          backend: this.secretMountPath.currentPath,
          responsePath: 'data.keys',
          page: Number(params.page) || 1,
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
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'keys', route: 'keys.index' },
    ];
    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
