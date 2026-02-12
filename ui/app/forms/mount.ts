/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import Form from 'vault/forms/form';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import FormField from 'vault/utils/forms/field';
import { WHITESPACE_WARNING } from 'vault/utils/forms/validators';

import type { Validations } from 'vault/app-types';
import type { SecretsEngineFormData } from 'vault/secrets/engine';
import type { EngineVersionInfo } from 'vault/utils/plugin-catalog-helpers';
import type { AuthMethodFormData } from 'vault/vault/auth/methods';

type ConfigWithPluginVersion = {
  plugin_version?: string;
  [key: string]: any;
};

// common fields and validations shared between secrets engine and auth methods (mounts)
// used in form classes for consistency and to avoid duplication
export default class MountForm<T extends SecretsEngineFormData | AuthMethodFormData> extends Form<T> {
  @tracked declare type: string;

  validations: Validations = {
    path: [
      { type: 'presence', message: "Path can't be blank." },
      {
        type: 'containsWhiteSpace',
        message: WHITESPACE_WARNING('path'),
        level: 'warn',
      },
    ],
  };

  fields = {
    path: new FormField('path', 'string'),
    description: new FormField('description', 'string', { editType: 'textarea' }),
    listingVisibility: new FormField('config.listing_visibility', 'boolean', {
      label: 'Use as preferred UI login method',
      editType: 'toggleButton',
      helperTextEnabled:
        'This mount will be included in the unauthenticated UI login endpoint and display as a preferred login method.',
      helperTextDisabled:
        'Turn on the toggle to use this auth mount as a preferred login method during UI login.',
    }),
    local: new FormField('local', 'boolean', {
      helpText:
        'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
    }),
    sealWrap: new FormField('seal_wrap', 'boolean', {
      helpText:
        'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
    }),
    defaultLeaseTtl: new FormField('config.default_lease_ttl', 'string', {
      label: 'Default Lease TTL',
      editType: 'ttl',
    }),
    maxLeaseTtl: new FormField('config.max_lease_ttl', 'string', { label: 'Max Lease TTL', editType: 'ttl' }),
    auditNonHmacRequestKeys: new FormField('config.audit_non_hmac_request_keys', 'string', {
      label: 'Request keys excluded from HMACing in audit',
      editType: 'stringArray',
      helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
    }),
    auditNonHmacResponseKeys: new FormField('config.audit_non_hmac_response_keys', 'string', {
      label: 'Response keys excluded from HMACing in audit',
      editType: 'stringArray',
      helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
    }),
    passthroughRequestHeaders: new FormField('config.passthrough_request_headers', 'string', {
      label: 'Allowed passthrough request headers',
      helpText: 'Headers to allow and pass from the request to the backend',
      editType: 'stringArray',
    }),
    allowedResponseHeaders: new FormField('config.allowed_response_headers', 'string', {
      label: 'Allowed response headers',
      helpText: 'Headers to allow, allowing a plugin to include them in the response.',
      editType: 'stringArray',
    }),
    pluginVersion: new FormField('config.plugin_version', 'string', {
      label: 'Plugin version',
      subText:
        'Specifies the semantic version of the plugin to use, e.g. "v1.0.0". If unspecified, the server will select any matching un-versioned plugin that may have been registered, the latest versioned plugin registered, or a built-in plugin in that order of precedence.',
    }),
  };

  // normalizes type for UI configuration purposes by:
  // 1. stripping `ns_` prefix (for namespaced types)
  // 2. mapping external plugins to their builtin equivalents for consistent UI experience
  get normalizedType() {
    const baseType = (this.type || '').replace(/^ns_/, '');
    return getEffectiveEngineType(baseType);
  }

  /**
   * Sets up plugin version configuration for the form.
   * Since plugin version is handled manually in the template, this method
   * only manages the data model setup.
   *
   * @param availableVersions - Array of available plugin versions
   */
  setupPluginVersionField(availableVersions: EngineVersionInfo[] | null | undefined) {
    if (!availableVersions || availableVersions.length === 0) {
      return;
    }

    // Initialize plugin_version as empty (default option)
    (this.data.config as ConfigWithPluginVersion).plugin_version = '';
  }

  /**
   * Updates the form data with the selected plugin version information.
   * For external plugins, this also updates the engine type to match the plugin name,
   * enabling proper mounting of external plugins with their specific names.
   *
   * @param versionInfo - The selected version information containing plugin name, version, and builtin status
   */
  setPluginVersionData(versionInfo: EngineVersionInfo) {
    // Set the version in config
    (this.data.config as ConfigWithPluginVersion).plugin_version = versionInfo.version;

    // For external plugins, update the type to the plugin name
    if (!versionInfo.isBuiltin) {
      this.type = versionInfo.pluginName;
    }
  }

  /**
   * Locates the version information object that matches a user-selected value.
   * This bridges the gap between the selected version value and the underlying plugin metadata needed for mounting.
   *
   * @param selectedValue - The selected version value from the UI dropdown (actual semantic version)
   * @param availableVersions - Available version options from the plugin catalog
   * @returns The matching version info or undefined if no match found
   */
  findVersionByLabel(
    selectedValue: string,
    availableVersions: EngineVersionInfo[]
  ): EngineVersionInfo | undefined {
    // Handle the empty value (default option) - return undefined so we don't send plugin_version
    if (!selectedValue || selectedValue === '') {
      return undefined;
    }

    return availableVersions.find((v) => v.version === selectedValue);
  }

  /**
   * Handles plugin version changes and updates the type if needed
   * This method should be called whenever the plugin version field changes
   */
  handlePluginVersionChange(availableVersions: EngineVersionInfo[]) {
    const config = this.data.config as ConfigWithPluginVersion;
    const selectedVersion = config?.plugin_version;
    if (!selectedVersion || !availableVersions) {
      return;
    }

    // Find the selected version info
    const selectedVersionInfo = this.findVersionByLabel(selectedVersion, availableVersions);
    if (selectedVersionInfo) {
      this.setPluginVersionData(selectedVersionInfo);
    }
  }

  toJSON() {
    const { config } = this.data;
    const data = {
      type: this.type,
      ...this.data,
      config: {
        ...(config || {}),
        force_no_cache: config?.force_no_cache ?? false,
        listing_visibility: config?.listing_visibility ? 'unauth' : 'hidden',
      },
    };

    // Remove plugin_version if it's empty (let server choose default)
    const configWithPluginVersion = data.config as ConfigWithPluginVersion;
    if (!configWithPluginVersion.plugin_version || configWithPluginVersion.plugin_version === '') {
      delete configWithPluginVersion.plugin_version;
    }

    // options are only relevant for kv/generic engines
    if (!['kv', 'generic'].includes(this.type)) {
      delete data.options;
    }

    return super.toJSON(data);
  }
}
