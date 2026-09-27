/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_REOPEN_CLICKED } from 'vault/utils/analytic-events';

import type RouterService from '@ember/routing/router-service';
import type AnalyticsService from 'vault/services/analytics';
import type WizardService from 'vault/services/wizard';

interface Engine {
  id: string;
  path: string;
  engineType: string;
  icon: string;
  displayName: string;
  version: string;
  [key: string]: unknown;
}

interface RouteModel {
  engines: Engine[];
  listViewConfig: object;
  page: number;
  pageSize: number;
}

export default class VaultClusterSecretsBackendController extends Controller {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly router: RouterService;
  @service declare readonly wizard: WizardService;

  declare model: RouteModel;

  // page refreshes the model (re-fetches from API) on change.
  // pageSize is client-side only — no model reload needed, but must survive route transitions.
  queryParams = ['page', 'pageSize'];
  @tracked page = 1;
  @tracked pageSize = 10;

  @tracked shouldRenderIntroModal = false;
  @tracked engineTypeFilters: string[] = [];
  @tracked engineVersionFilters: string[] = [];

  // search text for all three filter inputs
  @tracked pathSearchText = '';
  @tracked typeSearchText = '';
  @tracked versionSearchText = '';

  // Returns unique engine types matching the current type search text
  get secretEngineArrayByType(): { name: string; icon: string }[] {
    const engines = this.model?.engines ?? [];
    const filtered =
      this.typeSearchText.trim() !== ''
        ? engines.filter((e) =>
            getEffectiveEngineType(e.engineType).toLowerCase().includes(this.typeSearchText.toLowerCase())
          )
        : engines;

    const uniqueTypes = [...new Set(filtered.map((e) => getEffectiveEngineType(e.engineType)))];
    return uniqueTypes.map((type) => ({
      name: type,
      icon: engineDisplayData(type)?.glyph ?? 'lock',
    }));
  }

  // Returns unique engine versions from the type-filtered set, matching the version search text.
  // Version dropdown is only shown after a type filter is selected so this narrows from the
  // already-type-filtered list.
  get secretEngineArrayByVersions(): { version: string }[] {
    const engines = this.model?.engines ?? [];

    // Narrow to type-filtered engines first (mirrors the component behaviour)
    const typeFiltered =
      this.engineTypeFilters.length > 0
        ? engines.filter((e) => this.engineTypeFilters.includes(getEffectiveEngineType(e.engineType)))
        : engines;

    const searchFiltered =
      this.versionSearchText.trim() !== ''
        ? typeFiltered.filter((e) => e.version.toLowerCase().includes(this.versionSearchText.toLowerCase()))
        : typeFiltered;

    const uniqueVersions = [...new Set(searchFiltered.map((e) => e.version))];
    return uniqueVersions.map((version) => ({ version }));
  }

  // Engines passed to @model after all active filters are applied.
  get filteredEngines(): Engine[] {
    const engines = this.model?.engines ?? [];
    let result = engines;

    if (this.pathSearchText.trim() !== '') {
      const q = this.pathSearchText.toLowerCase();
      result = result.filter((e) => e.path.toLowerCase().includes(q));
    }

    if (this.engineTypeFilters.length > 0) {
      result = result.filter((e) => this.engineTypeFilters.includes(getEffectiveEngineType(e.engineType)));
    }

    if (this.engineVersionFilters.length > 0) {
      result = result.filter((e) => this.engineVersionFilters.includes(e.version));
    }

    return result;
  }

  @action
  setSearchText(type: string, event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    if (type === 'path') {
      this.pathSearchText = value;
    } else if (type === 'type') {
      this.typeSearchText = value;
    } else if (type === 'version') {
      this.versionSearchText = value;
    }
  }

  @action
  filterByEngineType(type: string): void {
    if (this.engineTypeFilters.includes(type)) {
      this.engineTypeFilters = this.engineTypeFilters.filter((t) => t !== type);
    } else {
      this.engineTypeFilters = [...this.engineTypeFilters, type];
    }
    // Clear version filters when type selection changes to avoid stale combinations
    this.engineVersionFilters = [];
  }

  @action
  filterByEngineVersion(version: string): void {
    if (this.engineVersionFilters.includes(version)) {
      this.engineVersionFilters = this.engineVersionFilters.filter((v) => v !== version);
    } else {
      this.engineVersionFilters = [...this.engineVersionFilters, version];
    }
  }

  @action
  clearAllFilters(): void {
    this.engineTypeFilters = [];
    this.engineVersionFilters = [];
    this.pathSearchText = '';
    this.typeSearchText = '';
    this.versionSearchText = '';
  }

  // True when only the default cubbyhole/ engine is present (no user-mounted engines yet).
  // Used to decide whether to show the intro button and wizard.
  get hasOnlyDefaultEngines(): boolean {
    const engines = this.model?.engines ?? [];
    return !engines.length || (engines.length === 1 && engines[0]?.path === 'cubbyhole/');
  }

  get showWizard(): boolean {
    return !this.wizard.isDismissed(WIZARD_ID_MAP.secretEngines) && this.hasOnlyDefaultEngines;
  }

  get showIntroButton(): boolean {
    return this.hasOnlyDefaultEngines;
  }

  @action
  showIntroPage(): void {
    this.analytics.trackEvent(INTRO_REOPEN_CLICKED, {
      namespace: 'intro-page',
      action: 'clicked',
      elementId: 'intro-reopen-button',
      channel: 'webpage',
      objectType: 'secrets-engine',
    });
    this.wizard.reset(WIZARD_ID_MAP.secretEngines);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshEngineList(): void {
    this.router.refresh('vault.cluster.secrets.backends');
  }
}
