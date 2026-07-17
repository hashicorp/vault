/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type { RoleRouteModel } from '../role';
import type SecretMountPath from 'vault/services/secret-mount-path';

export default class PkiExternalRolesRoleIndexRoute extends Route {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;

  redirect(model: RoleRouteModel) {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.external.roles.role.details',
      this.secretMountPath.currentPath,
      model.role_name
    );
  }
}
