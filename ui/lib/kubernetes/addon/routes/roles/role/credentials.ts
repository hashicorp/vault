/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';

export type KubernetesRoleCredentialsModel = ModelFrom<KubernetesRoleCredentialsRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesRoleCredentialsModel;
}

export default class KubernetesRoleCredentialsRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { name } = this.paramsFor('roles.role');
    const roleName = name as string;
    return { roleName };
  }

  setupController(controller: RouteController, resolvedModel: KubernetesRoleCredentialsModel) {
    super.setupController(controller, resolvedModel);

    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles', model: currentPath },
      { label: resolvedModel.roleName, route: 'roles.role.details' },
      { label: 'Credentials' },
    ];
  }
}
