/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module AuthConfigForm/Config
 * The `AuthConfigForm/Config` is the base form to configure auth methods.
 *
 * @example
 * <AuthConfigForm::Config @model={{this.model}} />
 *
 * @property model=null {DS.Model} - The corresponding auth model that is being configured.
 *
 */

export default class AuthConfigBase extends Component {
  @service flashMessages;
  @service router;

  @task
  @waitFor
  *saveModel(evt) {
    evt.preventDefault();
    try {
      yield this.args.model.save();
    } catch (err) {
      // AdapterErrors are handled by the error-message component
      // in the form
      if (err instanceof AdapterError === false) {
        throw err;
      }
      return;
    }
    this.router.transitionTo('vault.cluster.access.methods').followRedirects();
    this.flashMessages.success('The configuration was saved successfully.');
  }
}
