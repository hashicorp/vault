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

  get supportedBackends() {
    if (this.selectedEngineType) {
      return this.filterByEngineType();
    }
    if (this.selectedEngineName) {
      return this.filterByEngineName();
    }
    return (this.displayableBackends || [])
      .filter((backend) => LINKED_BACKENDS.includes(backend.get('engineType')))
      .sortBy('id');
  }

  get unsupportedBackends() {
    if (this.selectedEngineType) {
      return this.filterByEngineType(false);
    }
    if (this.selectedEngineType) {
      return this.filterByEngineName(false);
    }
    return (
      (this.displayableBackends || [])
        // exact same as supportedBackends but negated
        .filter((backend) => !LINKED_BACKENDS.includes(backend.get('engineType')))
        .sortBy('id')
    );
  }

  // FOR FILTER LISTS
  get secretEngineArrayByType() {
    return this.displayableBackends.map((modelObject) => ({
      name: modelObject.engineType,
      id: modelObject.engineType,
    }));
  }

  get secretEngineArrayByName() {
    return this.displayableBackends.map((modelObject) => ({
      name: modelObject.id,
      id: modelObject.id,
    }));
  }

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

  filterByEngineName(getSupportedBackends = true) {
    // 1. check if the engine is supported. We only have the name so we need to return the engineType.
    const engineObject = this.displayableBackends.find(
      (modelObject) => modelObject.id === this.selectedEngineName
    );
    const engineType = engineObject.engineType;
    const isSupported = LINKED_BACKENDS.includes(engineObject.engineType);
    if (isSupported) {
      return getSupportedBackends
        ? this.displayableBackends.filter((backend) => engineType === backend.get('engineType')).sortBy('id')
        : [];
    } else {
      return getSupportedBackends
        ? []
        : this.displayableBackends
            // exact same as supportedBackends but negated
            .filter((backend) => engineType === backend.get('engineType'))
            .sortBy('id');
    }
  }

  @action
  filterEngineType([type]) {
    this.selectedEngineType = type;
  }

  @action
  filterEngineName([name]) {
    this.selectedEngineName = name;
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
}
