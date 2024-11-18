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

import type FlashMessageService from 'vault/services/flash-messages';
import type { AuthEnableModel } from 'vault/routes/vault/cluster/settings/auth/enable';
import type ApiService from 'vault/services/api';
import type { AuthEnableMethodRequest } from 'vault/entities/auth-method';

/**
 * @module MountAuthMethodForm
 * The `MountAuthMethodForm` is used to mount an auth backend.
 *
 * @example ```js
 *   <MountAuthMethodForm @onMountSuccess={{this.onMountSuccess}} />```
 *
 * @param {function} onMountSuccess - A function that transitions once the Mount has been successfully posted.
 *
 */

interface Args {
  mountModel: AuthEnableModel;
  onMountSuccess: (type: string, path: string, useEngineRoute: boolean) => void;
}

export default class MountBackendForm extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  // validation related properties
  @tracked modelValidations = null;
  @tracked invalidFormAlert = null;

  @tracked errorMessage = '';

  checkPathChange(type: string) {
    if (!type) return;
    const mount = this.args.mountModel;
    const currentPath = mount.path;
    const mountTypes = methods().map((auth) => auth.type);
    // if the current path has not been altered by user,
    // change it here to match the new type
    if (!currentPath || mountTypes.includes(currentPath)) {
      mount.path = type;
    }
  }

  checkModelValidity(model: AuthEnableModel) {
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

    mountModel.data.type = type;

    const { error } = yield this.api.post('/sys/auth/{path}', {
      params: { path: { path } },
      body: mountModel.data as AuthEnableMethodRequest,
    });

    if (error) {
      if (error.errors) {
        const errors = error.errors.map((e: unknown) => {
          if (e && typeof e === 'object') {
            const error = e as Record<string, unknown>;
            return error['title'] || error['message'] || JSON.stringify(error);
          }
          return e;
        });

        this.errorMessage = errors;
      } else if (error.message) {
        this.errorMessage = error.message;
      } else {
        this.errorMessage = 'An error occurred, check the vault logs.';
      }
    } else {
      this.flashMessages.success(`Successfully mounted the ${type} auth method at ${path}.`);
      yield this.args.onMountSuccess(type, path, false);
    }
  }

  @action
  onKeyUp(name: string, value: string) {
    // eslint-disable-next-line
    // @ts-ignore
    this.args.mountModel.data[name] = value;
    this.checkModelWarnings();
  }

  @action
  setMountType(value: string) {
    this.args.mountModel.type = value;
    this.checkPathChange(value);
  }

  @action
  handleIdentityTokenKeyChange(value: string[] | string): void {
    // if array, it's coming from the search-select component, otherwise it hit the fallback component and will come in as a string.
    // eslint-disable-next-line
    // @ts-ignore
    this.args.mountModel.data.config.identityTokenKey = Array.isArray(value) ? value[0] : value;
  }
}
