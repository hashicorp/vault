/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import parseURL from 'core/utils/parse-url';
/**
 * @module OidcProviderForm
 * OidcProviderForm components are used to create and update OIDC providers
 *
 * @example
 * ```js
 * <OidcProviderForm @model={{this.model}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} model - oidc client model
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcProviderForm extends Component {
  @service store;
  @service flashMessages;
  @tracked modelValidations;
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked radioCardGroupValue =
    // If "*" is provided, all clients are allowed: https://developer.hashicorp.com/vault/api-docs/secret/identity/oidc-provider#parameters
    !this.args.model.allowedClientIds || this.args.model.allowedClientIds.includes('*')
      ? 'allow_all'
      : 'limited';

  // function passed to search select
  renderInfoTooltip(selection, dropdownOptions) {
    // if a client has been deleted it will not exist in dropdownOptions (response from search select's query)
    const clientExists = !!dropdownOptions.find((opt) => opt.clientId === selection);
    return !clientExists ? 'The application associated with this client_id no longer exists' : false;
  }

  // fired on did-insert from render modifier
  @action
  setIssuer(elem, [model]) {
    model.issuer = model.isNew ? '' : parseURL(model.issuer).origin;
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
        yield this.args.model.save();
        this.flashMessages.success(
          `Successfully ${isNew ? 'created' : 'updated'} the OIDC provider
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
