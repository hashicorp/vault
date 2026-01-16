/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { Breadcrumb } from 'vault/app-types';
import { KubernetesApplicationModel } from './application';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  secretsEngine: SecretsEngineResource;
}

type KubernetesErrorModel = ModelFrom<KubernetesErrorRoute>;

export default class KubernetesErrorRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  setupController(controller: RouteController, resolvedModel: KubernetesErrorModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
    ];
    const { secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    controller.secretsEngine = secretsEngine;
  }
}
