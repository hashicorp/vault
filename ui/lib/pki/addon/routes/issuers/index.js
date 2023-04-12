/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';

export default class PkiIssuersListRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store
      .query('pki/issuer', { backend: this.secretMountPath.currentPath })
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
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
    ];
    controller.message = PKI_DEFAULT_EMPTY_STATE_MSG;
  }
}
