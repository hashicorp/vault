/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action, set } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { presence } from 'vault/utils/forms/validators';
import { filterEnginesByMountCategory, isAddonEngine } from 'vault/utils/all-engines-metadata';
import { assert } from '@ember/debug';

import type FlashMessageService from 'vault/services/flash-messages';
import type Store from '@ember-data/store';
import type AuthMethodForm from 'vault/forms/auth/method';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type CapabilitiesService from 'vault/services/capabilities';
import type ApiService from 'vault/services/api';
import type { ApiError } from '@ember-data/adapter/error';
import type { ValidationMap } from 'vault/vault/app-types';

/**
 * @module MountBackendForm
 * The `MountBackendForm` is used to mount either a secret or auth backend.
 *
 * @example ```js
 *   <MountBackendForm @mountCategory="secret" @onMountSuccess={{this.onMountSuccess}} />```
 *
 * @param {function} onMountSuccess - A function that transitions once the Mount has been successfully posted.
 * @param {string} mountCategory - The type of engine to mount, either 'secret' or 'auth'.
 *
 */

type MountModel = SecretsEngineForm | AuthMethodForm;

interface Args {
  mountModel: MountModel;
  mountCategory: 'secret' | 'auth';
  onMountSuccess: (type: string, path: string, useEngineRoute: boolean) => void;
}

export default class MountBackendForm extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly api: ApiService;

  // validation related properties
  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormAlert: string | null = null;

  @tracked errorMessage: string | string[] = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    assert(`@mountCategory is required. Must be "auth" or "secret".`, presence(this.args.mountCategory));
  }

  checkPathChange(backendType: string) {
    if (!backendType) return;
    const { data } = this.args.mountModel;
    // mountCategory is usually 'secret' or 'auth', but sometimes an empty string is passed in (like when we click the cancel button).
    // In these cases, we should default to returning auth methods.
    const mountsByType = filterEnginesByMountCategory({
      mountCategory: this.args.mountCategory ?? 'auth',
      isEnterprise: true,
    }).map((engine) => engine.type);
    // if the current path has not been altered by user,
    // change it here to match the new type
    if (!data.path || mountsByType.includes(data.path)) {
      data.path = backendType;
    }
  }

  typeChangeSideEffect(type: string) {
    // If type PKI, set max lease to ~10years
    if (this.args.mountCategory === 'secret') {
      this.args.mountModel.data.config.max_lease_ttl = type === 'pki' ? '3650d' : 0;
    }
  }

  checkModelWarnings() {
    // check for warnings on change
    // since we only show errors on submit we need to clear those out and only send warning state
    const { mountModel } = this.args;
    const { state } = mountModel.toJSON();
    for (const key in state) {
      if (state[key]) {
        state[key].errors = [];
      }
    }
    this.modelValidations = state;
    this.invalidFormAlert = null;
  }

  async saveKvConfig(path: string, formData: SecretsEngineForm['data']) {
    const { options, kv_config = {} } = formData;
    const { max_versions, cas_required, delete_version_after } = kv_config;
    const isKvV2 = options?.version === 2 && ['kv', 'generic'].includes(this.args.mountModel.normalizedType);
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

  async onMountError(status: number, errors: ApiError[], message: string) {
    if (status === 403) {
      this.flashMessages.danger(
        'You do not have access to the sys/mounts endpoint. The secret engine was not mounted.'
      );
    } else if (errors) {
      this.errorMessage = errors.map((e) => {
        if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
        return e;
      });
    } else if (message) {
      this.errorMessage = message;
    } else {
      this.errorMessage = 'An error occurred, check the vault logs.';
    }
  }

  @task
  @waitFor
  *mountBackend(event: Event) {
    event.preventDefault();
    const { mountModel, mountCategory } = this.args;
    const { type } = mountModel;
    const { path } = mountModel.data;
    // only submit form if validations pass
    const { isValid, state, invalidFormMessage, data } = mountModel.toJSON();
    if (!isValid) {
      this.modelValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return;
    }

    try {
      if (mountCategory === 'secret') {
        yield this.api.sys.mountsEnableSecretsEngine(path, data);
        yield this.saveKvConfig(path, data as SecretsEngineForm['data']);
      } else {
        yield this.api.sys.authEnableMethod(path, data);
      }
      this.flashMessages.success(
        `Successfully mounted the ${mountModel.type} ${
          mountCategory === 'secret' ? 'secrets engine' : 'auth method'
        } at ${path}.`
      );
      // check whether to use the Ember engine route
      const version = (data as SecretsEngineForm['data']).options?.version;
      const useEngineRoute = isAddonEngine(mountModel.normalizedType, Number(version));
      this.args.onMountSuccess(type, path, useEngineRoute);
    } catch (error) {
      const { status, response, message } = yield this.api.parseError(error);
      this.onMountError(status, response.errors, message);
    }
  }

  @action
  onKeyUp(name: string, value: string) {
    set(this.args.mountModel.data, name, value);
    this.checkModelWarnings();
  }

  @action
  setMountType(value: string) {
    this.args.mountModel.type = value;
    this.typeChangeSideEffect(value);
    this.checkPathChange(value);
  }

  @action
  handleIdentityTokenKeyChange(value: string[] | string): void {
    // if array, it's coming from the search-select component, otherwise it hit the fallback component and will come in as a string.
    const { config } = this.args.mountModel.data;
    config.identity_token_key = Array.isArray(value) ? value[0] : value;
  }
}
