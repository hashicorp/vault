/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

/**
 * @module KvSecretEdit is used for creating a new version of a secret
 *
 * <Page::Secret::Edit
 *  @secret={{this.model.newVersion}}
 *  @previousVersion={{this.model.secret.version}}
 *  @currentVersion={{this.model.metadata.currentVersion}}
 *  @breadcrumbs={{this.breadcrumbs}
 * />
 *
 * @param {model} secret - Ember data model: 'kv/data', the new record for the new secret version saved by the form
 * @param {number} previousVersion - previous secret version number
 * @param {number} currentVersion - current secret version, comes from the metadata endpoint
 * @param {array} breadcrumbs - breadcrumb objects to render in page header
 */

/* eslint-disable no-undef */
export default class KvSecretEdit extends Component {
  @service controlGroup;
  @service flashMessages;
  @service router;

  @tracked showJsonView = false;
  @tracked showDiff = false;
  @tracked errorMessage;
  @tracked modelValidations;
  @tracked invalidFormAlert;
  originalSecret;

  constructor() {
    super(...arguments);
    this.originalSecret = JSON.stringify(this.args.secret.secretData || {});
  }

  get showOldVersionAlert() {
    const { currentVersion, previousVersion } = this.args;
    // isNew check prevents alert from flashing after save but before route transitions
    if (!currentVersion || !previousVersion || !this.args.secret.isNew) return false;
    if (currentVersion !== previousVersion) return true;
    return false;
  }

  get diffDelta() {
    const oldData = JSON.parse(this.originalSecret);
    const newData = this.args.secret.secretData;

    const diffpatcher = jsondiffpatch.create({});
    return diffpatcher.diff(oldData, newData);
  }

  get visualDiff() {
    if (!this.showDiff) return null;
    const newData = this.args.secret.secretData;
    return this.diffDelta
      ? jsondiffpatch.formatters.html.format(this.diffDelta, newData)
      : JSON.stringify(newData, undefined, 2);
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.secret.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { secret } = this.args;
        yield secret.save();
        this.flashMessages.success(`Successfully created new version of ${secret.path}.`);
        this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details', {
          queryParams: { version: secret?.version },
        });
      }
    } catch (error) {
      let message = errorMessage(error);
      if (error.message === 'Control Group encountered') {
        this.controlGroup.saveTokenFromError(error);
        const err = this.controlGroup.logFromError(error);
        message = err.content;
      }
      this.errorMessage = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  onCancel() {
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details');
  }
}
