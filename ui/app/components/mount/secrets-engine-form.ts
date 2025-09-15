/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { set } from '@ember/object';
import type FlashMessagesService from 'ember-cli-flash/services/flash-messages';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type Router from '@ember/routing/router';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type { ValidationMap } from 'vault/vault/app-types';
import { isAddonEngine } from 'vault/utils/all-engines-metadata';

interface Args {
  model: SecretsEngineForm;
  onMountSuccess?: (type: string, path: string, useEngineRoute: boolean) => void;
}

/**
 * @module Mount::SecretsEngineForm
 * Modern component for mounting secrets engines using the SecretsEngineForm.
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

  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked errorMessage: string | string[] = '';

  get mountForm(): SecretsEngineForm {
    return this.args.model;
  }

  @action
  onKeyUp(name: string, value: string) {
    set(this.mountForm.data, name, value);
  }

  async saveKvConfig(path: string, formData: SecretsEngineForm['data']) {
    const { options, kv_config = {} } = formData;
    const { max_versions, cas_required, delete_version_after } = kv_config;
    const isKvV2 = options?.version === 2 && ['kv', 'generic'].includes(this.mountForm.normalizedType);
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
    const mountModel = this.mountForm;
    const { type } = mountModel;
    const { path } = mountModel.data;

    // Only submit form if validations pass
    const { isValid, state, invalidFormMessage, data } = mountModel.toJSON();
    if (!isValid) {
      this.modelValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return;
    }

    this.errorMessage = '';
    this.modelValidations = null;
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
    const { config } = this.mountForm.data;
    config.identity_token_key = Array.isArray(value) ? value[0] : value;
  }

  @action
  goBack() {
    this.router.transitionTo('vault.cluster.secrets.mounts');
  }
}
