/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiCertificateForm from 'vault/forms/secrets/pki/certificate';

export default class PkiRoleSignRoute extends Route {
  @service secretMountPath;

  model() {
    const { role } = this.paramsFor('roles/role');
    return {
      role,
      form: new PkiCertificateForm('PkiSignWithRoleRequest', {}, { isNew: true }),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { role } = this.paramsFor('roles/role');
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: role, route: 'roles.role.details', models: [this.secretMountPath.currentPath, role] },
      { label: 'Sign Certificate' },
    ];
    // This is updated on successful generate in the controller
    controller.hasSubmitted = false;
  }
}
