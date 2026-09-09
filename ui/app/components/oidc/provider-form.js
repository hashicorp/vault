/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module OidcProviderForm
 * OidcProviderForm components are used to create and update OIDC providers
 *
 * @example
 * ```js
 * <OidcProviderForm @form={{this.model.form}} @scopes={{this.model.scopes}} />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} form - oidc provider form
 * @param {Array} scopes - array of scopes to populate dropdown in form
 * @param {Array} clients - array of clients to populate dropdown in form
 * @param {onCancel} onCancel - callback triggered when cancel button is clicked
 * @param {onSave} onSave - callback triggered on save success
 */

export default class OidcProviderForm extends Component {
  @service api;
  @service flashMessages;

  @tracked modelValidations;
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked radioCardGroupValue = 'limited';
  @tracked selectedClients = [];

  constructor() {
    super(...arguments);
    // If "*" is provided, all clients are allowed: https://developer.hashicorp.com/vault/api-docs/secret/identity/oidc-provider#parameters
    const { allowed_client_ids } = this.args.form.data;
    if (!allowed_client_ids || allowed_client_ids.includes('*')) {
      this.radioCardGroupValue = 'allow_all';
    }
    // initialize selectedClients for SearchSelect component with allowed_client_ids from form data
    this.updateSelectedClients();
  }

  get breadcrumbs() {
    const crumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'OIDC provider: Providers', route: 'vault.cluster.access.oidc.providers' },
    ];

    if (!this.args.form.isNew) {
      crumbs.push({
        label: this.args.form.data.name,
        route: 'vault.cluster.access.oidc.providers.provider.details',
        model: this.args.form.data.name,
      });
    }

    crumbs.push({ label: this.args.form.isNew ? 'Create provider' : 'Edit provider' });
    return crumbs;
  }

  // function passed to search select
  renderTooltip(selection, dropdownOptions) {
    // if a client has been deleted it will not exist in dropdownOptions (response from search select's query)
    const clientExists = !!dropdownOptions.find((opt) => opt.client_id === selection);
    return !clientExists ? 'The application associated with this client_id no longer exists' : false;
  }

  updateSelectedClients() {
    const { data } = this.args.form;
    this.selectedClients = data.allowed_client_ids?.map((clientId) =>
      this.args.clients.find((client) => client.client_id === clientId)
    );
  }

  @action
  handleClientSelection(selection) {
    const { data } = this.args.form;
    // when triggered from search-select component an array is passed
    // set selection as clients
    if (Array.isArray(selection)) {
      data.allowed_client_ids = selection.map((client) => client.client_id);
    } else {
      // otherwise update radio button value and reset clients so
      // UI always reflects a user's selection (including when no clients are selected)
      this.radioCardGroupValue = selection;
      data.allowed_client_ids = [];
    }
    // update selectedClients which appear in SearchSelect
    this.updateSelectedClients();
  }

  save = task(
    waitFor(async (event) => {
      event.preventDefault();
      try {
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;

        if (isValid) {
          if (this.radioCardGroupValue === 'allow_all') {
            data.allowed_client_ids = ['*'];
          }
          const { name, ...payload } = data;
          await this.api.identity.oidcWriteProvider(name, payload);
          this.flashMessages.success(
            `Successfully ${this.args.form.isNew ? 'created' : 'updated'} the OIDC provider ${name}.`
          );
          this.args.onSave();
        }
      } catch (error) {
        const { message } = await this.api.parseError(error);
        this.errorBanner = message;
        this.invalidFormAlert = 'There was an error submitting this form.';
      }
    })
  );
}
