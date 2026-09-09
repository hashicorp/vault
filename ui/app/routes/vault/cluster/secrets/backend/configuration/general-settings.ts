/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { getPluginVersionsFromEngineType } from 'vault/utils/plugin-catalog-helpers';
import SecretsEngineResource from 'vault/resources/secrets/engine';

import type Controller from '@ember/controller';
import type RouterService from '@ember/routing/router-service';
import type Transition from '@ember/routing/transition';
import type SecretEngineModel from 'vault/models/secret-engine';
import type ApiService from 'vault/services/api';
import type PluginCatalogService from 'vault/services/plugin-catalog';
import type UnsavedChangesService from 'vault/services/unsaved-changes';

interface RouteModel {
  secretsEngine: SecretEngineModel;
  versions: string[];
  config: Record<string, unknown>;
  pinnedVersion: string | null;
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

  async model() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend') as Record<string, unknown>;
    const response = await this.api.sys.internalUiReadMountInformation(backend as string);
    const secretsEngine = new SecretsEngineResource({
      ...(response as SecretsEngineResource),
      path: `${backend}/`,
    });
    const { data } = await this.pluginCatalog.getRawPluginCatalogData();
    const versions = getPluginVersionsFromEngineType(data?.secret, secretsEngine.type);

    // Fetch version data (pinned, current, and running versions)
    const pluginName = secretsEngine.type;
    const mountPath = secretsEngine.id;

    // Fetch both pinned version and mount info in parallel
    const [pinnedVersion, mountInfo] = await Promise.all([
      this.api.sys.pluginsCatalogPinsReadPinnedVersion(pluginName, 'secret').catch(() => {
        // Silently handle errors - pins are optional
        return null;
      }),
      this.api.sys.internalUiReadMountInformation(mountPath),
    ]);

    const pinnedVersionString = pinnedVersion?.version || null;

    // Update the secretsEngine properties directly
    secretsEngine.plugin_version = mountInfo?.plugin_version || '';
    secretsEngine.running_plugin_version = mountInfo?.running_plugin_version || '';

    const model = { secretsEngine, versions, pinnedVersion: pinnedVersionString };
    this.unsavedChanges.initialState = JSON.parse(JSON.stringify(model.secretsEngine));

    return { secretsEngine, versions, pinnedVersion: pinnedVersionString };
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);
    const { secretsEngine } = resolvedModel;

    const breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster', icon: 'vault' },
      { label: 'Secrets engines', route: 'vault.cluster.secrets' },
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

    const state = model ? (model['secretsEngine'] as Record<string, unknown> | undefined) : {};
    this.unsavedChanges.setup(state);
    // Only intercept transition if leaving THIS route and there are changes
    const targetRoute = transition?.to?.name ?? '';

    if (this.routeName !== targetRoute && this.unsavedChanges.hasChanges) {
      transition.abort();
      this.unsavedChanges.show(transition);
    }
    return true;
  }
}
