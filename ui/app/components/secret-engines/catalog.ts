/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { filterEnginesByMountCategory } from 'core/utils/all-engines-metadata';
import type { PluginCatalogData } from 'vault/services/plugin-catalog';
import {
  categorizeEnginesByStatus,
  EnhancedEngineDisplayData,
  enhanceEnginesWithCatalogData,
  MOUNT_CATEGORIES,
  PLUGIN_CATEGORIES,
  PLUGIN_TYPES,
} from 'vault/utils/plugin-catalog-helpers';

import type VersionService from 'vault/services/version';
import type FlashMessageService from 'vault/services/flash-messages';

/**
 * @module SecretEnginesCatalog
 * SecretEnginesCatalog component displays available secret engines in a catalog view
 * for selection when mounting a new secret engine.
 *
 * @example
 * ```js
 * <SecretEngines::Catalog @setMountType={{this.setMountType}} @pluginCatalogData={{this.pluginCatalogData}} @pluginCatalogError={{this.pluginCatalogError}} />
 * ```
 */

interface Args {
  setMountType: (type: string) => void;
  pluginCatalogData?: PluginCatalogData;
  pluginCatalogError?: boolean;
}

export default class SecretEnginesCatalogComponent extends Component<Args> {
  @service declare version: VersionService;
  @service declare flashMessages: FlashMessageService;

  @tracked selectedEngineType = '';
  @tracked showFlyout = false;
  @tracked flyoutPlugin: unknown = null;
  @tracked flyoutPluginType: string | null = null;
  @tracked keywords = '';
  @tracked secretTypeFilter: string | null = null;
  @tracked rotationTypeFilter: string | null = null;
  @tracked platformFilter: string | null = null;

  secretTypes = [
    { label: 'Database credentials', value: 'databaseCredentials' },
    { label: 'API keys & tokens', value: 'apiKeysTokens' },
    { label: 'Encryption keys', value: 'encryptionKeys' },
    { label: 'Certificates & PKI', value: 'certificatesPki' },
    { label: 'Static storage', value: 'staticStorage' },
    { label: 'Cloud credentials', value: 'cloudCredentials' },
    { label: 'SSH keys', value: 'sshKeys' },
  ];

  rotationTypes = [
    { label: 'Static credentials', capability: 'static' },
    { label: 'Dynamic credentials', capability: 'dynamic' },
    { label: 'Rotation-enabled', capability: 'rotating' },
  ];

  platforms = [
    { label: 'Common engines', category: PLUGIN_CATEGORIES.COMMON },
    { label: 'Identity and access', category: PLUGIN_CATEGORIES.IDENTITY },
    { label: 'Cryptography and data protection', category: PLUGIN_CATEGORIES.CRYPTO },
    { label: 'Cloud and infrastructure', category: PLUGIN_CATEGORIES.CLOUD_PLUS },
  ];

  badgeLegendEntries = [
    { badge: 'Static', description: 'Can store static credentials' },
    { badge: 'Rotating', description: 'Can perform auto-rotation' },
    { badge: 'Dynamic', description: 'Can generate temporary credentials' },
    { badge: 'Encryption', description: 'Can perform encryption as a service' },
    { badge: 'Signing', description: 'Can sign keys' },
    { badge: 'Tokenization', description: 'Can tokenize sensitive data' },
    { badge: 'Certificate authority', description: 'Can issue and manage certificates' },
  ];

  get breadcrumbs() {
    return [
      {
        label: 'Vault',
        icon: 'vault',
        route: 'vault.cluster',
      },
      {
        label: 'Secrets engines',
        route: 'vault.cluster.secrets.backends',
      },
      {
        label: 'Create a new secrets engine',
      },
    ];
  }

  get secretEngines() {
    // If an enterprise license is present, return all secret engines;
    // otherwise, return only the secret engines supported in OSS.
    const staticEngines = filterEnginesByMountCategory({
      mountCategory: MOUNT_CATEGORIES.SECRET,
      isEnterprise: !!this.version?.isEnterprise,
    });

    // If we have plugin catalog data, merge it with static engines to add catalog info
    if (this.args.pluginCatalogData) {
      const secretEnginesDetailed =
        this.args.pluginCatalogData.detailed?.filter((plugin) => plugin?.type === PLUGIN_TYPES.SECRET) || [];
      const databasePluginsDetailed =
        this.args.pluginCatalogData.detailed?.filter((plugin) => plugin?.type === PLUGIN_TYPES.DATABASE) ||
        [];

      return enhanceEnginesWithCatalogData(staticEngines, secretEnginesDetailed, databasePluginsDetailed);
    }

    return staticEngines;
  }

  get filteredEngines() {
    const kw = this.keywords.replace(/\s+/g, ' ').trim().toLowerCase();

    return this.secretEngines.filter((engine) => {
      if (kw) {
        const matchesKeyword =
          engine.displayName?.toLowerCase().includes(kw) ||
          engine.type?.toLowerCase().includes(kw) ||
          engine.description?.toLowerCase().includes(kw) ||
          engine.capabilities?.find((capability) => capability.toLowerCase().includes(kw));
        if (!matchesKeyword) return false;
      }

      if (this.secretTypeFilter) {
        if (!(engine.secretTypes ?? []).includes(this.secretTypeFilter)) return false;
      }

      if (this.rotationTypeFilter) {
        if (!(engine.capabilities ?? []).includes(this.rotationTypeFilter)) return false;
      }

      if (this.platformFilter && engine.category !== this.platformFilter) {
        return false;
      }

      return true;
    });
  }

  get pluginCategoriesList() {
    return [
      PLUGIN_CATEGORIES.COMMON,
      PLUGIN_CATEGORIES.IDENTITY,
      PLUGIN_CATEGORIES.CRYPTO,
      PLUGIN_CATEGORIES.CLOUD_PLUS,

      // TODO: enable external plugins once version selection is available (VAULT-39241)
      // PLUGIN_CATEGORIES.EXTERNAL,
    ];
  }

  get secretMountCategory() {
    return MOUNT_CATEGORIES.SECRET;
  }

  isDisabled = (type: EnhancedEngineDisplayData) => {
    return (
      (type.requiresEnterprise && !this.version.isEnterprise) ||
      (type.requiredFeature && !this.hasFeature(type.requiredFeature))
    );
  };

  showDeprecationBadge = (type: EnhancedEngineDisplayData) => {
    return type.deprecationStatus && type.deprecationStatus !== 'supported';
  };

  hasFeature(featureName: string) {
    return this.version.features?.includes(featureName) || false;
  }

  clearSelectedEngine() {
    this.selectedEngineType = '';
  }

  @action
  selectEngineType(type: string) {
    this.selectedEngineType = type;
  }

  @action
  handleSelection() {
    this.args.setMountType(this.selectedEngineType);
  }

  @action
  getMountTypesByCategory(category: string) {
    const allTypes = this.filteredEngines.filter((engine) => engine.category === category);
    return categorizeEnginesByStatus(allTypes);
  }

  @action
  handleDisabledPluginClick(plugin: unknown) {
    this.showFlyout = true;
    this.flyoutPlugin = plugin;
    this.flyoutPluginType = 'secret';
  }

  @action
  handleDisabledPluginKeyDown(plugin: unknown, event: KeyboardEvent) {
    // Only handle Enter and Space keys for accessibility
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      this.handleDisabledPluginClick(plugin);
    }
  }

  @action
  openExternalPluginsHelp() {
    this.showFlyout = true;
    this.flyoutPlugin = null;
    this.flyoutPluginType = 'secret';
  }

  @action
  closeFlyout() {
    this.showFlyout = false;
    this.flyoutPlugin = null;
    this.flyoutPluginType = null;
  }

  @action
  filterBySecretType(value: string) {
    this.secretTypeFilter = this.secretTypeFilter === value ? null : value;
    this.clearSelectedEngine();
  }

  @action
  filterByPlatform(category: string) {
    this.platformFilter = this.platformFilter === category ? null : category;
    this.clearSelectedEngine();
  }

  @action
  filterByRotationType(capability: string) {
    this.rotationTypeFilter = this.rotationTypeFilter === capability ? null : capability;
    this.clearSelectedEngine();
  }

  @action
  clearKeyword() {
    this.keywords = '';
    this.clearSelectedEngine();
  }

  @action
  clearAllFilters() {
    this.keywords = '';
    this.secretTypeFilter = null;
    this.rotationTypeFilter = null;
    this.platformFilter = null;
    this.clearSelectedEngine();
  }

  @action
  setSearchText(type: string, event: Event) {
    const target = event.target as HTMLInputElement;
    if (type === 'keywords') {
      if (target.value.trim()) {
        this.keywords = target.value;
        this.clearSelectedEngine();
      }
    }
  }
}
