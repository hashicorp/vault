/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
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
      issuers: this.store.query('pki/issuer', { backend }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'roles', route: 'roles.index' },
      { label: 'create' },
    ];
  }
}
