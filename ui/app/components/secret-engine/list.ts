/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { dropTask } from 'ember-concurrency';

import type FlashMessageService from 'vault/services/flash-messages';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';
import engineDisplayData from 'vault/helpers/engines-display-data';

/**
 * @module SecretEngineList handles the display of the list of secret engines, including the filtering.
 * 
 * @example
 * <SecretEngine::List
    @secretEngines={{this.model}}
    />
 *
 * @param {array} secretEngines - An array of Secret Engine models returned from query on the parent route.
 */

interface Args {
  secretEngines: Array<SecretsEngineResource>;
  onHover: (tooltip: string) => void;
}

export default class SecretEngineList extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;

  @tracked secretEngineOptions: Array<string> | [] = [];
  @tracked selectedEngineType = '';
  @tracked selectedEngineName = '';
  @tracked engineToDisable: SecretsEngineResource | undefined = undefined;
  @tracked tooltip = '';

  get displayableBackends() {
    return this.args.secretEngines.filter((backend) => backend.shouldIncludeInList);
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

  @action
  tooltipText(backend: SecretsEngineResource) {
    if (backend.isSupportedBackend) {
      if (backend.type === 'kv') {
        // If the backend is a KV engine, include the version in the tooltip.
        this.tooltip = `${engineDisplayData(backend.type)?.displayName} version ${backend.version}`;
      } else {
        this.tooltip = `${engineDisplayData(backend.type)?.displayName}`;
      }
    } else if (engineDisplayData(backend.type)?.type === 'generic') {
      // If a mounted engine type doesn't match any known type, the type is returned as 'generic' and set this tooltip.
      // Handles issue when a user externally mounts an engine that doesn't follow the expected naming conventions for what's in the binary, despite being a valid engine.
      this.tooltip =
        'This is an externally mounted plugin that is not recognized by the UI. Please use the CLI to manage this engine.';
    } else {
      // If the engine type is recognized but not supported, we only show configuration view and set this tooltip.
      this.tooltip =
        'The UI only supports configuration views for these secret engines. The CLI must be used to manage other engine resources.';
    }

    this.args.onHover?.(this.tooltip); // expose the tracked value
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
  *disableEngine(engine: SecretsEngineResource) {
    const { engineType, id, path } = engine;
    try {
      yield this.api.sys.mountsDisableSecretsEngine(id);
      this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
      this.router.transitionTo('vault.cluster.secrets.backends');
    } catch (err) {
      const { message } = yield this.api.parseError(err);
      this.flashMessages.danger(
        `There was an error disabling the ${engineType} Secrets Engines at ${path}: ${message}.`
      );
    } finally {
      this.engineToDisable = undefined;
    }
  }
}
