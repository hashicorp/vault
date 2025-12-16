/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';
import KubernetesRoleForm from 'vault/forms/secrets/kubernetes/role';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import type { KubernetesRole } from 'vault/secrets/kubernetes';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';

export type KubernetesRoleEditModel = ModelFrom<KubernetesRoleEditRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesRoleEditModel;
}

export default class KubernetesRoleEditRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;

  async model() {
    const { currentPath } = this.secretMountPath;
    const { name } = this.paramsFor('roles.role');
    const { data } = await this.api.secrets.kubernetesReadRole(name as string, currentPath);
    return new KubernetesRoleForm(data as KubernetesRole, { isNew: false });
  }

  setupController(controller: RouteController, resolvedModel: KubernetesRoleEditModel) {
    super.setupController(controller, resolvedModel);

    const { name } = this.paramsFor('roles.role');
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles', model: currentPath },
      { label: name as string, route: 'roles.role' },
      { label: 'Edit' },
    ];
  }
}
