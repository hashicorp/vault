/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module TotpEdit
 * `TotpEdit` is a component that allows you to create, view or delete a TOTP key.
 * When creating a key if `generate` and `exported` are true then after a successful save the UI renders a QR code for the generated key.
 * @example
 *   <TotpEdit @form={{this.form}} @mode={{this.mode}} @capabilities={{this.capabilities}} />
 *
 * @param {object} form - The TotpKeyForm instance.
 * @param {string} mode - The mode to render. Either 'create' or 'show'.
 * @param {object} capabilities - Capabilities object with canDelete, canRead flags.
 */
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class TotpEdit extends Component {
  @service router;
  @service flashMessages;
  @service api;

  @tracked hasGenerated = false;
  @tracked invalidFormAlert = '';
  @tracked modelValidations;
  @tracked key;

  constructor(owner, args) {
    super(owner, args);
    // In show mode, the route fetches the key and passes it via @form
    // This stores it in our tracked property for consistent data source
    if (args.mode === 'show') {
      this.key = args.form.data;
    }
  }

  displayFields = [
    { field: 'account_name', label: 'Account name' },
    { field: 'algorithm', label: 'Algorithm' },
    { field: 'digits', label: 'Digits' },
    { field: 'issuer', label: 'Issuer' },
    { field: 'period', label: 'Period' },
  ];

  generatedFields = [{ field: 'url', label: 'URL' }];

  breadcrumbs = [
    { label: 'Vault', text: 'Vault', icon: 'vault', path: 'vault.cluster.dashboard' },
    { text: 'Secrets engines', path: 'vault.cluster.secrets.backends' },
    this.args.root,
    { label: this.title, text: this.title },
  ];

  get title() {
    if (this.args.mode === 'create') {
      return 'Create a TOTP key';
    }
    return 'TOTP key';
  }

  get subtitle() {
    if (this.args.mode === 'create') return '';
    return this.args.form.data.name;
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  reset() {
    const { name } = this.args.form.data;
    this.transitionToRoute(SHOW_ROUTE, name);
  }

  @action
  async deleteKey() {
    try {
      const { name, backend } = this.args.form.data;
      await this.api.secrets.totpDeleteKey(name, backend);
      this.transitionToRoute(LIST_ROOT_ROUTE);
      this.flashMessages.success(`${name} was successfully deleted.`);
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(message);
    }
  }

  createKey = task(
    waitFor(async (event) => {
      event.preventDefault();
      const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;

      if (!isValid) return;
      try {
        const { name, backend, generate, exported } = this.args.form.data;
        const resp = await this.api.secrets.totpCreateKey(name, backend, data);

        if (generate && exported) {
          // stay in this template and show QR code returned from response
          if (resp?.data) {
            this.key = resp.data;
          }
          this.hasGenerated = true;
        } else {
          this.transitionToRoute(SHOW_ROUTE, name);
        }
        this.flashMessages.success('Successfully created key.');
      } catch (err) {
        const { message } = await this.api.parseError(err);
        this.flashMessages.danger(message);
      }
    })
  );
}
