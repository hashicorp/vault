/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { hash } from 'rsvp';

@withConfirmLeave('model.role', ['model.issuers'])
export default class PkiRolesCreateRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return hash({
      role: this.store.createRecord('pki/role', { backend }),
      issuers: this.store.query('pki/issuer', { backend }).catch((err) => {
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
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: 'create' },
    ];
  }
}
