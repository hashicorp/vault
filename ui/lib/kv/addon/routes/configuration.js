/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KvConfigurationRoute extends Route {
  @service api;

  async model() {
    const backend = this.modelFor('application');
    const {
      type,
      path,
      accessor,
      running_plugin_version,
      local,
      seal_wrap,
      config: { default_lease_ttl, max_lease_ttl },
      options: { version },
    } = await this.api.sys.internalUiReadMountInformation(backend.id);
    // display mount config if engine config request fails
    const engineConfig = await this.api.secrets.kvV2ReadConfiguration(backend.id).catch(() => {});

    return {
      ...engineConfig,
      type,
      path,
      accessor,
      running_plugin_version,
      local,
      seal_wrap,
      default_lease_ttl,
      max_lease_ttl,
      version,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { id } = this.modelFor('application');
    controller.backend = id;
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: id, route: 'list', model: id },
      { label: 'Configuration' },
    ];
  }
}
