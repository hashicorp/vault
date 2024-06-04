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
 * @module OidcClientForm
 * OidcClientForm components are used to create and update OIDC clients (a.k.a. applications)
 *
 * @example
 * ```js
 * <OidcClientForm @model={{this.model}} />
 * ```
 * @param {Object} model - oidc client model
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered on save success
 * @param {boolean} [isInline=false] - true when form is rendered within a modal
 */

export default class OidcClientForm extends Component {
  @service flashMessages;
  @tracked modelValidations;
  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked radioCardGroupValue =
    !this.args.model.assignments || this.args.model.assignments.includes('allow_all')
      ? 'allow_all'
      : 'limited';

  get modelAssignments() {
    const { assignments } = this.args.model;
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
      this.args.model.assignments = selection;
    } else {
      // otherwise update radio button value and reset assignments so
      // UI always reflects a user's selection (including when no assignments are selected)
      this.radioCardGroupValue = selection;
      this.args.model.assignments = [];
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
        if (this.radioCardGroupValue === 'allow_all') {
          // the backend permits 'allow_all' AND other assignments, though 'allow_all' will take precedence
          // the UI limits the config by allowing either 'allow_all' OR a list of other assignments
          // note: when editing the UI removes any additional assignments previously configured via CLI
          this.args.model.assignments = ['allow_all'];
        }
        // if TTL components are toggled off, set to default lease duration
        const { idTokenTtl, accessTokenTtl } = this.args.model;
        // value returned from API is a number, and string when from form action
        if (Number(idTokenTtl) === 0) this.args.model.idTokenTtl = '24h';
        if (Number(accessTokenTtl) === 0) this.args.model.accessTokenTtl = '24h';
        const { isNew, name } = this.args.model;
        yield this.args.model.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the application ${name}.`);
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }
}
