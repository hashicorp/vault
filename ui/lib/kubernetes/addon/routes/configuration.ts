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

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesApplicationModel;
}

export type KubernetesConfigureModel = ModelFrom<KubernetesConfigureRoute>;

export default class KubernetesConfigureRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;
  @service('app-router') declare readonly router: RouterService;

  model() {
    const { config, configError, secretsEngine, promptConfig } = this.modelFor(
      'application'
    ) as KubernetesApplicationModel;
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
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Configuration' },
    ];
  }
}
