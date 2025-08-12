/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import {
  isValidPluginCatalogResponse,
  addVersionsToEngines,
  categorizeEnginesByStatus,
} from 'vault/utils/plugin-catalog-helpers';

export default class MountBackendTypeForm extends Component {
  @service version;
  @service api;

  @tracked pluginCatalogData = null;
  @tracked pluginCatalogError = null;
  @tracked isLoadingPluginCatalog = false;
  @tracked showFlyout = false;
  @tracked flyoutPluginName = null;
  @tracked flyoutPluginType = null;
  @tracked flyoutDisplayName = null;

  constructor(owner, args) {
    super(owner, args);
    if (args.mountCategory === 'secret') {
      this.fetchPluginCatalog();
    }
  }

  async fetchPluginCatalog() {
    this.isLoadingPluginCatalog = true;
    this.pluginCatalogError = null;

    try {
      const response = await this.api.getPluginCatalog();

      if (isValidPluginCatalogResponse(response)) {
        this.pluginCatalogData = response.data;
      } else {
        this.pluginCatalogError = 'Invalid plugin catalog response';
      }
    } catch (error) {
      this.pluginCatalogError = error.message || 'Failed to load plugin catalog';
    } finally {
      this.isLoadingPluginCatalog = false;
    }
  }

  get secretEngines() {
    if (!this.pluginCatalogData) {
      return filterEnginesByMountCategory('secret');
    }

    const baseEngines = filterEnginesByMountCategory('secret');
    return addVersionsToEngines(baseEngines, this.pluginCatalogData.detailed);
  }

  get categorizedSecretEngines() {
    const engines = this.secretEngines || [];
    const categorized = categorizeEnginesByStatus(engines);
    const categories = ['generic', 'cloud', 'infra'];

    return categories.map((categoryName) => {
      const enabledInCategory = categorized.enabled.filter(
        (engine) => engine.pluginCategory === categoryName
      );
      const disabledInCategory = categorized.disabled.filter(
        (engine) => engine.pluginCategory === categoryName
      );

      return {
        category: categoryName,
        enabledEngines: enabledInCategory,
        disabledEngines: disabledInCategory,
      };
    });
  }

  get authMethods() {
    return filterEnginesByMountCategory('auth');
  }

  get mountTypes() {
    return this.args.mountCategory === 'secret' ? this.secretEngines : this.authMethods;
  }

  @action
  handleDisabledPluginClick(plugin) {
    this.flyoutPluginName = plugin.name;
    this.flyoutPluginType = plugin.type;
    this.flyoutDisplayName = plugin.displayName;
    this.showFlyout = true;
  }

  @action
  closeFlyout() {
    this.showFlyout = false;
    this.flyoutPluginName = null;
    this.flyoutPluginType = null;
    this.flyoutDisplayName = null;
  }
}
