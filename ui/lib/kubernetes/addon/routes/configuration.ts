/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type { KubernetesApplicationModel } from './application';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type { KubernetesConfigureRequest } from '@hashicorp/vault-client-typescript';
import type ApiService from 'vault/services/api';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesApplicationModel;
}

export type KubernetesConfigureModel = ModelFrom<KubernetesConfigureRoute>;

export default class KubernetesConfigureRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;

  @service('app-router') declare readonly router: RouterService;

  async model() {
    const { secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    let config: KubernetesConfigureRequest | undefined;
    let promptConfig = false;
    let configError: unknown;
    try {
      const { data } = await this.api.secrets.kubernetesReadConfiguration(secretsEngine.id);
      config = data as KubernetesConfigureRequest;
    } catch (error) {
      const { response, status } = await this.api.parseError(error);
      // not considering 404 an error since it triggers the cta
      if (status === 404) {
        promptConfig = true;
      } else {
        // ignore if the user does not have permission or other failures so as to not block the other operations
        // this error is thrown in the configuration route so we can display the error in the view
        configError = response;
      }
    }
    // in case of any error other than 404 we want to display that to the user
    if (configError) {
      throw configError;
    }
    return { secretsEngine, config, promptConfig };
  }

  afterModel(resolvedModel: KubernetesConfigureModel) {
    if (!resolvedModel.config) {
      this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.configure');
    }
  }

  setupController(controller: RouteController, resolvedModel: KubernetesConfigureModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Configuration' },
    ];
  }
}
