/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { Breadcrumb } from 'vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export type RoleRouteModel = ModelFrom<PkiExternalRolesRoleRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export default class PkiExternalRolesRoleRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { role_name } = this.paramsFor('external.roles.role') as { role_name: string };
    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      role_name,
    };
  }

  setupController(controller: RouteController, resolvedModel: RoleRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'Roles', route: 'external.roles', model: currentPath },
      { label: resolvedModel.role_name },
    ];
  }
}
