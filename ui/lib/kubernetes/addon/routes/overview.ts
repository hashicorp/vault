/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { ModelFrom } from 'vault/route';

import type { KubernetesApplicationModel } from './application';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Store from '@ember-data/store';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: KubernetesOverviewModel;
}

export type KubernetesOverviewModel = ModelFrom<KubernetesOverviewRoute>;

export default class KubernetesOverviewRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const backend = this.secretMountPath.currentPath;
    const { promptConfig, secretsEngine } = this.modelFor('application') as KubernetesApplicationModel;
    return hash({
      promptConfig,
      secretsEngine,
      roles: this.store.query('kubernetes/role', { backend }).catch(() => []),
    });
  }

  setupController(controller: RouteController, resolvedModel: KubernetesOverviewModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath },
    ];
  }
}
