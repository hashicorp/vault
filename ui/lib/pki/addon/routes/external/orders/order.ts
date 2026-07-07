/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/vault/route';

import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { Breadcrumb } from 'vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export type OrderRouteModel = ModelFrom<PkiExternalOrdersOrderRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export default class PkiExternalOrdersOrderRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { order_id } = this.paramsFor('external.orders.order') as { order_id: string };

    return {
      engine: this.modelFor('application') as SecretsEngineResource,
      order_id,
    };
  }

  setupController(controller: RouteController, resolvedModel: OrderRouteModel) {
    super.setupController(controller, resolvedModel);
    const { currentPath } = this.secretMountPath;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'external.overview', model: currentPath },
      { label: 'Orders', route: 'external.orders', model: currentPath },
      { label: resolvedModel.order_id },
    ];
  }
}
