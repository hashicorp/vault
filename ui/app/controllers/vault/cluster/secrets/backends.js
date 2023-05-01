/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
/* eslint ember/no-computed-properties-in-native-classes: 'warn' */
import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { filterBy } from '@ember/object/computed';
import { dropTask, task } from 'ember-concurrency';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';

const LINKED_BACKENDS = supportedSecretBackends();

export default class VaultClusterSecretsBackendController extends Controller {
  @service flashMessages;
  @filterBy('model', 'shouldIncludeInList') displayableBackends;

  @tracked secretEngineOptions = [];

  get supportedBackends() {
    return (this.displayableBackends || [])
      .filter((backend) => LINKED_BACKENDS.includes(backend.get('engineType')))
      .sortBy('id');
  }

  get unsupportedBackends() {
    return (this.displayableBackends || []).slice().removeObjects(this.supportedBackends).sortBy('id');
  }

  @task
  @dropTask
  *disableEngine(engine) {
    const { engineType, path } = engine;
    try {
      yield engine.destroyRecord();
      this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling the ${engineType} Secrets Engine at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }

  @action
  selectEngineType([engine]) {
    this.engineType = engine;
    if (!engine) {
      this.secretEngineOptions = [];
      // // on clear, also make sure auth method is cleared
      // this.selectedAuthMethod = null;
    } else {
      // Side effect: set auth namespaces
      // const mounts = this.filteredActivityByNamespace.mounts?.map((mount) => ({
      //   id: mount.label,
      //   name: mount.label,
      // }));
      this.secretEngineOptions = ['Test', 'test'];
    }
  }
}
