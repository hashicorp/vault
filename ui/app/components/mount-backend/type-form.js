/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import { isValidPluginCatalogResponse } from 'vault/utils/plugin-catalog-helpers';

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

  constructor(owner, args) {
    super(owner, args);
    // Only fetch plugin catalog for secret engines in Phase 1
    if (args.mountCategory === 'secret') {
      this.loadPluginCatalog();
    }
  }

  async loadPluginCatalog() {
    try {
      const response = await this.api.getPluginCatalog('secret');

      if (isValidPluginCatalogResponse(response)) {
        this.pluginCatalogData = response.data.detailed;
        console.log('Plugin catalog loaded successfully:', {
          pluginCount: this.pluginCatalogData.length,
          plugins: this.pluginCatalogData.map((p) => ({
            name: p.name,
            version: p.version,
            builtin: p.builtin,
          })),
        });
      } else {
        this.pluginCatalogError = new Error('Invalid response structure');
      }
    } catch (error) {
      this.pluginCatalogError = error;
    }
  }

  get secretEngines() {
    // If an enterprise license is present, return all secret engines;
    // otherwise, return only the secret engines supported in OSS.
    return filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: this.version.isEnterprise });
  }

  get authMethods() {
    // If an enterprise license is present, return all auth methods;
    // otherwise, return only the auth methods supported in OSS.
    return filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: this.version.isEnterprise });
  }

  get mountTypes() {
    return this.args.mountCategory === 'secret' ? this.secretEngines : this.authMethods;
  }
}
