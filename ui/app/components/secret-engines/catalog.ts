/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import {
  enhanceEnginesWithCatalogData,
  categorizeEnginesByStatus,
  MOUNT_CATEGORIES,
  PLUGIN_TYPES,
  PLUGIN_CATEGORIES,
} from 'vault/utils/plugin-catalog-helpers';
import type { PluginCatalogData } from 'vault/services/plugin-catalog';

import type VersionService from 'vault/services/version';

/**
 * @module SecretEnginesCatalog
 * SecretEnginesCatalog component displays available secret engines in a catalog view
 * for selection when mounting a new secret engine.
 *
 * @example
 * ```js
 * <SecretEngines::Catalog @setMountType={{this.setMountType}} @pluginCatalogData={{this.pluginCatalogData}} @pluginCatalogError={{this.pluginCatalogError}} />
 * ```
 */

interface Args {
  setMountType: (type: string) => void;
  pluginCatalogData?: PluginCatalogData;
  pluginCatalogError?: boolean;
}

export default class SecretEnginesCatalogComponent extends Component<Args> {
  @service declare version: VersionService;

  @tracked showFlyout = false;
  @tracked flyoutPlugin: unknown = null;
  @tracked flyoutPluginType: string | null = null;

  get breadcrumbs() {
    return [
      {
        label: 'Secrets engines',
        route: 'vault.cluster.secrets.backends',
      },
      {
        label: 'Enable secrets engine',
      },
    ];
  }

  get secretEngines() {
    // If an enterprise license is present, return all secret engines;
    // otherwise, return only the secret engines supported in OSS.
    const staticEngines = filterEnginesByMountCategory({
      mountCategory: MOUNT_CATEGORIES.SECRET,
      isEnterprise: !!this.version?.isEnterprise,
    });

    // If we have plugin catalog data, merge it with static engines to add catalog info
    if (this.args.pluginCatalogData) {
      const secretEnginesDetailed =
        this.args.pluginCatalogData.detailed?.filter((plugin) => plugin?.type === PLUGIN_TYPES.SECRET) || [];
      const databasePluginsDetailed =
        this.args.pluginCatalogData.detailed?.filter((plugin) => plugin?.type === PLUGIN_TYPES.DATABASE) ||
        [];

      return enhanceEnginesWithCatalogData(staticEngines, secretEnginesDetailed, databasePluginsDetailed);
    }

    return staticEngines;
  }

  get pluginCategoriesList() {
    return [
      PLUGIN_CATEGORIES.GENERIC,
      PLUGIN_CATEGORIES.CLOUD,
      PLUGIN_CATEGORIES.INFRA,

      // TODO: enable external plugins once version selection is available (VAULT-39241)
      // PLUGIN_CATEGORIES.EXTERNAL,
    ];
  }

  get secretMountCategory() {
    return MOUNT_CATEGORIES.SECRET;
  }

  @action
  getMountTypesByCategory(category: string) {
    try {
      const mountTypes = this.secretEngines;
      if (!mountTypes || !Array.isArray(mountTypes)) {
        return { enabled: [], disabled: [] };
      }

      const allTypes = mountTypes.filter((type: unknown) => {
        const engineType = type as { pluginCategory?: string };
        return engineType?.pluginCategory === category;
      });
      return categorizeEnginesByStatus(allTypes);
    } catch (error) {
      return { enabled: [], disabled: [] };
    }
  }

  @action
  handleDisabledPluginClick(plugin: unknown) {
    this.showFlyout = true;
    this.flyoutPlugin = plugin;
    this.flyoutPluginType = 'secret';
  }

  @action
  handleDisabledPluginKeyDown(plugin: unknown, event: KeyboardEvent) {
    // Only handle Enter and Space keys for accessibility
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      this.handleDisabledPluginClick(plugin);
    }
  }

  @action
  openExternalPluginsHelp() {
    this.showFlyout = true;
    this.flyoutPlugin = null;
    this.flyoutPluginType = 'secret';
  }

  @action
  closeFlyout() {
    this.showFlyout = false;
    this.flyoutPlugin = null;
    this.flyoutPluginType = null;
  }
}
