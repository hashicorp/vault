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

export default class VaultClusterSecretsBackendController extends Controller {
  @service flashMessages;
  @filterBy('model', 'shouldIncludeInList') displayableBackends;

  @tracked secretEngineOptions = [];
  @tracked selectedEngineType = null;
  @tracked selectedEngineName = null;

  // ARG TODO: check multiple types of same secret Engine
  // check pagination.
  // ARG TMRW: solve for when both filters are set
  // test with hundreds of mounts

  get sortedDisplayableBackends() {
    // show supported secret engines first and then organize those by id.
    const sortedDisplayableBackends = this.displayableBackends.sort(
      (a, b) => b.isSupportedBackend - a.isSupportedBackend || a.id - b.id
    );
    // filter list by engine type, ex: 'kv'
    if (this.selectedEngineType) {
      return sortedDisplayableBackends.filter(
        (backend) => this.selectedEngineType === backend.get('engineType')
      );
    }
    // filter list by engine name, ex: 'secret'
    if (this.selectedEngineName) {
      // const engineObject = sortedDisplayableBackends.find(
      //   (modelObject) => modelObject.id === this.selectedEngineName
      // );
      return sortedDisplayableBackends.filter((backend) => this.selectedEngineName === backend.get('id'));
    }
    // if no filter, return full sorted list
    return sortedDisplayableBackends;
  }

  get secretEngineArrayByType() {
    const arrayOfAllEngineTypes = this.sortedDisplayableBackends.map((modelObject) => modelObject.engineType);
    // filter out repeated engineTypes (e.g. [secret, secret] => [secret])
    const arrayOfUniqueEngineTypes = [...new Set(arrayOfAllEngineTypes)];

    return arrayOfUniqueEngineTypes.map((engineType) => ({
      name: engineType,
      id: engineType,
    }));
  }

  get secretEngineArrayByName() {
    return this.sortedDisplayableBackends.map((modelObject) => ({
      name: modelObject.id,
      id: modelObject.id,
    }));
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
