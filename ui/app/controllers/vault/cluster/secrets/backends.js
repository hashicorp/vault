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
import { dropTask } from 'ember-concurrency';

export default class VaultClusterSecretsBackendController extends Controller {
  @service flashMessages;
  @filterBy('model', 'shouldIncludeInList') displayableBackends;

  @tracked secretEngineOptions = [];
  @tracked selectedEngineType = null;
  @tracked selectedEngineName = null;

  // check pagination.
  // test with hundreds of mounts

  get sortedDisplayableBackends() {
    // show supported secret engines first and then organize those by id.
    const sortedDisplayableBackends = this.displayableBackends.sort(
      (a, b) => b.isSupportedBackend - a.isSupportedBackend || a.id - b.id
    );

    // filter list by engine type, ex: 'kv'
    if (this.selectedEngineType) {
      // check first if the user has also filtered by name.
      if (this.selectedEngineName) {
        return this.filterByName(sortedDisplayableBackends);
      }
      // otherwise filter by engine type
      return this.filterByEngineType(sortedDisplayableBackends);
    }

    // filter list by engine name, ex: 'secret'
    if (this.selectedEngineName) {
      return this.filterByName(sortedDisplayableBackends);
    }
    // no filters, return full sorted list.
    return sortedDisplayableBackends;
  }

  filterByName(backendList) {
    return backendList.filter((backend) => this.selectedEngineName === backend.get('id'));
  }

  filterByEngineType(backendList) {
    return backendList.filter((backend) => this.selectedEngineType === backend.get('engineType'));
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
