/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';
import { SecretsApiPkiExternalCaListRoleActiveOrdersListEnum } from '@hashicorp/vault-client-typescript';

import type { RoleRouteModel } from '../role';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type RoleActiveOrdersRouteModel = ModelFrom<PkiExternalRolesRoleActiveOrdersRoute>;

export default class PkiExternalRolesRoleActiveOrdersRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { engine, role_name } = this.modelFor('external.roles.role') as RoleRouteModel;
    let activeOrders: string[] = [];
    try {
      const resp = await this.api.secrets.pkiExternalCaListRoleActiveOrders(
        role_name,
        this.secretMountPath.currentPath,
        SecretsApiPkiExternalCaListRoleActiveOrdersListEnum.TRUE
      );
      activeOrders = resp.keys ?? [];
    } catch (e) {
      // Catch 404s and render empty state instead; throw all other errors.
      const error = await this.api.parseError(e);
      if (error.status !== 404) {
        throw e;
      }
    }

    return {
      engine,
      activeOrders,
      role_name,
    };
  }
}
