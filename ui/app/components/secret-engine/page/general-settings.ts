/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { convertToSeconds } from 'core/utils/duration-utils';
import { task } from 'ember-concurrency';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type Router from '@ember/routing/router';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type UnsavedChangesService from 'vault/services/unsaved-changes';

const CHARACTER_LIMIT = 500;

/**
 * @module GeneralSettingsComponent is used to configure the SSH secret engine.
 *
 * @example
 * ```js
 * <Secrets:Page:GeneralSettings
 *    @model={{this.model}}
 *  />
 * ```
 *
 * @param {string} secretsEngine - secrets engine resource
 * @param {string} versions - list of versions for a given secret engine
 */

interface Args {
  model: {
    secretsEngine: SecretsEngineResource;
    versions: string[];
    pinnedVersion: string | null;
  };
}

export default class GeneralSettingsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly unsavedChanges: UnsavedChangesService;

  @tracked errorMessage: string | string[] | null = null;
  @tracked errors: string[] = [];
  @tracked invalidFormAlert: string | null = null;
  @tracked showUnsavedChangesModal = false;
  @tracked pinnedVersion: string | null = null;

  @tracked defaultLeaseUnit = '';
  @tracked maxLeaseUnit = '';

  originalModel = JSON.parse(JSON.stringify(this.args.model));

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.pinnedVersion = args.model.pinnedVersion;
  }

  /**
   * Handles version data updates from the version card component
   * Updates both tracked properties and the original model reference
   */
  @action
  handleVersionReload(data: {
    pinnedVersion: string | null;
    pluginVersion: string | null;
    runningVersion: string | null;
  }) {
    // Update tracked pinnedVersion
    this.pinnedVersion = data.pinnedVersion;

    // Update the secretsEngine properties directly
    this.args.model.secretsEngine.running_plugin_version = data.runningVersion || '';
    this.args.model.secretsEngine.plugin_version = data.pluginVersion || '';

    // Update the original model reference so hasPluginVersionChanged works correctly
    this.originalModel.secretsEngine.running_plugin_version = data.runningVersion;
    this.originalModel.secretsEngine.plugin_version = data.pluginVersion;
  }

  /**
   * Creates a model object with version data for the version component
   */
  get versionModel() {
    return {
      secretsEngine: this.args.model.secretsEngine,
      versions: this.args.model.versions,
      pinnedVersion: this.pinnedVersion,
    };
  }

  get modalChangedFields() {
    const changedFieldsCopy = [...this.unsavedChanges.changedFields];
    const configIndex = this.unsavedChanges.changedFields.indexOf('config');

    if (configIndex === -1) return this.unsavedChanges.changedFields;

    changedFieldsCopy[configIndex] = 'Secrets duration';

    return changedFieldsCopy;
  }

  get configRoute() {
    const engine = this.args.model.secretsEngine;
    const isKvv2 = engine.version === 2 && engine.effectiveEngineType === 'kv';
    const engineMetadata = engineDisplayData(engine.effectiveEngineType);

    // Kvv2 is configurable but shares metadata with Kvv1 so isConfigurable is left unset
    if (engineMetadata.isConfigurable || isKvv2) {
      return engineMetadata.configRoute || 'configuration.plugin-settings';
    } else {
      return false;
    }
  }

  validateTtl(ttlValue: FormDataEntryValue | number | null) {
    if (isNaN(Number(ttlValue))) {
      return false;
    }

    return true;
  }

  validateDescription(description: FormDataEntryValue | string) {
    return description?.toString().length <= CHARACTER_LIMIT;
  }

  validateForm() {
    const { defaultLeaseTime, maxLeaseTime, description } = this.formData;

    const errorMessages = [];

    if (!this.validateTtl(defaultLeaseTime) || !this.validateTtl(maxLeaseTime)) {
      errorMessages.push('TTL should only contain numbers.');
    }

    if (description && !this.validateDescription(description)) {
      const charactersExceedBy = description?.toString().length - CHARACTER_LIMIT;
      errorMessages.push(`Engine description exceeds character limit by ${charactersExceedBy}.`);
    }

    this.errors = errorMessages;

    if (errorMessages.length) return false;
    return true;
  }

  hasTtlValueChanged(ttlTime: number, ttlUnit: string, ttlKey: 'max_lease_ttl' | 'default_lease_ttl') {
    const defaultLeaseInSecs = convertToSeconds(ttlTime, ttlUnit);
    if (defaultLeaseInSecs === this?.originalModel?.secretsEngine?.config[ttlKey]) {
      return false;
    }

    return true;
  }

  get formData() {
    const form = document.getElementById('general-settings-form');
    const fd = new FormData(form as HTMLFormElement);
    const fdDefaultLeaseTime = Number(fd.get('default_lease_ttl-time'));
    const fdDefaultLeaseUnit = fd.get('default_lease_ttl-unit')?.toString() || 's';
    const fdMaxLeaseTime = Number(fd.get('max_lease_ttl-time'));
    const fdMaxLeaseUnit = fd.get('max_lease_ttl-unit')?.toString() || 's';

    return {
      defaultLeaseTime: fdDefaultLeaseTime,
      defaultLeaseUnit: fdDefaultLeaseUnit,
      maxLeaseTime: fdMaxLeaseTime,
      maxLeaseUnit: fdMaxLeaseUnit,
      description: fd.get('description'),
      version: fd.get('plugin-version'),
    };
  }

  formatTuneParams() {
    const { defaultLeaseTime, defaultLeaseUnit, maxLeaseTime, maxLeaseUnit, description, version } =
      this.formData;

    const hasDefaultTtlValueChanged = this.hasTtlValueChanged(
      defaultLeaseTime,
      defaultLeaseUnit,
      'default_lease_ttl'
    );

    const hasMaxTtlValueChanged = this.hasTtlValueChanged(maxLeaseTime, maxLeaseUnit, 'max_lease_ttl');
    const hasDescriptionChanged = description !== this?.originalModel?.secretsEngine?.description;

    // Clean the version string by removing "(Pinned)" label
    const cleanVersion = version ? version.toString().replace(' (Pinned)', '') : null;
    const hasPluginVersionChanged =
      cleanVersion && cleanVersion !== this?.originalModel?.secretsEngine?.running_plugin_version;

    const defaultLeaseTtl = hasDefaultTtlValueChanged ? `${defaultLeaseTime}${defaultLeaseUnit}` : undefined;
    const maxLeaseTtl = hasMaxTtlValueChanged ? `${maxLeaseTime}${maxLeaseUnit}` : undefined;
    const pluginVersion = hasPluginVersionChanged ? cleanVersion : undefined;
    const pluginDescription = hasDescriptionChanged ? description : undefined;

    // Determine if we need to override pinned version
    // This is required when updating to a version that differs from the pinned version
    const overridePinnedVersion =
      hasPluginVersionChanged && this.pinnedVersion && this.pinnedVersion !== cleanVersion;

    // Determine if user is selecting the pinned version (to send override_pinned_version: false)
    const selectedPinnedVersion =
      hasPluginVersionChanged && this.pinnedVersion && this.pinnedVersion === cleanVersion;

    return {
      defaultLeaseTtl,
      maxLeaseTtl,
      pluginVersion,
      pluginDescription,
      hasPluginVersionChanged,
      overridePinnedVersion,
      selectedPinnedVersion,
    };
  }

  saveGeneralSettings = task(async (event?) => {
    // event is an optional arg because saveGeneralSettings can be called in the closeAndHandle function.
    // There are instances where we will save in the modal where that doesn't require an event.
    if (event) event.preventDefault();

    if (!this.validateForm()) return;

    try {
      const {
        defaultLeaseTtl,
        maxLeaseTtl,
        pluginVersion,
        pluginDescription,
        hasPluginVersionChanged,
        overridePinnedVersion,
        selectedPinnedVersion,
      } = this.formatTuneParams();

      const tunePayload = {
        description: pluginDescription as string | undefined,
        default_lease_ttl: defaultLeaseTtl,
        max_lease_ttl: maxLeaseTtl,
        plugin_version: selectedPinnedVersion ? undefined : pluginVersion, // Don't send plugin_version if user is selecting the pinned version to avoid API error
        override_pinned_version: overridePinnedVersion ? true : selectedPinnedVersion ? false : undefined,
      };

      /*
       * Using raw request instead of client API methods.
       *
       * The client API methods (mountsTuneConfigurationParameters, mountsTuneConfigurationParametersRaw)
       * do not currently include the new override_pinned_version flag in the request body.
       *
       * Temporary Solution: Use raw request (this.api.request.post) to ensure the complete
       * payload with override_pinned_version is always sent correctly.
       *
       * TODO: Fast follow-up PR to update this to use proper client API methods once the
       * client version is bumped.
       */
      await this.api.request.post(`/sys/mounts/${this.args?.model?.secretsEngine?.id}/tune`, tunePayload);

      // If plugin version changed, reload the mount to ensure the new version is active
      if (hasPluginVersionChanged) {
        await this.reloadMount();
      }

      this.flashMessages.success('Engine settings successfully updated.', { title: 'Configuration saved' });

      this.unsavedChanges.showModal = false;
      this.unsavedChanges.transition('vault.cluster.secrets.backend.configuration.general-settings');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  });

  /**
   * Reloads the mount by calling the plugin reload API
   * This is necessary when the plugin version changes to ensure the new version is active
   */
  private async reloadMount(): Promise<void> {
    try {
      const mountPath = this.args?.model?.secretsEngine?.id;
      if (!mountPath) {
        throw new Error('Mount path is required to reload mount');
      }

      // Use the dedicated API service method to reload the plugin mount
      await this.api.sys.pluginsReloadBackends({ mounts: [mountPath] });
    } catch (error) {
      // The TUNE was successful, but the reload failed
      this.errorMessage =
        'Plugin version was updated successfully, but the mount could not be automatically reloaded. You may need to restart Vault or manually reload the plugin.';
    }
  }

  @action
  discardChanges() {
    const currentRouteName = this.router.currentRouteName;
    this.router.transitionTo(currentRouteName);
  }
}
