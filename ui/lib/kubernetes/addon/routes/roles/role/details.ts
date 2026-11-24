/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type Controller from '@ember/controller';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type CapabilitiesService from 'vault/services/capabilities';
import type { Breadcrumb } from 'vault/app-types';
import type { KubernetesRole } from 'vault/vault/secrets/kubernetes';

export type KubernetesRoleDetailsModel = ModelFrom<KubernetesRoleDetailsRoute>;

interface RouteController extends Controller {
  model: KubernetesRoleDetailsModel;
  breadcrumbs: Array<Breadcrumb>;
}

export default class KubernetesRoleDetailsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model() {
    const { currentPath } = this.secretMountPath;
    const { name } = this.paramsFor('roles.role');
    const { data } = await this.api.secrets.kubernetesReadRole(name as string, currentPath);
    // fetch capabilities for role
    const pathMap = {
      role: this.capabilities.pathFor('kubernetesRole', { backend: currentPath, name }),
      creds: this.capabilities.pathFor('kubernetesCreds', { backend: currentPath, name }),
    };
    const capabilities = await this.capabilities.fetch(Object.values(pathMap));

    return {
      role: data as KubernetesRole,
      capabilities: {
        canUpdate: capabilities[pathMap.role]?.canUpdate,
        canDelete: capabilities[pathMap.role]?.canDelete,
        canGenerateCreds: capabilities[pathMap.creds]?.canCreate,
      },
    };
  }

  setupController(controller: RouteController, resolvedModel: KubernetesRoleDetailsModel) {
    super.setupController(controller, resolvedModel);

    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles', model: currentPath },
      { label: resolvedModel.role.name },
    ];
  }
}
