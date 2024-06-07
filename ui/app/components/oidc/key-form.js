/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

/**
 * @module OidcKeyForm
 * OidcKeyForm components are used to create and update OIDC providers
 *
 * @example
 * ```js
 * <OidcKeyForm @model={{this.model}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} model - oidc client model
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcKeyForm extends Component {
  @service store;
  @service flashMessages;
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked modelValidations;
  @tracked radioCardGroupValue =
    // If "*" is provided, all clients are allowed: https://developer.hashicorp.com/vault/api-docs/secret/identity/oidc-provider#parameters
    !this.args.model.allowedClientIds || this.args.model.allowedClientIds.includes('*')
      ? 'allow_all'
      : 'limited';

  get filterDropdownOptions() {
    // query object sent to search-select so only clients that reference this key appear in dropdown
    return { paramKey: 'key', filterFor: [this.args.model.name] };
  }

  @action
  handleClientSelection(selection) {
    // if array then coming from search-select component, set selection as model clients
    if (Array.isArray(selection)) {
      this.args.model.allowedClientIds = selection.map((client) => client.clientId);
    } else {
      // otherwise update radio button value and reset clients so
      // UI always reflects a user's selection (including when no clients are selected)
      this.radioCardGroupValue = selection;
      this.args.model.allowedClientIds = [];
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { isNew, name } = this.args.model;
        if (this.radioCardGroupValue === 'allow_all') {
          this.args.model.allowedClientIds = ['*'];
        }
        // if TTL components are toggled off, set to default lease duration
        const { rotationPeriod, verificationTtl } = this.args.model;
        // value returned from API is a number, and string when from form action
        if (Number(rotationPeriod) === 0) this.args.model.rotationPeriod = '24h';
        if (Number(verificationTtl) === 0) this.args.model.verificationTtl = '24h';
        yield this.args.model.save();
        this.flashMessages.success(
          `Successfully ${isNew ? 'created' : 'updated'} the key
          ${name}.`
        );
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
