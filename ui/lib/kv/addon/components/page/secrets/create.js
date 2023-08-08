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
import errorMessage from 'vault/utils/error-message';

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
  @tracked errors;
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
    this.errors = [];
    const { isValid, state } = this.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : 'There was a problem submitting this form.';

    const { secret, metadata } = this.args;
    // users must have kv/data create to create metadata in the UI
    if (isValid) {
      try {
        // try saving secret data first
        yield secret.save();
        this.flashMessages.success(`Successfully saved secret data for: ${secret.path}.`);
      } catch (error) {
        this.errors.pushObject({ endpoint: 'kv/data', message: errorMessage(error) });
      }

      // only attempt to save metadata secret data was saved
      if (secret.createdTime) {
        try {
          metadata.path = secret.path;
          yield metadata.save();
          this.flashMessages.success(`Successfully saved metadata.`);
        } catch (error) {
          this.flashMessages.danger(`POST kv/metadata: ${errorMessage(error)}`, {
            sticky: true,
          });
        }
      }

      // prevent transition if there are errors with secret data
      if (this.errors.length) {
        this.invalidFormAlert = 'There was an error submitting this form.';
      } else {
        this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', secret.path);
      }
    }
  }

  @action
  onCancel() {
    const { path } = this.args.secret;
    pathIsFromDirectory(path)
      ? this.router.transitionTo('vault.cluster.secrets.backend.kv.list-directory', path)
      : this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
  }
}
