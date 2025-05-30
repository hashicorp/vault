/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { dropTask } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

import type FlashMessageService from 'vault/services/flash-messages';
import SecretEngineModel from 'vault/models/secret-engine';

/**
 * @module SecretEngineList handles the display of the list of secret engines, including the filtering.
 * 
 * @example
 * <SecretEngine::List
    @secretEngineModels={{this.model}}
    />
 *
 * @param {array} secretEngineModels - An array of Secret Engine models returned from query on the parent route.
 */

interface Args {
  secretEngineModels: Array<SecretEngineModel>;
}

export default class SecretEngineList extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked secretEngineOptions: Array<string> | [] = [];
  @tracked selectedEngineType = '';
  @tracked selectedEngineName = '';
  @tracked engineToDisable: SecretEngineModel | undefined = undefined;

  get displayableBackends() {
    return this.args.secretEngineModels.filter((backend) => backend.shouldIncludeInList);
  }

  get sortedDisplayableBackends() {
    // show supported secret engines first and then organize those by id.
    const sortedBackends = this.displayableBackends.sort(
      (a, b) => Number(b.isSupportedBackend) - Number(a.isSupportedBackend) || a.id.localeCompare(b.id)
    );

    // return an options list to filter by engine type, ex: 'kv'
    if (this.selectedEngineType) {
      // check first if the user has also filtered by name.
      if (this.selectedEngineName) {
        return sortedBackends.filter((backend) => this.selectedEngineName === backend.id);
      }
      // otherwise filter by engine type
      return sortedBackends.filter((backend) => this.selectedEngineType === backend.engineType);
    }

    // return an options list to filter by engine name, ex: 'secret'
    if (this.selectedEngineName) {
      return sortedBackends.filter((backend) => this.selectedEngineName === backend.id);
    }
    // no filters, return full sorted list.
    return sortedBackends;
  }

  // Filtering & searching
  get secretEngineArrayByType() {
    const arrayOfAllEngineTypes = this.sortedDisplayableBackends.map((modelObject) => modelObject.engineType);
    // filter out repeated engineTypes (e.g. [kv, kv] => [kv])
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
  filterEngineType(type: string[]) {
    const [selectedType] = type;
    this.selectedEngineType = selectedType || '';
  }

  @action
  filterEngineName(name: string[]) {
    const [selectedName] = name;
    this.selectedEngineName = selectedName || '';
  }

  @dropTask
  *disableEngine(engine: SecretEngineModel) {
    const { engineType, path } = engine;
    try {
      yield engine.destroyRecord();
      this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling the ${engineType} Secrets Engines at ${path}: ${errorMessage(err)}.`
      );
    } finally {
      this.engineToDisable = undefined;
    }
  }
}
