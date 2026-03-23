/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action, set } from '@ember/object';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { filterEnginesByMountCategory } from 'core/utils/all-engines-metadata';
import { task } from 'ember-concurrency';
import { MOUNT_CATEGORIES } from 'vault/utils/plugin-catalog-helpers';

import type { ApiError } from '@ember-data/adapter/error';
import type Store from '@ember-data/store';
import type AuthMethodForm from 'vault/forms/auth/method';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type FlashMessageService from 'vault/services/flash-messages';
import type { ValidationMap } from 'vault/vault/app-types';

/**
 * @module MountBackendForm
 * The `MountBackendForm` is used to mount authentication methods.
 *
 * @example ```js
 *   <MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />```
 *
 * @param {AuthMethodForm} mountModel - The authentication method form.
 * @param {function} onMountSuccess - A function that transitions once the Mount has been successfully posted.
 *
 */

interface TypeInfo {
  type: string;
  displayName: string;
  glyph: string;
}
interface Args {
  mountModel: AuthMethodForm;
  mountTypes: TypeInfo[] | undefined;
  onMountSuccess: (type: string, path: string, useEngineRoute: boolean) => void;
}

const AUTH_MOUNT_CATEGORY = MOUNT_CATEGORIES.AUTH;

export default class MountBackendForm extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly api: ApiService;

  // validation related properties
  @tracked modelValidations: ValidationMap | null = null;
  @tracked invalidFormAlert: string | null = null;

  @tracked errorMessage: string | string[] = '';

  get mountForm(): AuthMethodForm {
    return this.args.mountModel;
  }

  get showEnable(): boolean {
    return !!this.mountForm.type;
  }

  get title(): string {
    const typeInfo = (this.args?.mountTypes || []).find(
      (mountType: TypeInfo) => this.mountForm.type === mountType.type
    );
    if (this.showEnable) {
      return `Enable ${typeInfo?.displayName || this.mountForm.type} Authentication Method`;
    }
    return 'Enable an Authentication Method';
  }

  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }

  checkPathChange(backendType: string) {
    if (!backendType) return;
    const { data } = this.mountForm;
    // Always use auth mount category since this component only handles auth methods
    const mountsByType = filterEnginesByMountCategory({
      mountCategory: AUTH_MOUNT_CATEGORY,
      isEnterprise: true,
    }).map((engine) => engine.type);

    // if the current path has not been altered by user (is empty or matches a default mount type),
    // change it here to match the new type
    if (!data.path || mountsByType.includes(data.path)) {
      data.path = backendType;
    }
  }

  checkModelWarnings() {
    // check for warnings on change
    // since we only show errors on submit we need to clear those out and only send warning state
    const mountModel = this.mountForm;
    const { state } = mountModel.toJSON();
    for (const key in state) {
      if (state[key]) {
        state[key].errors = [];
      }
    }
    this.modelValidations = state;
    this.invalidFormAlert = null;
  }

  async onMountError(status: number, errors: ApiError[], message: string) {
    if (status === 403) {
      this.flashMessages.danger(
        'You do not have access to the sys/auth endpoint. The auth method was not mounted.'
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
    const mountModel = this.mountForm;
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
      yield this.api.sys.authEnableMethod(path, data);
      this.flashMessages.success(`Successfully mounted the ${mountModel.type} auth method at ${path}.`);
      this.args.onMountSuccess(type, path, false);
    } catch (error) {
      const { status, response, message } = yield this.api.parseError(error);
      this.onMountError(status, response.errors, message);
    }
  }

  @action
  onKeyUp(name: string, value: string) {
    set(this.mountForm.data, name, value);
    this.checkModelWarnings();
  }

  @action
  setMountType(value: string) {
    this.mountForm.type = value;
    this.checkPathChange(value);
  }

  @action
  handleIdentityTokenKeyChange(value: string[] | string): void {
    // if array, it's coming from the search-select component, otherwise it hit the fallback component and will come in as a string.
    const { config } = this.mountForm.data;
    config.identity_token_key = Array.isArray(value) ? value[0] : value;
  }
}
