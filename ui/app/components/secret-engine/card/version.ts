/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  model: {
    secretsEngine: SecretsEngineResource;
    versions: string[];
    pinnedVersion: string | null;
  };
  onReload?: (data: {
    pinnedVersion: string | null;
    pluginVersion: string | null;
    runningVersion: string | null;
  }) => void;
}

/**
 * @module SecretEngine::Card::Version
 * SecretEngine::Card::Version component displays the version information for a secrets engine
 * and provides a dropdown to update the plugin version.
 *
 * @example
 * ```js
 * <SecretEngine::Card::Version @model={{this.model}} />
 * ```
 * @param {object} model - Object containing secretsEngine, versions, and version data
 * @param {SecretsEngineResource} model.secretsEngine - The secrets engine resource
 * @param {string[]} model.versions - Array of available plugin versions
 * @param {string|null} model.pinnedVersion - Currently pinned plugin version
 */
export default class SecretEngineCardVersionComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked selectedPluginVersion: string | null = null;

  /**
   * The pinned version for the secret engine type (ie. 'kv', 'keymgmt', 'vault-plugin-secrets-keymgmt'), if a pinned version is set.
   *
   * This value may be null if no pinned version is set for the current engine type.
   */
  get pinnedVersion(): string | null {
    return this.args.model.pinnedVersion;
  }

  /**
   * The running plugin version for the current secrets engine mount.
   *
   * This is the version that is currently active and running for the mount.
   * This value should never be null, it is set by the API when the plugin is mounted.
   */
  get runningVersion(): string | null {
    return this.args.model.secretsEngine.running_plugin_version;
  }

  /**
   * The configured plugin version for the current secrets engine mount.
   *
   * This is the version that is set in the configuration for the mount.
   * This value should never be null, it is set by the API when the plugin is mounted.
   */
  get pluginVersion(): string | null {
    return this.args.model.secretsEngine.plugin_version;
  }

  /**
   * Determines if all version data is available and non-null to evaluate version states.
   */
  get areAllVersionsSet(): boolean {
    return !!(this.pinnedVersion && this.runningVersion && this.pluginVersion);
  }

  /**
   * Determines if there's a manual override of the pinned version.
   *
   * This occurs when the running version matches the configured plugin version, but differs from the pinned version.
   * This indicates that the user has manually overridden the pinned version.
   */
  get hasManualPinnedOverride(): boolean {
    return !!(
      this.areAllVersionsSet &&
      this.runningVersion !== this.pinnedVersion &&
      this.runningVersion === this.pluginVersion
    );
  }

  /**
   * Determines if the configured plugin version is overridden by the pinned version.
   *
   * This occurs when the running version matches the pinned version, but differs from the configured plugin version.
   * This indicates that the pinned version is taking precedence over the configured version.
   */
  get configuredVersionIsOverridden(): boolean {
    return !!(
      this.areAllVersionsSet &&
      this.runningVersion === this.pinnedVersion &&
      this.runningVersion !== this.pluginVersion
    );
  }

  /**
   * Determines if there's a version mismatch between plugin and running versions.
   *
   * When all versions are different, a reload is required to resolve the mismatch.
   */
  get hasVersionMismatch(): boolean {
    return !!(
      this.areAllVersionsSet &&
      this.runningVersion !== this.pinnedVersion &&
      this.runningVersion !== this.pluginVersion
    );
  }

  /**
   * Returns the info message to display, or null if no info message should be shown.
   */
  get infoMessage(): string | null {
    if (this.hasManualPinnedOverride) {
      return `This engine has a manual override of the pinned version (${this.pinnedVersion}) by version ${this.runningVersion}.`;
    }
    if (this.configuredVersionIsOverridden) {
      return `Configured plugin version (${this.pluginVersion}) is overridden by the pinned version (${this.pinnedVersion}).`;
    }
    if (this.hasVersionMismatch) {
      return `Pinned plugin version (${this.pinnedVersion}) is overridden by the running version (${this.runningVersion}).`;
    }
    return null;
  }

  /**
   * Returns the version mismatch alert configuration, or null if no alert should be shown.
   */
  get versionMismatchAlert(): { title: string; description: string; showReloadButton: boolean } | null {
    if (this.hasVersionMismatch) {
      return {
        title: 'Version mismatch detected',
        description: `This plugin is configured to use version ${this.pluginVersion} but is currently running version ${this.runningVersion}. Reload the plugin to sync the running version with the configured version.`,
        showReloadButton: true,
      };
    }
    return null;
  }

  /**
   * Returns the override pinned alert configuration, or null if no alert should be shown.
   */
  get overridePinnedAlert(): { title: string; description: string; showReloadButton: boolean } | null {
    if (this.willOverridePinnedVersion) {
      return {
        title: 'Override pinned version',
        description: `You have selected ${this.selectedPluginVersion}, but version ${this.pinnedVersion} is pinned for this plugin. Updating to this version will override the pinned version for this mount.`,
        showReloadButton: false,
      };
    }
    return null;
  }

  /**
   * Checks if the selected version would override the pinned version
   */
  get willOverridePinnedVersion(): boolean {
    if (!this.selectedPluginVersion || !this.pinnedVersion) {
      return false;
    }

    // Remove (Pinned) suffix if it exists for comparison
    const selectedVersion = this.selectedPluginVersion.replace(' (Pinned)', '');
    return selectedVersion !== this.pinnedVersion;
  }

  /**
   * Filtered version options that exclude the current running version and show "(Pinned)" for pinned versions.
   * This prevents users from selecting the same version they're already running.
   */
  get filteredVersionOptions(): string[] {
    const { versions } = this.args.model;
    const currentRunningVersion = this.runningVersion;
    const pinnedVersion = this.pinnedVersion;

    if (!versions || !Array.isArray(versions)) {
      return [];
    }

    // Filter out empty strings and the current running version so users can only select different versions
    const filteredVersions = versions.filter(
      (version) => version !== '' && version !== currentRunningVersion
    );

    // Add "(Pinned)" label to pinned versions
    return filteredVersions.map((version) =>
      pinnedVersion && version === pinnedVersion ? `${version} (Pinned)` : version
    );
  }

  /**
   * Reloads the plugin mount to sync running version with configured version
   */
  reloadPlugin = task(async () => {
    try {
      const mountPath = this.args.model.secretsEngine.id;
      await this.api.sys.pluginsReloadBackends({ mounts: [mountPath] });

      this.flashMessages.success(
        'Plugin reloaded successfully. The running version should now match the configured version.',
        {
          title: 'Plugin reloaded',
        }
      );

      // Refresh version data after reload
      const { secretsEngine } = this.args.model;
      const pluginName = secretsEngine.type;

      const [pinnedVersion, mountInfo] = await Promise.all([
        this.api.sys.pluginsCatalogPinsReadPinnedVersion(pluginName, 'secret').catch(() => {
          // Silently handle errors - pins are optional
          return null;
        }),
        this.api.sys.internalUiReadMountInformation(mountPath),
      ]);

      // Notify parent component to refresh the route data
      if (this.args.onReload) {
        this.args.onReload({
          pinnedVersion: pinnedVersion?.version || null,
          pluginVersion: mountInfo?.plugin_version || null,
          runningVersion: mountInfo?.running_plugin_version || null,
        });
      }
    } catch (error) {
      const message = await this.api.parseError(error);
      this.flashMessages.danger(`Failed to reload plugin: ${message.message}`, {
        title: 'Reload Failed',
      });
    }
  });

  @action
  onVersionSelect(event: Event) {
    const target = event.target as HTMLSelectElement;
    this.selectedPluginVersion = target.value || null;
  }
}
