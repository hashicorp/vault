/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import { dropTask } from 'ember-concurrency';
import sortObjects from 'vault/utils/sort-objects';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_REOPEN_CLICKED } from 'vault/utils/analytic-events';

import type AnalyticsService from 'vault/services/analytics';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';
import type { Breadcrumb } from 'vault/vault/app-types';
import type AuthMethodResource from 'vault/resources/auth/method';
import engineDisplayData from 'core/helpers/engines-display-data';

/**
 * @module PageAuthMethods
 * PageAuthMethods component handles the display and management of authentication methods.
 *
 * @param {object} model - contains methods array and capabilities
 * @param {array} breadcrumbs - breadcrumb navigation items
 */

interface Args {
  model: {
    methods: AuthMethodResource[];
    capabilities: unknown;
  };
  breadcrumbs: Breadcrumb[];
}

export default class PageAuthMethodsComponent extends Component<Args> {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly wizard: WizardService;

  @tracked authMethodOptions = [];
  @tracked methodToDisable: AuthMethodResource | null = null;
  @tracked shouldRenderIntroModal = false;
  @tracked searchText = '';
  @tracked authTypeFilters: Array<string> = [];

  // search text for dropdown filters
  @tracked typeSearchText = '';
  wizardId = WIZARD_ID_MAP.authMethods;

  tableColumns = [
    {
      key: 'path',
      label: 'Auth name',
      isSortable: true,
      customTableItem: true,
    },
    {
      key: 'accessor',
      label: 'Accessor',
    },
    {
      key: 'description',
      label: 'Description',
    },
    {
      key: 'popupMenu',
      label: 'Action',
      width: '8%',
    },
  ];

  // list returned by getter is sorted in template
  get authMethodList() {
    const { methods } = this.args.model;

    let filteredMethodList = methods
      .slice()
      .sort((a, b) => Number(b) - Number(a) || a.id.localeCompare(b.id));

    // check for any auth type filters
    if (this.authTypeFilters.length > 0) {
      filteredMethodList = filteredMethodList.filter((method) => {
        return this.authTypeFilters.includes(method.type);
      });
    }

    // if there is search text, filter path name by that
    if (this.searchText.trim() !== '') {
      filteredMethodList = filteredMethodList.filter((backend) =>
        backend.path.toLowerCase().includes(this.searchText.toLowerCase())
      );
    }
    // no filters, return full sorted list.
    return filteredMethodList;
  }

  // Returns filter options for engine type dropdown
  get typeFilterOptions() {
    const { methods } = this.args.model;
    // if there is search text, filter types by that
    if (this.typeSearchText.trim() !== '') {
      return methods.filter((method) => {
        return method.id.toLowerCase().includes(this.typeSearchText.toLowerCase());
      });
    }
    return methods;
  }

  get authMethodArrayByType() {
    const arrayOfAllAuthTypes = this.typeFilterOptions.map((modelObject) => modelObject.type);
    // filter out repeated auth types (e.g. [userpass, userpass] => [userpass])
    const arrayOfUniqueAuthTypes = [...new Set(arrayOfAllAuthTypes)];

    return arrayOfUniqueAuthTypes.map((authType) => ({
      name: authType,
      id: authType,
      icon: engineDisplayData(authType).glyph ?? 'lock',
    }));
  }

  get authMethodArrayByName() {
    return this.authMethodList.map((modelObject) => ({
      name: modelObject.id,
      id: modelObject.id,
    }));
  }

  get hasOnlyDefaultMethods() {
    return this.args.model.methods.length === 1;
  }

  get showContent() {
    // Show when the 1) wizard is not shown OR 2) wizard intro modal is shown
    // This ensures the wizard intro modal is shown on top of the list view and the background content is not blank behind the modal
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  get showIntroButton() {
    return this.showContent && this.hasOnlyDefaultMethods;
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && this.hasOnlyDefaultMethods;
  }

  get showPageHeader() {
    return !this.showWizard || this.wizard.isIntroVisible(this.wizardId);
  }

  getAuthMethodData = (path: string) => {
    return this.authMethodList.find((method) => method.path === path);
  };

  @action
  setSearchText(type: string, event: Event) {
    const target = event.target as HTMLInputElement;
    if (type === 'type') {
      this.typeSearchText = target.value;
    } else {
      this.searchText = target.value;
    }
  }

  @action
  filterByAuthType(type: string) {
    if (this.authTypeFilters.includes(type)) {
      this.authTypeFilters = this.authTypeFilters.filter((t) => t !== type);
    } else {
      this.authTypeFilters = [...this.authTypeFilters, type];
    }
  }

  @action
  clearAllFilters() {
    this.authTypeFilters = [];
  }

  @action
  showIntroPage() {
    this.analytics.trackEvent(INTRO_REOPEN_CLICKED, {
      namespace: 'intro-page',
      action: 'clicked',
      elementId: 'intro-reopen-button',
      channel: 'webpage',
      objectType: 'auth-method',
    });
    // Reset the wizard dismissal state to allow re-entering the wizard
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  async refreshMethodsList() {
    this.router.refresh('vault.cluster.access.methods');
  }

  @dropTask
  *disableMethod(method: AuthMethodResource) {
    const { type, path } = method;
    try {
      yield this.api.sys.authDisableMethod(path);
      this.flashMessages.success(`The ${type} Auth Method at ${path} has been disabled.`);
      this.refreshMethodsList();
    } catch (err) {
      const { message } = yield this.api.parseError(err);
      this.flashMessages.danger(`There was an error disabling Auth Method at ${path}: ${message}.`);
    } finally {
      this.methodToDisable = null;
    }
  }

  // template helper
  sortMethods = (methods: AuthMethodResource[]) => {
    // make sure there are methods to sort otherwise slice with throw an error
    if (!Array.isArray(methods) || methods.length === 0) return [];
    return sortObjects(methods.slice(), 'path');
  };
}
