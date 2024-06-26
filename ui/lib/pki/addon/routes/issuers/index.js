/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiIssuersListRoute extends Route {
  @service store;
  @service secretMountPath;

  model(params) {
    const page = Number(params.page) || 1;
    return this.store
      .lazyPaginatedQuery('pki/issuer', {
        backend: this.secretMountPath.currentPath,
        responsePath: 'data.keys',
        page,
        skipCache: page === 1,
        isListView: true,
      })
      .then((issuersModel) => {
        return { issuersModel, parentModel: this.modelFor('issuers') };
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return { parentModel: this.modelFor('issuers') };
        } else {
          throw err;
        }
      });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
    ];
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
