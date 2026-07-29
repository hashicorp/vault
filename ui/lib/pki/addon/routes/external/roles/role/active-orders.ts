/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type { RoleRouteModel } from '../role';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type RoleActiveOrdersRouteModel = ModelFrom<PkiExternalRolesRoleActiveOrdersRoute>;

export default class PkiExternalRolesRoleActiveOrdersRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const parentModel = this.modelFor('external.roles.role') as RoleRouteModel;
    const { activeOrders } = parentModel;

    if (activeOrders.error && activeOrders.error.status !== 404) {
      throw activeOrders.error;
    }

    return {
      ...parentModel,
      activeOrders: activeOrders.list,
    };
  }
}
