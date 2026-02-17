/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { dropTask } from 'ember-concurrency';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type VersionService from 'vault/services/version';
import type WizardService from 'vault/services/wizard';
import { WIZARD_ID } from '../wizard/secret-engines/secret-engines-wizard';

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
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;
  @service declare readonly wizard: WizardService;

  @tracked secretEngineOptions: Array<string> | [] = [];
  @tracked engineToDisable: SecretsEngineResource | undefined = undefined;
  @tracked enginesToDisable: Array<SecretsEngineResource> | null = null;

  @tracked engineTypeFilters: Array<string> = [];
  @tracked engineVersionFilters: Array<string> = [];
  @tracked searchText = '';

  // search text for dropdown filters
  @tracked typeSearchText = '';
  @tracked versionSearchText = '';

  @tracked selectedItems = Array<string>();

  @tracked shouldRenderIntroModal = false;
  wizardId = WIZARD_ID;

  tableColumns = [
    {
      key: 'path',
      label: 'Engine path',
      isSortable: true,
      width: '250px',
      customTableItem: true,
    },
    {
      key: 'accessor',
      label: 'Accessor',
      width: '175px',
    },
    {
      key: 'description',
      label: 'Description',
      width: '300px',
    },
    {
      key: 'running_plugin_version',
      label: 'Version',
      isSortable: true,
      width: '170px',
    },
    {
      key: 'popupMenu',
      label: 'Action',
      width: '75px',
    },
  ];

  get breadcrumbs() {
    return [
      {
        label: 'Vault',
        route: 'vault.cluster.dashboard',
        icon: 'vault',
      },
      {
        label: 'Secrets engines',
      },
    ];
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
      sortedBackends = sortedBackends.filter((backend) => {
        const effectiveType = getEffectiveEngineType(backend.engineType);
        return this.engineTypeFilters.includes(effectiveType);
      });
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
      return this.displayableBackends.filter((backend) => {
        const effectiveType = getEffectiveEngineType(backend.engineType);
        return effectiveType.toLowerCase().includes(this.typeSearchText.toLowerCase());
      });
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
    const arrayOfAllEffectiveTypes = this.typeFilterOptions.map((modelObject) =>
      getEffectiveEngineType(modelObject.engineType)
    );
    // filter out repeated effective types (e.g. [kv, kv] => [kv])
    const arrayOfUniqueEffectiveTypes = [...new Set(arrayOfAllEffectiveTypes)];

    return arrayOfUniqueEffectiveTypes.map((effectiveType) => ({
      name: effectiveType,
      id: effectiveType,
      icon: engineDisplayData(effectiveType)?.glyph ?? 'lock',
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

  // The backend does not directly indicate which engines were mounted by default and which have been mounted by the user
  // Currently the cubbyhole/, sys/, identity/ engines are mounted by default. (secret/ is mounted in dev mode as well)
  // The sys/ and identity/ engines are non-displayable engines.
  // While not ideal, we can check whether there are other engines than the default cubbyhole/ engine
  // to determine whether we should show the intro page
  get hasOnlyDefaultEngines() {
    const listedEngines = this.sortedDisplayableBackends;
    return !listedEngines.length || (listedEngines.length === 1 && listedEngines[0]?.path === 'cubbyhole/');
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && this.hasOnlyDefaultEngines;
  }

  @action
  showIntroPage() {
    // Reset the wizard dismissal state to allow re-entering the wizard
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshSecretEngineList() {
    this.router.refresh('vault.cluster.secrets.backends');
  }

  // Returns engine resource data for a given engine path, needed to get icon and other metadata from SecretEnginesResource
  getEngineResourceData = (enginePath: string) => {
    return this.displayableBackends.find((backend) => backend.path === enginePath);
  };

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
    } else if (!ALL_ENGINES.find((engine) => engine.type === backend.type)) {
      // If a mounted engine type doesn't match any known type in our static metadata, set this tooltip.
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

  @action
  updateSelectedItems(tableData: { selectedRowsKeys: string[] }) {
    this.selectedItems = tableData.selectedRowsKeys;
  }

  async disableSingleEngine(engine: SecretsEngineResource) {
    const { engineType, id, path } = engine;
    try {
      await this.api.sys.mountsDisableSecretsEngine(id);
      this.flashMessages.success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(
        `There was an error disabling the ${engineType} Secrets Engine at ${path}: ${message}.`
      );
    }
  }

  @dropTask
  *disableMultipleEngines(enginePathsToDisable: Array<string>) {
    const enginesToDisable = this.displayableBackends.filter((engine: SecretsEngineResource) =>
      enginePathsToDisable.includes(engine.path)
    );
    try {
      for (const engine of enginesToDisable) {
        yield this.disableSingleEngine(engine);
      }

      // Navigate once all operations are complete
      this.router.transitionTo('vault.cluster.secrets.backends');
    } finally {
      this.enginesToDisable = null;
    }
  }

  @dropTask
  *disableEngine(engine: SecretsEngineResource) {
    try {
      yield this.disableSingleEngine(engine);
      this.router.transitionTo('vault.cluster.secrets.backends');
    } finally {
      this.engineToDisable = undefined;
    }
  }
}
