/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { convertToSeconds } from 'core/utils/duration-utils';
import { action } from '@ember/object';

import type Router from '@ember/routing/router';
import type FlashMessageService from 'vault/services/flash-messages';
import type ApiService from 'vault/services/api';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
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

  @tracked defaultLeaseUnit = '';
  @tracked maxLeaseUnit = '';

  originalModel = JSON.parse(JSON.stringify(this.args.model));

  get modalChangedFields() {
    const changedFieldsCopy = [...this.unsavedChanges.changedFields];
    const configIndex = this.unsavedChanges.changedFields.indexOf('config');

    if (configIndex === -1) return this.unsavedChanges.changedFields;

    changedFieldsCopy[configIndex] = 'Secrets duration';

    return changedFieldsCopy;
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
    const hasPluginVersionChanged =
      version && version !== this?.originalModel?.secretsEngine?.running_plugin_version;

    const defaultLeaseTtl = hasDefaultTtlValueChanged ? `${defaultLeaseTime}${defaultLeaseUnit}` : undefined;
    const maxLeaseTtl = hasMaxTtlValueChanged ? `${maxLeaseTime}${maxLeaseUnit}` : undefined;
    const pluginVersion = hasPluginVersionChanged ? version : undefined;
    const pluginDescription = hasDescriptionChanged ? description : undefined;

    return {
      defaultLeaseTtl,
      maxLeaseTtl,
      pluginVersion,
      pluginDescription,
    };
  }

  saveGeneralSettings = task(async (event?) => {
    // event is an optional arg because saveGeneralSettings can be called in the closeAndHandle function.
    // There are instances where we will save in the modal where that doesn't require an event.
    if (event) event.preventDefault();

    if (!this.validateForm()) return;

    try {
      const { defaultLeaseTtl, maxLeaseTtl, pluginVersion, pluginDescription } = this.formatTuneParams();

      await this.api.sys.mountsTuneConfigurationParameters(this.args?.model?.secretsEngine?.id, {
        description: pluginDescription as string | undefined,
        default_lease_ttl: defaultLeaseTtl,
        max_lease_ttl: maxLeaseTtl,
        plugin_version: pluginVersion as string | undefined,
      });

      this.flashMessages.success('Engine settings successfully updated.', { title: 'Configuration saved' });

      this.unsavedChanges.showModal = false;
      this.router.transitionTo(this.router.currentRouteName);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.errorMessage = message;
    }
  });

  @action
  discardChanges() {
    const currentRouteName = this.router.currentRouteName;
    this.router.transitionTo(currentRouteName);
  }
}
