/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';
import timestamp from 'core/utils/timestamp';

import type { Breadcrumb } from 'vault/app-types';
import type { RoleRouteModel } from '../role';
import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export type RoleOrderRouteModel = ModelFrom<PkiExternalRolesRoleOrderRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export default class PkiExternalRolesRoleOrderRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  async model({ order_id }: { order_id: string }) {
    const { role_name } = this.modelFor('external.roles.role') as RoleRouteModel;

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      order_id,
      role_name,
      responseTimestamp: timestamp.now(),
    };
  }
  setupController(controller: RouteController, resolvedModel: RoleOrderRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'Roles', route: 'external.roles', model: currentPath },
      {
        label: resolvedModel.role_name,
        route: 'external.roles.role',
        models: [currentPath, resolvedModel.role_name],
      },
      { label: resolvedModel.order_id },
    ];
  }
}
