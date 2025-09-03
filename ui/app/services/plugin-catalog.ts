/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { service } from '@ember/service';
import { sanitizePath } from 'core/utils/sanitize-path';

import type ApiService from 'vault/services/api';
import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';

export interface PluginCatalogPlugin {
  name: string;
  type: string;
  builtin: boolean;
  version: string;
  deprecation_status?: string;
  oci_image?: string;
  runtime?: string;
}

export interface PluginCatalogData {
  detailed: Array<PluginCatalogPlugin>;
}

export interface EnhancedPluginCatalogData extends PluginCatalogData {
  secret: Array<string>;
  auth: Array<string>;
  database: Array<string>;
}

export interface PluginCatalogResponse {
  data: EnhancedPluginCatalogData | null;
  error: boolean;
}

export default class PluginCatalogService extends Service {
  @service declare readonly api: ApiService;
  @service declare readonly auth: AuthService;
  @service declare readonly namespace: NamespaceService;

  /**
   * Fetches the plugin catalog from the Vault API
   * Uses the API service middleware for authentication, namespacing, and error handling
   * @returns Promise resolving to plugin catalog data and error state
   */
  async fetchPluginCatalog(): Promise<PluginCatalogResponse> {
    try {
      const response = await this.api.sys.pluginsCatalogListPlugins({
        headers: {
          token: this.auth.currentToken,
          namespace: sanitizePath(this.namespace.path),
        },
      });

      if (response && response.detailed && Array.isArray(response.detailed) && response.detailed.length > 0) {
        const detailedPlugins = response.detailed as PluginCatalogPlugin[];

        const secretPlugins = detailedPlugins
          .filter((plugin) => plugin.type === 'secret')
          .map((plugin) => plugin.name);

        const authPlugins = detailedPlugins
          .filter((plugin) => plugin.type === 'auth')
          .map((plugin) => plugin.name);

        const databasePlugins = detailedPlugins
          .filter((plugin) => plugin.type === 'database')
          .map((plugin) => plugin.name);

        return {
          data: {
            secret: secretPlugins,
            auth: authPlugins,
            database: databasePlugins,
            detailed: detailedPlugins,
          },
          error: false,
        };
      }

      return {
        data: null,
        error: true,
      };
    } catch (error) {
      return {
        data: null,
        error: true,
      };
    }
  }

  /**
   * Gets plugins of a specific type from the catalog
   * @param type - The plugin type ('secret', 'auth', 'database')
   * @returns Promise resolving to array of plugin names for the specified type
   */
  async getPluginsByType(type: keyof Omit<EnhancedPluginCatalogData, 'detailed'>): Promise<Array<string>> {
    const response = await this.fetchPluginCatalog();
    return response.data?.[type] || [];
  }

  /**
   * Gets detailed plugin information for all plugins
   * @returns Promise resolving to array of detailed plugin objects
   */
  async getDetailedPlugins(): Promise<Array<PluginCatalogPlugin>> {
    const response = await this.fetchPluginCatalog();
    return response.data?.detailed || [];
  }

  /**
   * Gets detailed plugin information filtered by type
   * @param type - The plugin type to filter by ('secret', 'auth', 'database')
   * @returns Promise resolving to array of detailed plugin objects of the specified type
   */
  async getDetailedPluginsByType(type: string): Promise<Array<PluginCatalogPlugin>> {
    const response = await this.fetchPluginCatalog();
    return response.data?.detailed.filter((plugin) => plugin.type === type) || [];
  }

  /**
   * Gets a specific plugin by name and type
   * @param name - The plugin name
   * @param type - The plugin type ('secret', 'auth', 'database')
   * @returns Promise resolving to the plugin object or undefined if not found
   */
  async getPlugin(name: string, type: string): Promise<PluginCatalogPlugin | undefined> {
    const response = await this.fetchPluginCatalog();
    return response.data?.detailed.find((plugin) => plugin.name === name && plugin.type === type);
  }
}
