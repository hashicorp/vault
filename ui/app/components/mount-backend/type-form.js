/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { filterEnginesByMountCategory } from 'core/utils/all-engines-metadata';
import keys from 'core/utils/keys';
import {
  categorizeEnginesByStatus,
  enhanceEnginesWithCatalogData,
  MOUNT_CATEGORIES,
  PLUGIN_CATEGORIES,
  PLUGIN_TYPES,
} from 'vault/utils/plugin-catalog-helpers';

/**
 *
 * @module MountBackendTypeForm
 * MountBackendTypeForm components are used to display type options for
 * mounting either an auth method or secret engine.
 *
 * @example
 * ```js
 * <MountBackend::TypeForm @setMountType={{this.setMountType}} @mountCategory="secret" @pluginCatalogData={{this.pluginCatalogData}} @pluginCatalogError={{this.pluginCatalogError}} />
 * ```
 * @param {CallableFunction} setMountType - function will receive the mount type string. Should update the model type value
 * @param {string} [mountCategory=auth] - mount category can be `auth` or `secret`
 * @param {object} [pluginCatalogData] - plugin catalog data fetched from the route model
 * @param {boolean} [pluginCatalogError] - indicates if there was an error fetching plugin catalog data
 */

export default class MountBackendTypeForm extends Component {
  @service version;

  @tracked showFlyout = false;
  @tracked flyoutPlugin = null;
  @tracked flyoutPluginType = null;

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
        this.args.pluginCatalogData?.detailed?.filter((plugin) => plugin?.type === PLUGIN_TYPES.SECRET) || [];
      const databasePluginsDetailed =
        this.args.pluginCatalogData?.detailed?.filter((plugin) => plugin?.type === PLUGIN_TYPES.DATABASE) ||
        [];

      return enhanceEnginesWithCatalogData(staticEngines, secretEnginesDetailed, databasePluginsDetailed);
    }

    return staticEngines;
  }

  get authMethods() {
    // If an enterprise license is present, return all auth methods;
    // otherwise, return only the auth methods supported in OSS.
    return filterEnginesByMountCategory({
      mountCategory: MOUNT_CATEGORIES.AUTH,
      isEnterprise: !!this.version?.isEnterprise,
    });
  }

  get mountTypes() {
    if (this.args.mountCategory === MOUNT_CATEGORIES.SECRET) {
      return this.secretEngines || [];
    }
    return this.authMethods || [];
  }

  get pluginCategoriesList() {
    return [
      PLUGIN_CATEGORIES.GENERIC,
      PLUGIN_CATEGORIES.CLOUD,
      PLUGIN_CATEGORIES.INFRA,
      // PLUGIN_CATEGORIES.EXTERNAL, // TODO: enable external plugins once version selection is available (VAULT-39241)
    ];
  }

  get secretMountCategory() {
    return MOUNT_CATEGORIES.SECRET;
  }

  @action
  getMountTypesByCategory(category) {
    try {
      if (!this.args) {
        return { enabled: [], disabled: [] };
      }

      const mountTypes = this.mountTypes;
      if (!mountTypes || !Array.isArray(mountTypes)) {
        return { enabled: [], disabled: [] };
      }

      const allTypes = mountTypes.filter((type) => type?.pluginCategory === category);
      return categorizeEnginesByStatus(allTypes);
    } catch (error) {
      return { enabled: [], disabled: [] };
    }
  }

  @action
  handleDisabledPluginKeyDown(plugin, event) {
    // Only handle Enter and Space keys for accessibility
    if (event.key === keys.ENTER || event.key === keys.SPACE) {
      event.preventDefault();
      this.handleDisabledPluginClick(plugin);
    }
  }

  @action
  handleDisabledPluginClick(plugin) {
    this.showFlyout = true;
    this.flyoutPlugin = plugin;
    this.flyoutPluginType = this.args?.mountCategory;
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
