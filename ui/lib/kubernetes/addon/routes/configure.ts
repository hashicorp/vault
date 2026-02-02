/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import KubernetesConfigForm from 'vault/forms/secrets/kubernetes/config';
import { ModelFrom } from 'vault/route';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type { KubernetesApplicationModel } from './application';
import type { Breadcrumb } from 'vault/app-types';
import type Controller from '@ember/controller';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesConfigureModel;
}
export type KubernetesConfigureModel = ModelFrom<KubernetesConfigureRoute>;

export default class KubernetesConfigureRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { config, secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    const data = config || { disable_local_ca_jwt: false };
    return { form: new KubernetesConfigForm(data, { isNew: !config }), secretsEngine };
  }

  setupController(controller: RouteController, resolvedModel: KubernetesConfigureModel) {
    super.setupController(controller, resolvedModel);

    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Configure' },
    ];
  }
}
