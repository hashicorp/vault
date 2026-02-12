/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { capitalize } from '@ember/string';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

import type Router from '@ember/routing/router';
import type FlashMessagesService from 'ember-cli-flash/services/flash-messages';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type VersionService from 'vault/services/version';
import { isAddonEngine } from 'vault/utils/all-engines-metadata';
import { getExternalPluginNameFromBuiltin } from 'vault/utils/external-plugin-helpers';
import type { EngineVersionInfo } from 'vault/utils/plugin-catalog-helpers';
import { sortVersions } from 'vault/utils/version-utils';
import type { ValidationMap } from 'vault/vault/app-types';

// Extended config interface for plugin mounting
interface ExtendedMountConfig {
  plugin_version?: string;
  override_pinned_version?: boolean;
  [key: string]: any;
}

enum PluginRegistrationType {
  BUILTIN = 'builtin',
  EXTERNAL = 'external',
}

interface Args {
  model: {
    form: SecretsEngineForm;
    availableVersions?: EngineVersionInfo[];
    hasUnversionedPlugins?: boolean;
    pinnedVersion?: string | null;
  };
  onMountSuccess?: (type: string, path: string, useEngineRoute: boolean) => void;
}

/**
 * @module Mount::SecretsEngineForm
 * Modern component for mounting secrets engines using the SecretsEngineForm.
 *
 * Plugin version handling:
 * - Plugin type (built-in/external) is selected via radio cards
 * - Version dropdown appears only for external plugins
 * - When version changes, onPluginVersionChange updates model and type
 * - The model's handlePluginVersionChange method updates the type to use external plugin name if needed
 *
 * @example
 * ```hbs
 * <Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />
 * ```
 */
export default class MountSecretsEngineFormComponent extends Component<Args> {
  @service declare flashMessages: FlashMessagesService;
  @service declare api: ApiService;
  @service declare capabilities: CapabilitiesService;
  @service declare router: Router;
  @service declare version: VersionService;

  @tracked formValidations: ValidationMap | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked errorMessage: string | string[] = '';
  @tracked pluginRegistrationType: 'builtin' | 'external' = PluginRegistrationType.BUILTIN;
  @tracked selectedPluginVersion = '';

  _originalBuiltinType = '';

  // Plugin registration type constants
  PluginRegistrationType = PluginRegistrationType;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    // Store the original builtin type for restoration when switching back from external
    this._originalBuiltinType = this.args.model.form.normalizedType;

    // Initialize plugin version
    this.configObject.plugin_version = '';
  }

  // Helper to get config object with proper typing
  get configObject() {
    return this.args.model.form.data.config as ExtendedMountConfig;
  }

  // Check if current plugin registration type is builtin
  get isBuiltinPlugin(): boolean {
    return this.pluginRegistrationType === PluginRegistrationType.BUILTIN;
  }

  // Check if current plugin registration type is external
  get isExternalPlugin(): boolean {
    return this.pluginRegistrationType === PluginRegistrationType.EXTERNAL;
  }

  get breadcrumbs() {
    const breadcrumbs: { label: string; route?: string; icon?: string }[] = [
      { label: 'Vault', route: 'vault.cluster', icon: 'vault' },
      { label: 'Secrets engines', route: 'vault.cluster.secrets.backends' },
      { label: 'Enable secrets engine', route: 'vault.cluster.secrets.enable' },
    ];

    if (this.args?.model?.form?.normalizedType) {
      breadcrumbs.push({ label: capitalize(this.args?.model?.form?.normalizedType) });
    }

    return breadcrumbs;
  }

  get pluginTypeOptions() {
    return [
      {
        type: this.PluginRegistrationType.BUILTIN,
        icon: 'server',
        label: 'Built-in plugin',
        description:
          'Preregistered plugins shipped with Vault. The plugin version is tied to your Vault version and cannot be specified.',
        dataTestAttr: 'builtin',
        disabled: false,
        showBadge: false,
        showAlert: false,
      },
      {
        type: this.PluginRegistrationType.EXTERNAL,
        icon: 'download',
        label: 'External plugin',
        description:
          'External plugins manually registered in your plugin catalog. If multiple versions are registered, you can specify which version to enable.',
        dataTestAttr: 'external',
        disabled: this.shouldDisableExternal,
        showBadge: !this.version.isEnterprise,
        showAlert: this.shouldShowNoExternalVersionsMessage,
      },
    ];
  }

  @action
  onKeyUp(name: string, value: string) {
    (this.args.model.form.data as any)[name] = value;
  }

  // Get pinned version for current external plugin
  get pinnedVersionForCurrentPlugin(): string | null {
    return this.args.model.pinnedVersion || null;
  }

  // Check if External radio should be disabled (only built-in versions available or no Enterprise license)
  get shouldDisableExternal(): boolean {
    // Disable if no Enterprise license
    if (!this.version.isEnterprise) {
      return true;
    }

    // Disable if no external versions available
    if (!this.args.model.availableVersions) {
      return true;
    }
    return !this.args.model.availableVersions.some((version) => !version.isBuiltin);
  }

  // Get the external plugin name for the current engine type
  get externalPluginName(): string | null {
    const engineType = this.args.model.form.normalizedType;
    return getExternalPluginNameFromBuiltin(engineType);
  }

  // Check if we should show info message for disabled external card due to no external versions
  get shouldShowNoExternalVersionsMessage(): boolean {
    return this.version.isEnterprise && this.shouldDisableExternal;
  }

  // Check if plugin version field should be shown
  get shouldShowPluginVersionField(): boolean {
    // Only show for external plugins
    if (!this.isExternalPlugin) {
      return false;
    }

    // Only show if we have external versions
    return this.getExternalVersionList().length > 0;
  }

  // Get external version options with default pinned version
  get filteredVersionOptions(): string[] {
    const versionList = this.getExternalVersionList();
    if (versionList.length === 0) {
      return [];
    }

    // Sort versions with pinned version first if it exists and pins are loaded
    const pinnedVersion = this.pinnedVersionForCurrentPlugin;
    if (pinnedVersion && versionList.includes(pinnedVersion)) {
      const sortedVersions = [pinnedVersion, ...versionList.filter((v) => v !== pinnedVersion)];
      return sortedVersions;
    }

    // Sort by semantic version (highest first)
    return sortVersions(versionList, true);
  }

  // Extract common version list filtering logic
  private getExternalVersionList(): string[] {
    const versions = this.args.model.availableVersions;
    if (!versions || !Array.isArray(versions)) {
      return [];
    }

    // Filter external versions and exclude empty strings
    const externalVersions = versions.filter((version) => !version.isBuiltin && version.version !== '');
    return externalVersions.map((version) => version.version);
  }

  // Check if the currently selected version differs from the pinned version
  get shouldShowPinWarning(): boolean {
    if (!this.isExternalPlugin) {
      return false;
    }

    const pinnedVersion = this.pinnedVersionForCurrentPlugin;
    const currentVersion = this.selectedPluginVersion;

    // If there's no pinned version, no warning needed
    if (!pinnedVersion) {
      return false;
    }

    // If the current version is undefined/empty, no warning needed
    if (!currentVersion) {
      return false;
    }

    // Show warning if there's a pinned version and it's different from current selection
    return pinnedVersion !== currentVersion;
  }

  // Update override flag based on version selection
  updateOverridePinnedVersionFlag() {
    // For builtin plugins, ensure override flag is not sent
    if (this.isBuiltinPlugin) {
      delete this.configObject.override_pinned_version;
      return;
    }

    const pinnedVersion = this.pinnedVersionForCurrentPlugin;
    const currentVersion = this.configObject.plugin_version;

    if (pinnedVersion && currentVersion && pinnedVersion !== currentVersion) {
      // User selected a version different from pinned - include both parameters
      this.configObject.plugin_version = currentVersion;
      this.configObject.override_pinned_version = true;
    } else if (pinnedVersion && currentVersion === pinnedVersion) {
      // User is using the pinned version - omit both parameters (backend will use pin)
      delete this.configObject.plugin_version;
      delete this.configObject.override_pinned_version;
    } else {
      // No pinned version exists - include plugin_version but not override flag
      this.configObject.plugin_version = currentVersion;
      delete this.configObject.override_pinned_version;
    }
  }

  // Save KV configuration if applicable
  @action
  async saveKvConfig(path: string, formData: SecretsEngineForm['data']) {
    const { options, kv_config = {} } = formData;
    const { max_versions, cas_required, delete_version_after } = kv_config;
    const isKvV2 = options?.version === 2 && ['kv', 'generic'].includes(this.args.model.form.normalizedType);
    const hasConfig = max_versions || cas_required || delete_version_after;

    if (isKvV2 && hasConfig) {
      try {
        const { canUpdate } = await this.capabilities.for('kvConfig', { path });
        if (canUpdate) {
          await this.api.secrets.kvV2Configure(path, kv_config);
        } else {
          this.flashMessages.warning(
            'You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.'
          );
        }
      } catch (e) {
        const { message } = await this.api.parseError(e);
        this.flashMessages.warning(
          `The secret engine was mounted, but the configuration settings were not saved. ${message}`
        );
      }
    }
  }

  // Handle mount errors
  @action
  async onMountError(status: number, errors: unknown[] | undefined, message: string) {
    if (status === 403) {
      this.flashMessages.danger(
        'You do not have access to the sys/mounts endpoint. The secret engine was not mounted.'
      );
    } else if (errors) {
      this.errorMessage = errors.map((e) => {
        if (typeof e === 'object' && e !== null) {
          const errorObj = e as { title?: string; message?: string };
          return errorObj.title || errorObj.message || JSON.stringify(e);
        }
        return String(e);
      });
    } else if (message) {
      this.errorMessage = message;
    } else {
      this.errorMessage = 'An error occurred, check the vault logs.';
    }
  }

  @task
  *mountBackend(event: Event) {
    event.preventDefault();
    const mountModel = this.args.model.form;
    const { type } = mountModel;
    const { path } = mountModel.data;

    // Handle plugin version change before validation in case onKeyUp wasn't called
    if (this.args.model.availableVersions && this.configObject.plugin_version) {
      mountModel.handlePluginVersionChange(this.args.model.availableVersions);
    }

    // Only submit form if validations pass
    const { isValid, state, invalidFormMessage, data } = mountModel.toJSON();

    if (!isValid) {
      this.formValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return;
    }

    this.errorMessage = '';
    this.formValidations = null;
    this.invalidFormAlert = null;

    try {
      // Mount the secrets engine
      yield this.api.sys.mountsEnableSecretsEngine(path, data);

      // Save KV config if applicable
      yield this.saveKvConfig(path, data);

      this.flashMessages.success(`Successfully mounted the ${mountModel.type} secrets engine at ${path}.`);

      // Determine if we should use engine routes
      const version = data.options?.version;
      const useEngineRoute = isAddonEngine(mountModel.normalizedType, Number(version));

      // Call success callback or navigate
      if (this.args.onMountSuccess) {
        this.args.onMountSuccess(type, path, useEngineRoute);
      } else {
        // Default navigation
        if (useEngineRoute) {
          this.router.transitionTo('vault.cluster.secrets.backend.index', path);
        } else {
          this.router.transitionTo('vault.cluster.secrets.backend.list-root', path);
        }
      }
    } catch (error) {
      const { status, response, message } = yield this.api.parseError(error);
      this.onMountError(status, response.errors, message);
    }
  }

  @action
  handleIdentityTokenKeyChange(value: string[] | string): void {
    // if array, it's coming from the search-select component, otherwise it hit the fallback component and will come in as a string.
    const { config } = this.args.model.form.data;
    config.identity_token_key = Array.isArray(value) ? value[0] : value;
  }

  @action
  goBack() {
    this.router.transitionTo('vault.cluster.secrets.enable');
  }

  // Set default plugin version for external plugins
  private setDefaultPluginVersion() {
    const versionList = this.getExternalVersionList();
    if (versionList.length === 0) {
      this.configObject.plugin_version = '';
      this.selectedPluginVersion = '';
      return;
    }

    // Check for pinned version first (pins should be loaded from constructor)
    const pinnedVersion = this.pinnedVersionForCurrentPlugin;

    if (pinnedVersion && versionList.includes(pinnedVersion)) {
      // Use pinned version if available in catalog
      this.selectedPluginVersion = pinnedVersion;
      this.configObject.plugin_version = pinnedVersion;
    } else {
      // Use highest semantic version from catalog if no pin or pin not available
      const sortedVersions = sortVersions(versionList, true);
      const topVersion = sortedVersions[0] || '';
      this.selectedPluginVersion = topVersion;
      this.configObject.plugin_version = topVersion;
    }

    // Update override flag based on final selection
    this.updateOverridePinnedVersionFlag();
  }

  @action
  setPluginType(type: 'builtin' | 'external') {
    this.pluginRegistrationType = type;

    // Update the model type based on selection
    if (type === PluginRegistrationType.BUILTIN) {
      // Use the stored original built-in type (e.g., 'keymgmt')
      this.args.model.form.type = this._originalBuiltinType;
      // Clear plugin version and override flag for built-in plugins
      this.selectedPluginVersion = '';
      this.configObject.plugin_version = '';
      delete this.configObject.override_pinned_version;
    } else {
      // Use the external plugin name (e.g., 'vault-plugin-secrets-keymgmt')
      this.args.model.form.type = this.externalPluginName || '';

      // Set appropriate plugin version based on available versions
      this.setDefaultPluginVersion();
    }
  }

  @action
  onPluginVersionChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    const value = target.value;

    this.selectedPluginVersion = value;
    this.configObject.plugin_version = value;

    // Update override flag when user manually changes version
    this.updateOverridePinnedVersionFlag();

    // Update the type based on the selected version
    if (this.args.model.availableVersions) {
      this.args.model.form.handlePluginVersionChange(this.args.model.availableVersions);
    }
  }
}
