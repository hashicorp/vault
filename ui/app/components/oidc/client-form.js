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
 * @module OidcClientForm
 * OidcClientForm components are used to create and update OIDC clients (a.k.a. applications)
 *
 * @example
 * ```js
 * <OidcClientForm @form={{this.model}} />
 * ```
 * @param {Object} form - oidc client form
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered on save success
 */

export default class OidcClientForm extends Component {
  @service flashMessages;
  @service api;

  @tracked modelValidations;
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked radioCardGroupValue = 'limited';

  constructor() {
    super(...arguments);
    const { assignments } = this.args.form.data;
    if (!assignments || assignments.includes('allow_all')) {
      this.radioCardGroupValue = 'allow_all';
    }
  }

  get breadcrumbs() {
    const { isNew } = this.args.form;
    const { name } = this.args.form.data;
    const crumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'OIDC provider: Applications', route: 'vault.cluster.access.oidc.clients' },
    ];

    if (!isNew) {
      crumbs.push({
        label: name,
        route: 'vault.cluster.access.oidc.clients.client.details',
        model: name,
      });
    }

    crumbs.push({ label: isNew ? 'Create application' : 'Edit application' });
    return crumbs;
  }

  get modelAssignments() {
    const { assignments } = this.args.form.data;
    if (assignments.includes('allow_all') && assignments.length === 1) {
      return [];
    } else {
      return assignments;
    }
  }

  @action
  handleAssignmentSelection(selection) {
    // if array then coming from search-select component, set selection as model assignments
    if (Array.isArray(selection)) {
      this.args.form.data.assignments = selection;
    } else {
      // otherwise update radio button value and reset assignments so
      // UI always reflects a user's selection (including when no assignments are selected)
      this.radioCardGroupValue = selection;
      this.args.form.data.assignments = [];
    }
  }

  @action
  onKeyChange([key]) {
    this.args.form.data.key = key;
  }

  save = task(
    waitFor(async (event) => {
      event.preventDefault();
      try {
        const { isNew } = this.args.form;
        const { isValid, state, invalidFormMessage, data } = this.args.form.toJSON();
        this.modelValidations = isValid ? null : state;
        this.invalidFormAlert = invalidFormMessage;

        if (isValid) {
          if (this.radioCardGroupValue === 'allow_all') {
            // the backend permits 'allow_all' AND other assignments, though 'allow_all' will take precedence
            // the UI limits the config by allowing either 'allow_all' OR a list of other assignments
            // note: when editing the UI removes any additional assignments previously configured via CLI
            data.assignments = ['allow_all'];
          }
          // if TTL components are toggled off, set to default lease duration
          const { id_token_ttl, access_token_ttl } = data;
          // value returned from API is a number, and string when from form action
          if (Number(id_token_ttl) === 0) data.id_token_ttl = '24h';
          if (Number(access_token_ttl) === 0) data.access_token_ttl = '24h';

          const { name, ...payload } = data;
          await this.api.identity.oidcWriteClient(name, payload);
          this.flashMessages.success(
            `Successfully ${isNew ? 'created' : 'updated'} the application ${name}.`
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
