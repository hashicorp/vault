/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';

/**
 * @module KvSecretsCreate renders the form for creating a new secret. 
 * 
 * <Page::Secrets::Create
 *  @secret={{this.model.secret}}
 *  @breadcrumbs={{this.breadcrumbs}}
  /> 
 *
 * @param {model} secret - Ember data model: 'kv/data'  
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvSecretsCreate extends Component {
  @service flashMessages;
  @service router;
  @tracked showJsonView = false;
  @tracked errorMessage;
  @tracked modelValidations;
  @tracked invalidFormAlert;

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.secret.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        yield this.args.secret.save();
        const { path } = this.args.secret;
        this.flashMessages.success(`Successfully created secret ${path}`);
        this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', path);
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorMessage = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
    this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
  }
}
