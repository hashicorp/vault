/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { filterBy } from '@ember/object/computed';
import { computed } from '@ember/object';
import Controller from '@ember/controller';
import { task } from 'ember-concurrency';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { inject as service } from '@ember/service';
const LINKED_BACKENDS = supportedSecretBackends();

export default Controller.extend({
  flashMessages: service(),
  displayableBackends: filterBy('model', 'shouldIncludeInList'),

  supportedBackends: computed('displayableBackends', 'displayableBackends.[]', function () {
    return (this.displayableBackends || [])
      .filter((backend) => LINKED_BACKENDS.includes(backend.get('engineType')))
      .sortBy('id');
  }),

  unsupportedBackends: computed(
    'displayableBackends',
    'displayableBackends.[]',
    'supportedBackends',
    'supportedBackends.[]',
    function () {
      return (this.displayableBackends || []).slice().removeObjects(this.supportedBackends).sortBy('id');
    }
  ),

  disableEngine: task(function* (engine) {
    const { engineType, path } = engine;
    try {
      yield engine.destroyRecord();
      this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling the ${engineType} Secrets Engine at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }).drop(),
});
