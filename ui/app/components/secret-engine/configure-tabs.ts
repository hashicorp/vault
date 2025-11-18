/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

import type { EngineDisplayData } from 'vault/utils/all-engines-metadata';

/**
 * @module ConfigureTabs
 * These tabs render in the shared general-settings, plugin-settings and edit routes of the secret engines headers.
 *
 * @param {string} [configRoute] - only passed when rendering the vault.cluster.secrets.backend.configuration.edit route to highlight the tab for that view
 * @param {object} engineMetadata - engine specific metadata
 * @param {boolean} [isConfigured] - whether an engine has been configured. if configured, plugin settings exist
 * @param {object} path - the secret engine mount path, sometimes referred to as the engine "id" or "backend"
 */

interface Args {
  configRoute?: string;
  engineMetadata: EngineDisplayData;
  isConfigured: boolean;
  path: string;
}

export default class ConfigureTabs extends Component<Args> {
  routePrefix = 'vault.cluster.secrets.backend.';

  // The plugin settings tab only renders if an engine is configurable.
  // `configEditRoute` and `configReadRoute` are defined for engines with custom route patterns, like Ember engines.
  //  Otherwise navigates to default 'backend.configuration' routes.
  get pluginSettingsRoute() {
    const { engineMetadata, isConfigured } = this.args;

    // If the engine is configurable, but not configured, navigate to its edit route
    if (engineMetadata.isConfigurable && !isConfigured) {
      const route = engineMetadata.configEditRoute || 'configuration.edit';
      return this.routePrefix + route;
    }

    // For configured engines, determine route based on context:
    // - If @configRoute is passed (from edit.hbs) the user has navigated to edit and this ensures the tab is highlighted.
    // - Otherwise, route to the read view (`configReadRoute` or 'plugin-settings')
    const route = this.args.configRoute || engineMetadata?.configReadRoute || 'configuration.plugin-settings';
    return this.routePrefix + route;
  }
}
