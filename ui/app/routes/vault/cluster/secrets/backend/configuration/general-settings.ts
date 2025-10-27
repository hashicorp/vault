/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getPluginVersionsFromEngineType } from 'vault/utils/plugin-catalog-helpers';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type PluginCatalogService from 'vault/services/plugin-catalog';
import type UnsavedChangesService from 'vault/services/unsaved-changes';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type RouterService from '@ember/routing/router-service';
import type SecretEngineModel from 'vault/models/secret-engine';

interface RouteModel {
  secretsEngine: SecretEngineModel;
  versions: string[];
}

interface RouteController extends Controller {
  model: Record<string, unknown> | undefined;
  changedFields: string[];
}

export default class SecretsBackendConfigurationGeneralSettingsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly pluginCatalog: PluginCatalogService;
  @service declare readonly unsavedChanges: UnsavedChangesService;

  oldModel: Record<string, unknown> | undefined;

  async model() {
    const secretsEngine = this.modelFor('vault.cluster.secrets.backend') as SecretsEngineResource;
    const { data } = await this.pluginCatalog.getRawPluginCatalogData();
    const versions = getPluginVersionsFromEngineType(data?.secret, secretsEngine.type);

    const model = { secretsEngine, versions };
    this.unsavedChanges.initialState = JSON.parse(JSON.stringify(model.secretsEngine));

    return { secretsEngine, versions };
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

  @action
  willTransition(transition: Transition) {
    // eslint-disable-next-line ember/no-controller-access-in-routes
    const controller = this.controllerFor(this.routeName) as RouteController;

    const { model } = controller;

    // Only intercept transition if leaving THIS route and there are changes
    const targetRoute = transition?.to?.name ?? '';
    // Saving transitions to the index route, so we do not one to intercept the transition then
    if (this.routeName !== targetRoute && this.oldModel && model) {
      const oldModel = this.oldModel['secretsEngine'] as Record<string, unknown> | undefined;
      const currentModel = model['secretsEngine'] as Record<string, unknown> | undefined;
      this.unsavedChanges.setupProperties(oldModel, currentModel);
      this.unsavedChanges.getDiff();

      if (this.unsavedChanges.hasChanges) {
        transition.abort();
        this.unsavedChanges.showModal = true;
      }
    }
    return true;
  }
}
