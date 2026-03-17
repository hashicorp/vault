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
import { WIZARD_ID } from '../wizard/methods/methods-wizard';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';
import type { Breadcrumb } from 'vault/vault/app-types';
import type AuthMethodResource from 'vault/resources/auth/method';

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
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly wizard: WizardService;

  @tracked authMethodOptions = [];
  @tracked selectedAuthType: string | null = null;
  @tracked selectedAuthName: string | null = null;
  @tracked methodToDisable: AuthMethodResource | null = null;
  @tracked shouldRenderIntroModal = false;

  wizardId = WIZARD_ID;

  // list returned by getter is sorted in template
  get authMethodList() {
    const { methods } = this.args.model;
    // return an options list to filter by engine type, ex: 'kv'
    if (this.selectedAuthType) {
      // check first if the user has also filtered by name.
      // names are individualized across type so you can't have the same name for an aws auth method and userpass.
      // this means it's fine to filter by first type and then name or just name.
      if (this.selectedAuthName) {
        return methods.filter((method) => this.selectedAuthName === method.id);
      }
      // otherwise filter by auth type
      return methods.filter((method) => this.selectedAuthType === method.type);
    }
    // return an options list to filter by auth name, ex: 'my-userpass'
    if (this.selectedAuthName) {
      return methods.filter((method) => this.selectedAuthName === method.id);
    }
    // no filters, return full list
    return methods;
  }

  get authMethodArrayByType() {
    const arrayOfAllAuthTypes = this.authMethodList.map((modelObject) => modelObject.type);
    // filter out repeated auth types (e.g. [userpass, userpass] => [userpass])
    const arrayOfUniqueAuthTypes = [...new Set(arrayOfAllAuthTypes)];

    return arrayOfUniqueAuthTypes.map((authType) => ({
      name: authType,
      id: authType,
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
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(WIZARD_ID));
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

  @action
  filterAuthType([type]: [string]) {
    this.selectedAuthType = type;
  }

  @action
  filterAuthName([name]: [string]) {
    this.selectedAuthName = name;
  }

  @action
  showIntroPage() {
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
