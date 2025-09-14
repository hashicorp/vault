/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { convertToSeconds } from 'core/utils/duration-utils';

import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

export const CUSTOM = 'Custom';
export const SYSTEM_DEFAULT = 'System default';

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
  };
}

export default class GeneralSettingsComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked showUnsavedChangesModal = false;
  @tracked changedFields: string[] = [];

  originalModel = JSON.parse(JSON.stringify(this.args.model));

  getUnsavedChanges(newModel: SecretsEngineResource, originalModel: SecretsEngineResource) {
    for (const key in this.args.model.secretsEngine) {
      const secretsEngineKeyType = key as keyof typeof this.args.model.secretsEngine;

      if (secretsEngineKeyType === 'options') {
        return;
      }

      if (secretsEngineKeyType === 'config') {
        const { defaultLeaseTime, defaultLeaseUnit, maxLeaseTime, maxLeaseUnit } = this.getFormData();

        const hasDefaultTtlValueChanged = this.hasTtlValueChanged(
          defaultLeaseTime,
          defaultLeaseUnit,
          'default_lease_ttl'
        );
        const hasMaxTtlValueChanged = this.hasTtlValueChanged(maxLeaseTime, maxLeaseUnit, 'max_lease_ttl');

        if (
          (hasDefaultTtlValueChanged || hasMaxTtlValueChanged) &&
          !this.changedFields.includes('Lease Duration')
        ) {
          this.changedFields.push('Lease Duration');
        }
      } else {
        if (newModel[secretsEngineKeyType] !== originalModel[secretsEngineKeyType]) {
          this.changedFields.push(key);
        }
      }
    }
  }

  validateTtl(ttlValue: FormDataEntryValue | number | null) {
    if (isNaN(Number(ttlValue))) {
      this.errorMessage = 'Only use numbers for this setting.';
      return false;
    }

    return true;
  }

  hasTtlValueChanged(ttlTime: number, ttlUnit: string, ttlKey: 'max_lease_ttl' | 'default_lease_ttl') {
    const defaultLeaseInSecs = convertToSeconds(ttlTime, ttlUnit);
    if (defaultLeaseInSecs === this?.originalModel?.secretsEngine?.config[ttlKey]) {
      return false;
    }

    return true;
  }

  hasDescriptionChanged(description: FormDataEntryValue | null) {
    return description !== this?.originalModel?.secretsEngine?.description;
  }

  hasPluginVersionChanged(version: FormDataEntryValue | null) {
    return version && version !== this?.originalModel?.secretsEngine?.running_plugin_version;
  }

  getFormData() {
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

  hasUnsavedChanges() {
    const { defaultLeaseTime, defaultLeaseUnit, maxLeaseTime, maxLeaseUnit, description, version } =
      this.getFormData();

    const hasDefaultTtlValueChanged = this.hasTtlValueChanged(
      defaultLeaseTime,
      defaultLeaseUnit,
      'default_lease_ttl'
    );
    const hasMaxTtlValueChanged = this.hasTtlValueChanged(maxLeaseTime, maxLeaseUnit, 'max_lease_ttl');

    return (
      hasDefaultTtlValueChanged ||
      hasMaxTtlValueChanged ||
      this.hasDescriptionChanged(description) ||
      this.hasPluginVersionChanged(version)
    );
  }

  formatTuneParams() {
    const { defaultLeaseTime, defaultLeaseUnit, maxLeaseTime, maxLeaseUnit, description, version } =
      this.getFormData();

    const hasDefaultTtlValueChanged = this.hasTtlValueChanged(
      defaultLeaseTime,
      defaultLeaseUnit,
      'default_lease_ttl'
    );

    const hasMaxTtlValueChanged = this.hasTtlValueChanged(maxLeaseTime, maxLeaseUnit, 'max_lease_ttl');

    const defaultLeaseTtl = hasDefaultTtlValueChanged ? `${defaultLeaseTime}${defaultLeaseUnit}` : undefined;
    const maxLeaseTtl = hasMaxTtlValueChanged ? `${maxLeaseTime}${maxLeaseUnit}` : undefined;
    const pluginVersion = this.hasPluginVersionChanged(version) ? version : undefined;
    const pluginDescription = this.hasDescriptionChanged(description) ? description : undefined;

    return {
      defaultLeaseTtl,
      maxLeaseTtl,
      pluginVersion,
      pluginDescription,
    };
  }

  @action
  openUnsavedChangesModal() {
    if (this.hasUnsavedChanges()) {
      this.getUnsavedChanges(this.args?.model?.secretsEngine, this?.originalModel?.secretsEngine);
      this.showUnsavedChangesModal = true;
    } else {
      this.showUnsavedChangesModal = false;
    }
  }

  @action
  closeUnsavedChangesModal() {
    this.showUnsavedChangesModal = !this.showUnsavedChangesModal;
    this.changedFields = [];
  }

  @action
  discardChanges() {
    this.closeUnsavedChangesModal();
    this.router.transitionTo(this.args?.model?.secretsEngine?.backendConfigurationLink);
  }

  saveGeneralSettings = task(async (event) => {
    event.preventDefault();

    const { defaultLeaseTime, maxLeaseTime } = this.getFormData();

    if (!this.validateTtl(defaultLeaseTime) || !this.validateTtl(maxLeaseTime)) {
      this.errorMessage = 'Only use numbers for this setting.';
      return;
    }

    try {
      const { defaultLeaseTtl, maxLeaseTtl, pluginVersion, pluginDescription } = this.formatTuneParams();

      await this.api.sys.mountsTuneConfigurationParameters(this.args?.model?.secretsEngine?.id, {
        description: pluginDescription as string | undefined,
        default_lease_ttl: defaultLeaseTtl,
        max_lease_ttl: maxLeaseTtl,
        plugin_version: pluginVersion as string | undefined,
      });

      this.flashMessages.success('Engine settings successfully updated.');
      this.router.transitionTo(this.args?.model?.secretsEngine?.backendConfigurationLink);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
      this.flashMessages.danger(`Try again or check your network connection. ${message}`);
    }
  });
}
