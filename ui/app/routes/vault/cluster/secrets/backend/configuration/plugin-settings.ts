/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Transition from '@ember/routing/transition';

interface RouteModel {
  secretsEngine: SecretEngineModel;
  versions: string[];
  config: Record<string, unknown>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

export default class SecretsBackendConfigurationPluginSettingsRoute extends Route {
  @service declare readonly api: ApiService;

  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    const { config } = this.modelFor('vault.cluster.secrets.backend.configuration') as Record<
      string,
      unknown
    >;

    return { secretsEngine, config };
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);
    const { secretsEngine } = resolvedModel;

    const breadcrumbs = [
      { label: 'Secrets', route: 'vault.cluster.secrets' },
      { label: secretsEngine.id, route: 'vault.cluster.secrets.backend.list-root', model: secretsEngine.id },
      { label: 'Configuration' },
    ];

    controller.set('breadcrumbs', breadcrumbs);
  }
}
