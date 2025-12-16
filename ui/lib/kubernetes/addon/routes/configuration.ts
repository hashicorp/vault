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

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesApplicationModel;
}

export type KubernetesConfigureModel = ModelFrom<KubernetesConfigureRoute>;

export default class KubernetesConfigureRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { config, configError, secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    // in case of any error other than 404 we want to display that to the user
    if (configError) {
      throw configError;
    }
    return { secretsEngine, config };
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
