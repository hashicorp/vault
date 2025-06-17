/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { presence } from 'vault/utils/forms/validators';
import { filterEnginesByMountCategory, isAddonEngine } from 'vault/utils/all-engines-metadata';
import { assert } from '@ember/debug';
import { ResponseError } from '@hashicorp/vault-client-typescript';
import AdapterError from '@ember-data/adapter/error';

import type FlashMessageService from 'vault/services/flash-messages';
import type Store from '@ember-data/store';
import type { AuthEnableModel } from 'vault/routes/vault/cluster/settings/auth/enable';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type CapabilitiesService from 'vault/services/capabilities';
import type ApiService from 'vault/services/api';
import type { ApiError } from '@ember-data/adapter/error';

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

type MountModel = SecretsEngineForm | AuthEnableModel;

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
  @tracked modelValidations = null;
  @tracked invalidFormAlert = null;

  @tracked errorMessage: string | string[] = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    assert(`@mountCategory is required. Must be "auth" or "secret".`, presence(this.args.mountCategory));
  }

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    if (noTeardown && this.args.mountCategory === 'auth' && this.args?.mountModel?.isNew) {
      this.args.mountModel.unloadRecord();
    }
    super.willDestroy();
  }

  checkPathChange(backendType: string) {
    if (!backendType) return;
    const mount = this.args.mountModel;
    const currentPath = mount.path;
    // mountCategory is usually 'secret' or 'auth', but sometimes an empty string is passed in (like when we click the cancel button).
    // In these cases, we should default to returning auth methods.
    const mountsByType = filterEnginesByMountCategory({
      mountCategory: this.args.mountCategory ?? 'auth',
      isEnterprise: true,
    }).map((engine) => engine.type);
    // if the current path has not been altered by user,
    // change it here to match the new type
    if (!currentPath || mountsByType.includes(currentPath)) {
      mount.path = backendType;
    }
  }

  typeChangeSideEffect(type: string) {
    if (this.args.mountCategory !== 'secret') return;
    // If type PKI, set max lease to ~10years
    this.args.mountModel.config.maxLeaseTtl = type === 'pki' ? '3650d' : 0;
  }

  checkModelValidity(model: MountModel) {
    const { mountCategory } = this.args;
    const { isValid, state, invalidFormMessage, data } =
      mountCategory === 'secret' ? model.toJSON() : model.validate();
    this.modelValidations = state;
    this.invalidFormAlert = invalidFormMessage;
    return { isValid, data };
  }

  checkModelWarnings() {
    // check for warnings on change
    // since we only show errors on submit we need to clear those out and only send warning state
    const { mountCategory, mountModel } = this.args;
    const { state } = mountCategory === 'secret' ? mountModel.toJSON() : mountModel.validate();
    for (const key in state) {
      state[key].errors = [];
    }
    this.modelValidations = state;
    this.invalidFormAlert = null;
  }

  async saveKvConfig(path: string, formData: SecretsEngineForm['data']) {
    const { options, kvConfig = {} } = formData;
    const { maxVersions, casRequired, deleteVersionAfter } = kvConfig;
    const isKvV2 = options?.version === 2 && ['kv', 'generic'].includes(this.args.mountModel.engineType);
    const hasConfig = maxVersions || casRequired || deleteVersionAfter;

    if (isKvV2 && hasConfig) {
      try {
        const { canUpdate } = await this.capabilities.for('kvConfig', { path });
        if (canUpdate) {
          await this.api.secrets.kvV2Configure(path, kvConfig);
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
    const { type, path } = mountModel;
    // only submit form if validations pass
    const { isValid, data: formData } = this.checkModelValidity(mountModel);
    if (!isValid) {
      return;
    }

    try {
      if (mountCategory === 'secret') {
        yield this.api.sys.mountsEnableSecretsEngine(path, formData);
        yield this.saveKvConfig(path, formData);
      } else {
        yield mountModel.save();
      }
      this.flashMessages.success(
        `Successfully mounted the ${type} ${
          this.args.mountCategory === 'secret' ? 'secrets engine' : 'auth method'
        } at ${path}.`
      );
      // check whether to use the Ember engine route
      const useEngineRoute = isAddonEngine(mountModel.engineType, Number(formData?.options?.version));
      this.args.onMountSuccess(type, path, useEngineRoute);
    } catch (error) {
      if (error instanceof ResponseError) {
        const { status, response, message } = yield this.api.parseError(error);
        this.onMountError(status, response.errors, message);
      } else {
        const err = error as AdapterError;
        this.onMountError(err.httpStatus, err.errors, err.message);
      }
    }
  }

  @action
  onKeyUp(name: string, value: string) {
    this.args.mountModel[name] = value;
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
    this.args.mountModel.config.identityTokenKey = Array.isArray(value) ? value[0] : value;
  }
}
