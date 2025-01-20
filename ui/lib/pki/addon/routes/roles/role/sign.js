/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

withConfirmLeave();
export default class PkiRoleSignRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const { role } = this.paramsFor('roles/role');
    return this.store.createRecord('pki/certificate/sign', {
      role,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { role } = this.paramsFor('roles/role');
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: role, route: 'roles.role.details', models: [this.secretMountPath.currentPath, role] },
      { label: 'sign certificate' },
    ];
    // This is updated on successful generate in the controller
    controller.hasSubmitted = false;
  }
}
