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
  @tracked selectedEngineType = null;
  @tracked selectedEngineName = null;

  // ARG TODO cleaner way to write this??
  filterByEngineType(getSupportedBackends = true) {
    // 1. Need to confirm if the selected engineType is a supported engine type or unsupported
    const isSupported = LINKED_BACKENDS.includes(this.selectedEngineType);
    if (isSupported) {
      // 2. if the filtered engineType is supported return a value only if the call is coming from get supportedBackends
      return getSupportedBackends
        ? this.displayableBackends
            .filter((backend) => this.selectedEngineType === backend.get('engineType'))
            .sortBy('id')
        : [];
    } else {
      // 3. if the filtered engineType is unsupported return nothing if the call is coming from the get supportedBackend but return if call is from  get unsupportedBackends
      return getSupportedBackends
        ? []
        : this.displayableBackends
            // exact same as supportedBackends but negated
            .filter((backend) => this.selectedEngineType === backend.get('engineType'))
            .sortBy('id');
    }
  }

  get supportedBackends() {
    // if list has been filtered by engineType
    if (this.selectedEngineType) {
      return this.filterByEngineType();
    }
    return (this.displayableBackends || [])
      .filter((backend) => LINKED_BACKENDS.includes(backend.get('engineType')))
      .sortBy('id');
  }

  get unsupportedBackends() {
    // if list has been filtered by engineType
    if (this.selectedEngineType) {
      return this.filterByEngineType(false);
    }
    return (
      (this.displayableBackends || [])
        // exact same as supportedBackends but negated
        .filter((backend) => !LINKED_BACKENDS.includes(backend.get('engineType')))
        .sortBy('id')
    );
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

  get secretEngineArray() {
    // if list has not been filtered by name or engineType
    return this.displayableBackends.map((modelObject) => ({
      name: modelObject.engineType,
      id: modelObject.engineType,
    }));
  }

  @action
  selectEngineType([engine]) {
    this.selectedEngineType = engine;
    // filter the list
  }
}
