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

import type FlashMessageService from 'vault/services/flash-messages';
import type Store from '@ember-data/store';
import type AdapterError from '@ember-data/adapter/error';

import type { AuthEnableModel } from 'vault/routes/vault/cluster/settings/auth/enable';
import type { MountSecretBackendModel } from 'vault/routes/vault/cluster/settings/mount-secret-backend';

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

type MountModel = MountSecretBackendModel | AuthEnableModel;

interface Args {
  mountModel: MountModel;
  mountType: 'secret' | 'auth';
  onMountSuccess: (type: string, path: string, useEngineRoute: boolean) => void;
}

export default class MountBackendForm extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;

  // validation related properties
  @tracked modelValidations = null;
  @tracked invalidFormAlert = null;

  @tracked errorMessage: string | string[] = '';

  willDestroy() {
    // components are torn down after store is unloaded and will cause an error if attempt to unload record
    const noTeardown = this.store && !this.store.isDestroying;
    if (noTeardown && this.args?.mountModel?.isNew) {
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
    if (this.args.mountType !== 'secret') return;
    if (type === 'pki') {
      // If type PKI, set max lease to ~10years
      this.args.mountModel.config.maxLeaseTtl = '3650d';
    } else {
      // otherwise reset
      this.args.mountModel.config.maxLeaseTtl = 0;
    }
  }

  checkModelValidity(model: MountModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidations = state;
    this.invalidFormAlert = invalidFormMessage;
    return isValid;
  }

  checkModelWarnings() {
    // check for warnings on change
    // since we only show errors on submit we need to clear those out and only send warning state
    const { state } = this.args.mountModel.validate();
    for (const key in state) {
      state[key].errors = [];
    }
    this.modelValidations = state;
    this.invalidFormAlert = null;
  }

  async showWarningsForKvv2() {
    try {
      const capabilities = await this.store.findRecord('capabilities', `${this.args.mountModel.path}/config`);
      if (!capabilities?.canUpdate) {
        // config error is not thrown from secret-engine adapter, so handling here
        this.flashMessages.warning(
          'You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.'
        );
        // remove the config data from the model otherwise it will persist in the store even though network request failed.
        [
          this.args.mountModel.maxVersions,
          this.args.mountModel.casRequired,
          this.args.mountModel.deleteVersionAfter,
        ] = [0, false, 0];
      }
    } catch (e) {
      // Show different warning if we're not sure the config saved
      this.flashMessages.warning(
        'You may not have access to the config endpoint. The secret engine was mounted, but the configuration settings may not be saved.'
      );
    }
    return;
  }

  @task
  @waitFor
  *mountBackend(event: Event) {
    event.preventDefault();
    const mountModel = this.args.mountModel;
    const { type, path } = mountModel;
    // only submit form if validations pass
    if (!this.checkModelValidity(mountModel)) {
      return;
    }

    const changedAttrKeys = Object.keys(mountModel.changedAttributes());
    const updatesConfig =
      changedAttrKeys.includes('casRequired') ||
      changedAttrKeys.includes('deleteVersionAfter') ||
      changedAttrKeys.includes('maxVersions');

    try {
      yield mountModel.save();
    } catch (error) {
      const err = error as AdapterError;

      if (err.httpStatus === 403) {
        this.flashMessages.danger(
          'You do not have access to the sys/mounts endpoint. The secret engine was not mounted.'
        );
        return;
      }
      if (err.errors) {
        const errors = err.errors.map((e) => {
          if (typeof e === 'object') return e.title || e.message || JSON.stringify(e);
          return e;
        });
        this.errorMessage = errors;
      } else if (err.message) {
        this.errorMessage = err.message;
      } else {
        this.errorMessage = 'An error occurred, check the vault logs.';
      }
      return;
    }
    if (mountModel.isV2KV && updatesConfig) {
      yield this.showWarningsForKvv2();
    }
    this.flashMessages.success(
      `Successfully mounted the ${type} ${
        this.args.mountType === 'secret' ? 'secrets engine' : 'auth method'
      } at ${path}.`
    );
    // Check whether to use the engine route, since KV version 1 does not
    const useEngineRoute = isAddonEngine(mountModel.engineType, mountModel.version);
    yield this.args.onMountSuccess(type, path, useEngineRoute);
    return;
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
