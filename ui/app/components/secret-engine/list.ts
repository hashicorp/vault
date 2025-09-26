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
import type VersionService from 'vault/services/version';
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
}

export default class SecretEngineList extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;

  @tracked secretEngineOptions: Array<string> | [] = [];
  @tracked engineToDisable: SecretsEngineResource | undefined = undefined;

  @tracked engineTypeFilters: Array<string> = [];
  @tracked engineVersionFilters: Array<string> = [];
  @tracked searchText = '';

  // search text for dropdown filters
  @tracked typeSearchText = '';
  @tracked versionSearchText = '';

  get clusterName() {
    return this.version.clusterName;
  }

  get displayableBackends() {
    return this.args.secretEngines.filter((backend) => backend.shouldIncludeInList);
  }

  get sortedDisplayableBackends() {
    // show supported secret engines first and then organize those by id.
    let sortedBackends = this.displayableBackends
      .slice()
      .sort(
        (a, b) => Number(b.isSupportedBackend) - Number(a.isSupportedBackend) || a.id.localeCompare(b.id)
      );

    // filters by engine type, ex: 'kv'
    if (this.engineTypeFilters.length > 0) {
      sortedBackends = sortedBackends.filter((backend) =>
        this.engineTypeFilters.includes(backend.engineType)
      );
    }

    // filters by engine version, ex: 'v1.21.0...'
    if (this.engineVersionFilters.length > 0) {
      sortedBackends = sortedBackends.filter((backend) =>
        this.engineVersionFilters.includes(backend.running_plugin_version)
      );
    }

    // if there is search text, filter path name by that
    if (this.searchText.trim() !== '') {
      sortedBackends = sortedBackends.filter((backend) =>
        backend.path.toLowerCase().includes(this.searchText.toLowerCase())
      );
    }
    // no filters, return full sorted list.
    return sortedBackends;
  }

  // Returns filter options for engine type dropdown
  get typeFilterOptions() {
    // if there is search text, filter types by that
    if (this.typeSearchText.trim() !== '') {
      return this.displayableBackends.filter((backend) =>
        backend.engineType.toLowerCase().includes(this.typeSearchText.toLowerCase())
      );
    }

    return this.displayableBackends;
  }

  // Returns filter options for version dropdown
  get versionFilterOptions() {
    // if there is search text, filter versions by that
    if (this.versionSearchText.trim() !== '') {
      // filtered by sorted backends array since an engine type filter has to be selected first
      return this.sortedDisplayableBackends.filter((backend) =>
        backend.running_plugin_version.toLowerCase().includes(this.versionSearchText.toLowerCase())
      );
    }
    return this.sortedDisplayableBackends;
  }

  // Returns filtered engines list by type
  get secretEngineArrayByType() {
    const arrayOfAllEngineTypes = this.typeFilterOptions.map((modelObject) => modelObject.engineType);
    // filter out repeated engineTypes (e.g. [kv, kv] => [kv])
    const arrayOfUniqueEngineTypes = [...new Set(arrayOfAllEngineTypes)];

    return arrayOfUniqueEngineTypes.map((engineType) => ({
      name: engineType,
      id: engineType,
      icon: engineDisplayData(engineType)?.glyph ?? 'lock',
    }));
  }

  // Returns filtered engines list by version
  get secretEngineArrayByVersions() {
    const arrayOfAllEngineVersions = this.versionFilterOptions.map(
      (modelObject) => modelObject.running_plugin_version
    );
    // filter out repeated engineVersions (e.g. [1.0, 1.0] => [1.0])
    const arrayOfUniqueEngineVersions = [...new Set(arrayOfAllEngineVersions)];
    return arrayOfUniqueEngineVersions.map((version) => ({
      version,
      id: version,
    }));
  }

  generateToolTipText = (backend: SecretsEngineResource) => {
    const displayData = engineDisplayData(backend.type);

    if (!displayData) {
      return;
    } else if (backend.isSupportedBackend) {
      if (backend.type === 'kv') {
        // If the backend is a KV engine, include the version in the tooltip.
        return `${displayData.displayName} version ${backend.version}`;
      } else {
        return `${displayData.displayName}`;
      }
    } else if (displayData.type === 'unknown') {
      // If a mounted engine type doesn't match any known type, the type is returned as 'unknown' and set this tooltip.
      // Handles issue when a user externally mounts an engine that doesn't follow the expected naming conventions for what's in the binary, despite being a valid engine.
      return `This engine's type is not recognized by the UI. Please use the CLI to manage this engine.`;
    } else {
      // If the engine type is recognized but not supported, we only show configuration view and set this tooltip.
      return 'The UI only supports configuration views for these secret engines. The CLI must be used to manage other engine resources.';
    }
  };

  @action
  setSearchText(type: string, event: Event) {
    const target = event.target as HTMLInputElement;
    if (type === 'type') {
      this.typeSearchText = target.value;
    } else if (type === 'version') {
      this.versionSearchText = target.value;
    } else {
      this.searchText = target.value;
    }
  }

  @action
  filterByEngineType(type: string) {
    if (this.engineTypeFilters.includes(type)) {
      this.engineTypeFilters = this.engineTypeFilters.filter((t) => t !== type);
    } else {
      this.engineTypeFilters = [...this.engineTypeFilters, type];
    }
  }

  @action
  filterByEngineVersion(version: string) {
    if (this.engineVersionFilters.includes(version)) {
      this.engineVersionFilters = this.engineVersionFilters.filter((v) => v !== version);
    } else {
      this.engineVersionFilters = [...this.engineVersionFilters, version];
    }
  }

  @action
  clearAllFilters() {
    this.engineTypeFilters = [];
    this.engineVersionFilters = [];
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
