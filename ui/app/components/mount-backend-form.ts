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
import { methods } from 'vault/helpers/mountable-auth-methods';
import { isAddonEngine, allEngines } from 'vault/helpers/mountable-secret-engines';
import { ResponseError } from '@hashicorp/vault-client-typescript';

import type FlashMessageService from 'vault/services/flash-messages';
import type Store from '@ember-data/store';
import type AdapterError from '@ember-data/adapter/error';
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
 *   <MountBackendForm @mountType="secret" @onMountSuccess={{this.onMountSuccess}} />```
 *
 * @param {function} onMountSuccess - A function that transitions once the Mount has been successfully posted.
 * @param {string} [mountType=auth] - The type of backend we want to mount.
 *
 */

type MountModel = SecretsEngineForm | AuthEnableModel;

interface Args {
  mountModel: MountModel;
  mountType: 'secret' | 'auth';
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

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    if (noTeardown && this.args.mountType === 'auth' && this.args?.mountModel?.isNew) {
      this.args.mountModel.unloadRecord();
    }
    super.willDestroy();
  }

  checkPathChange(type: string) {
    if (!type) return;
    const mount = this.args.mountModel;
    const currentPath = mount.path;
    const mountTypes =
      this.args.mountType === 'secret'
        ? allEngines().map((engine) => engine.type)
        : methods().map((auth) => auth.type);
    // if the current path has not been altered by user,
    // change it here to match the new type
    if (!currentPath || mountTypes.includes(currentPath)) {
      mount.path = type;
    }
  }

  typeChangeSideEffect(type: string) {
    if (this.args.mountType === 'secret') {
      // If type PKI, set max lease to ~10years
      this.args.mountModel.config.maxLeaseTtl = type === 'pki' ? '3650d' : 0;
    }
  }

  checkModelValidity(model: MountModel) {
    const { mountType } = this.args;
    const { isValid, state, invalidFormMessage, data } =
      mountType === 'secret' ? model.toJSON() : model.validate();
    this.modelValidations = state;
    this.invalidFormAlert = invalidFormMessage;
    return { isValid, data };
  }

  checkModelWarnings() {
    // check for warnings on change
    // since we only show errors on submit we need to clear those out and only send warning state
    const { mountType, mountModel } = this.args;
    const { state } = mountType === 'secret' ? mountModel.toJSON() : mountModel.validate();
    for (const key in state) {
      state[key].errors = [];
    }
    this.modelValidations = state;
    this.invalidFormAlert = null;
  }

  async saveKvConfig(path: string, formData: SecretsEngineForm['data']) {
    const { options, kvConfig = {} } = formData;
    const { maxVersions, casRequired, deleteVersionAfter } = kvConfig;
    const isKV = options?.version === 2 && ['kv', 'generic'].includes(this.args.mountModel.engineType);
    const hasConfig = maxVersions || casRequired || deleteVersionAfter;

    if (isKV && hasConfig) {
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
    const { mountModel, mountType } = this.args;
    const { type, path } = mountModel;
    // only submit form if validations pass
    const { isValid, data: formData } = this.checkModelValidity(mountModel);
    if (!isValid) {
      return;
    }

    try {
      if (mountType === 'secret') {
        yield this.api.sys.mountsEnableSecretsEngine(path, formData);
        yield this.saveKvConfig(path, formData);
      } else {
        yield mountModel.save();
      }
      this.flashMessages.success(
        `Successfully mounted the ${type} ${
          this.args.mountType === 'secret' ? 'secrets engine' : 'auth method'
        } at ${path}.`
      );
      // Check whether to use the engine route, since KV version 1 does not
      const useEngineRoute = isAddonEngine(mountModel.engineType, Number(formData.options?.version));
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
