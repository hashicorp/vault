/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ValidationMap } from 'vault/vault/app-types';
import errorMessage from 'vault/utils/error-message';

import type CaConfigModel from 'vault/models/ssh/ca-config';
import type Router from '@ember/routing/router';
import type Store from '@ember-data/store';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module ConfigureSshComponent is used to configure the SSH secret engine.
 *
 * @example
 * ```js
 * <SecretEngine::ConfigureSsh
 *    @model={{this.model.ssh-ca-config}}
 *    @id={{this.model.id}}
 *  />
 * ```
 *
 * @param {string} model - SSH ca-config model
 * @param {string} id - name of the SSH secret engine, ex: 'ssh-123'
 */

interface Args {
  model: CaConfigModel;
  id: string;
}

export default class ConfigureSshComponent extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked errorMessage: string | null = null;
  @tracked invalidFormAlert: string | null = null;
  @tracked modelValidations: ValidationMap | null = null;

  @task
  *save(event: Event) {
    event.preventDefault();
    this.resetErrors();
    const { id, model } = this.args;
    const isValid = this.validate(model);

    if (!isValid) {
      this.flashMessages.danger('Please correct the errors in the form before submitting.');
      return;
    }
    // Check if any of the models attributes have changed.
    // If no changes to the model, transition and notify user.
    // Otherwise, save the model.
    const attributesChanged = Object.keys(model.changedAttributes()).length > 0;
    if (!attributesChanged) {
      this.flashMessages.danger('No changes detected.');
      this.transition(id);
    }

    try {
      yield model.save();
      this.flashMessages.success(`Successfully saved ${id}'s root configuration.`);
    } catch (error) {
      this.errorMessage = errorMessage(error);
    }

    this.transition(id);
  }

  transition(id: string) {
    // prevent transition if there are errors with root configuration
    if (this.errorMessage) {
      this.invalidFormAlert = 'There was an error submitting this form.';
    } else {
      this.router.transitionTo('vault.cluster.secrets.backend.configuration', id);
    }
  }

  validate(model: CaConfigModel) {
    const { isValid, state, invalidFormMessage } = model.validate();
    this.modelValidations = isValid ? null : state;
    this.invalidFormAlert = isValid ? '' : invalidFormMessage;
    return isValid;
  }

  resetErrors() {
    this.flashMessages.clearMessages();
    this.errorMessage = null;
    this.invalidFormAlert = null;
  }

  @action
  onCancel() {
    // clear errors because they're canceling out of the workflow.
    this.resetErrors();
    this.transition(this.args.id);
  }
}
