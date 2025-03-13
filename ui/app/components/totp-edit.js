/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

/**
 * @module TotpEdit
 * `TotpEdit` is a component that allows you to create, view or delete a TOTP key or view the QR code of the key.
 *
 * @example
 * ```js
 *   <TotpEdit @model={{this.model}} @mode={{this.mode}} />
 * ```
 * @param {object} model - The totp edit model.
 * @param {string} mode - The mode to render. Either 'create' or 'show'.
 */
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class TotpEdit extends Component {
  @service router;

  @tracked hasGenerated = false;
  @tracked invalidFormAlert = '';
  @tracked modelValidations;

  successCallback;

  get keyFormFields() {
    const shared = ['name', 'generate', 'issuer', 'accountName'];
    const generated = [...shared, 'exported'];
    const nonGenerated = [...shared, 'url', 'key'];
    return this.args.model.generate ? generated : nonGenerated;
  }

  get mode() {
    return this.args.mode || 'show';
  }

  get model() {
    return this.args.model;
  }

  persist(method, successCallback) {
    // TODO refactor this further
    return this.model[method]().then(() => {
      if (!this.model.isError) {
        if (this.model.backend === 'totp' && this.model.generate) {
          this.hasGenerated = true;
          this.successCallback = successCallback;
        } else {
          successCallback(this.model);
        }
      }
    });
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  reset() {
    this.model.unloadRecord();
    this.successCallback(null);
  }

  @action
  delete() {
    this.persist('destroyRecord', () => {
      this.transitionToRoute(LIST_ROOT_ROUTE);
    });
  }

  @action
  create(event) {
    event.preventDefault();
    const { isValid, state, invalidFormMessage } = this.model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = invalidFormMessage;
    if (isValid) {
      const modelId = this.model.name;

      // TODO verify url resolves for a key before submitting -> confusing error message
      this.persist('save', () => {
        this.transitionToRoute(SHOW_ROUTE, modelId);
      });
    }
  }
}
