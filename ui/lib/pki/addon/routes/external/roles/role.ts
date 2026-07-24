/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type { Breadcrumb } from 'vault/app-types';
import type { ExternalRouteModel } from 'pki/routes/external';
import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type RoleRouteModel = ModelFrom<PkiExternalRolesRoleRoute>;

interface RouteController extends Controller {
  generateCrumbs: (isOrderRoute: boolean, path: string, roleName: string) => Array<Breadcrumb>;
}

export default class PkiExternalRolesRoleRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  async model({ role_name }: { role_name: string }) {
    const { engine } = this.modelFor('external') as ExternalRouteModel;
    return {
      engine,
      role_name,
    };
  }

  // Breadcrumbs and tabs are rendered in the parent template so they remain visible even when a
  // child route throws an error. The "roles.role.order" child route needs to render different breadcrumbs
  // so the template calls this function with the current route state (via matches-current-url).
  // This also avoids us having to duplicate HeaderTabs in each child or using controllerFor to reach across routes.
  generateCrumbs = (isOrderRoute: boolean, roleName: string): Array<Breadcrumb> => {
    const path = this.secretMountPath.currentPath;
    const base: Array<Breadcrumb> = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: path, route: 'external.overview', model: path },
      { label: 'Roles', route: 'external.roles', model: path },
    ];
    if (isOrderRoute) {
      return [
        ...base,
        { label: roleName, route: 'external.roles.role', models: [path, roleName] },
        { label: 'View order' },
      ];
    }
    return [...base, { label: roleName }];
  };

  setupController(controller: RouteController, resolvedModel: RoleRouteModel) {
    super.setupController(controller, resolvedModel);
    controller.generateCrumbs = this.generateCrumbs;
  }
}
