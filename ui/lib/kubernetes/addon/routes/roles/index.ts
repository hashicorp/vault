/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';
import { SecretsApiKubernetesListRolesListEnum } from '@hashicorp/vault-client-typescript';

import type { KubernetesApplicationModel } from '../application';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import type Transition from '@ember/routing/transition';
import type CapabilitiesService from 'vault/services/capabilities';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesRolesModel;
}

export type KubernetesRolesModel = ModelFrom<KubernetesRolesRoute>;

export default class KubernetesRolesRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model(_params: unknown, transition: Transition) {
    const { promptConfig, secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    const model = { promptConfig, secretsEngine, roles: [], capabilities: {} };
    const { currentPath } = this.secretMountPath;

    try {
      // filter roles based on pageFilter value
      const { pageFilter } = (transition.to?.queryParams || {}) as { pageFilter?: string };
      const { keys } = await this.api.secrets.kubernetesListRoles(
        currentPath,
        SecretsApiKubernetesListRolesListEnum.TRUE
      );
      const roles = pageFilter
        ? keys?.filter((key) => key.toLowerCase().includes(pageFilter.toLowerCase()))
        : keys;

      // fetch capabilities for filtered roles
      const paths = roles?.map((role) =>
        this.capabilities.pathFor('kubernetesRole', { backend: currentPath, name: role })
      );
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        ...model,
        roles,
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status !== 404) {
        throw error;
      }
    }

    return model;
  }

  setupController(controller: RouteController, resolvedModel: KubernetesRolesModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Roles' },
    ];
  }
}
