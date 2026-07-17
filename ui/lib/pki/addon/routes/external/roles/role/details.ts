/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type { RoleRouteModel } from '../role';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type RoleDetailsRouteModel = ModelFrom<PkiExternalRolesRoleDetailsRoute>;

export default class PkiExternalRolesRoleDetailsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { role_name } = this.modelFor('external.roles.role') as RoleRouteModel;
    const role = await this.api.secrets.pkiExternalCaReadRole(role_name, this.secretMountPath.currentPath);
    return { role };
  }
}
