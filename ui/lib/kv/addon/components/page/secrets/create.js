/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import { pathIsFromDirectory } from 'vault/lib/kv-breadcrumbs';

/**
 * @module KvSecretCreate is used for creating the initial version of a secret
 *
 * <Page::Secrets::Create
 *    @secret={{this.model.secret}}
 *    @metadata={{this.model.metadata}}
 *    @breadcrumbs={{this.breadcrumbs}}
 *  />
 *
 * @param {model} secret - Ember data model: 'kv/data', the new record saved by the form
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {array} breadcrumbs - breadcrumb objects to render in page header
 */

export default class KvSecretCreate extends Component {
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

  @action
  pathValidations() {
    // check path attribute warnings on key up
    const { state } = this.args.secret.validate();
    if (state?.path?.warnings) {
      // only set model validations if warnings exist
      this.modelValidations = state;
    }
  }

  validate() {
    const dataValidations = this.args.secret.validate();
    const metadataValidations = this.args.metadata.validate();
    const state = { ...dataValidations.state, ...metadataValidations.state };
    const failed = !dataValidations.isValid || !metadataValidations.isValid;
    return { state, isValid: !failed };
  }

  @task
  *save(event) {
    event.preventDefault();
    this.flashMessages.clearMessages();
    try {
      const { isValid, state } = this.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = isValid ? '' : 'There was an error submitting this form.';
      if (isValid) {
        const { secret } = this.args;
        yield this.args.secret.save();
        this.flashMessages.success(`Successfully created secret ${secret.path}.`);
        this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', secret.path);
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorMessage = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  onCancel() {
    pathIsFromDirectory(this.args.secret?.path)
      ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', this.args.secret.path)
      : this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
  }
}
