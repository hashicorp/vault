/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';
import timestamp from 'core/utils/timestamp';

import type { Breadcrumb } from 'vault/app-types';
import type { ExternalRouteModel } from 'pki/routes/external';
import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export type RolesIndexRouteModel = ModelFrom<PkiExternalRolesIndexRoute>;

export default class PkiExternalRolesIndexRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { engine, rolesResp } = this.modelFor('external') as ExternalRouteModel;
    if (rolesResp.error.message) {
      throw rolesResp.error;
    }
    return {
      engine,
      roles: rolesResp.keys,
      responseTimestamp: timestamp.now(),
    };
  }

  setupController(controller: RouteController, resolvedModel: RolesIndexRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'Roles' },
    ];
  }
}
