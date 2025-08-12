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

/**
 *
 * @module MountBackendTypeForm
 * MountBackendTypeForm components are used to display type options for
 * mounting either an auth method or secret engine.
 *
 * @example
 * ```js
 * <MountBackend::TypeForm @setMountType={{this.setMountType}} @mountCategory="secret" />
 * ```
 * @param {CallableFunction} setMountType - function will receive the mount type string. Should update the model type value
 * @param {string} [mountCategory=auth] - mount category can be `auth` or `secret`
 */

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
    // Only fetch plugin catalog for secret engines in Phase 1
    if (args.mountCategory === 'secret') {
      this.loadPluginCatalog();
    }
  }

  async loadPluginCatalog() {
    this.isLoadingPluginCatalog = true;
    try {
      const response = await this.api.getPluginCatalog();

      if (isValidPluginCatalogResponse(response)) {
        this.pluginCatalogData = response.data.detailed;
      } else {
        this.pluginCatalogError = new Error('Invalid response structure');
      }
    } catch (error) {
      this.pluginCatalogError = error;
    } finally {
      this.isLoadingPluginCatalog = false;
    }
  }

  get secretEngines() {
    // If an enterprise license is present, return all secret engines;
    // otherwise, return only the secret engines supported in OSS.
    const staticEngines = filterEnginesByMountCategory({
      mountCategory: 'secret',
      isEnterprise: this.version.isEnterprise,
    });

    // If we have plugin catalog data, merge it with static engines to add version info
    if (this.pluginCatalogData) {
      return addVersionsToEngines(staticEngines, this.pluginCatalogData);
    }

    return staticEngines;
  }

  get authMethods() {
    // If an enterprise license is present, return all auth methods;
    // otherwise, return only the auth methods supported in OSS.
    return filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: this.version.isEnterprise });
  }

  get mountTypes() {
    return this.args.mountCategory === 'secret' ? this.secretEngines : this.authMethods;
  }

  get genericMountTypes() {
    const allTypes = this.mountTypes.filter((type) => type.pluginCategory === 'generic');
    return categorizeEnginesByStatus(allTypes);
  }

  get cloudMountTypes() {
    const allTypes = this.mountTypes.filter((type) => type.pluginCategory === 'cloud');
    return categorizeEnginesByStatus(allTypes);
  }

  get infraMountTypes() {
    const allTypes = this.mountTypes.filter((type) => type.pluginCategory === 'infra');
    return categorizeEnginesByStatus(allTypes);
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
