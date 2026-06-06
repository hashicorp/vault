/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';
import { SecretsApiKubernetesListRolesListEnum } from '@hashicorp/vault-client-typescript';

import type { KubernetesApplicationModel } from './application';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesOverviewModel;
}

export type KubernetesOverviewModel = ModelFrom<KubernetesOverviewRoute>;

export default class KubernetesOverviewRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { currentPath } = this.secretMountPath;
    const { promptConfig, secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;

    const { keys } = await this.api.secrets
      .kubernetesListRoles(currentPath, SecretsApiKubernetesListRolesListEnum.TRUE)
      .catch(() => ({ keys: [] }));

    return {
      promptConfig,
      secretsEngine,
      roles: keys,
    };
  }

  setupController(controller: RouteController, resolvedModel: KubernetesOverviewModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath },
    ];
  }
}
