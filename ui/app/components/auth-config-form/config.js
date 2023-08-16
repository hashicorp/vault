/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module AuthConfigForm/Config
 * The `AuthConfigForm/Config` is the base form to configure auth methods.
 *
 * @example
 * ```js
 * {{auth-config-form/config model.model}}
 * ```
 *
 * @property model=null {DS.Model} - The corresponding auth model that is being configured.
 *
 */

const AuthConfigBase = Component.extend({
  tagName: '',
  model: null,

  flashMessages: service(),
  router: service(),
  saveModel: task(
    waitFor(function* () {
      try {
        yield this.model.save();
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
    })
  ),
});

AuthConfigBase.reopenClass({
  positionalParams: ['model'],
});

export default AuthConfigBase;
