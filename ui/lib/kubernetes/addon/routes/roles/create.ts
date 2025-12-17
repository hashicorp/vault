/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';
import KubernetesRoleForm from 'vault/forms/secrets/kubernetes/role';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';

export type KubernetesRolesCreateModel = ModelFrom<KubernetesRolesCreateRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesRolesCreateModel;
}

export default class KubernetesRolesCreateRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    return new KubernetesRoleForm({}, { isNew: true });
  }

  setupController(controller: RouteController, resolvedModel: KubernetesRolesCreateModel) {
    super.setupController(controller, resolvedModel);

    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles', model: currentPath },
      { label: 'Create' },
    ];
  }
}
