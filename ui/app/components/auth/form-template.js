/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { allSupportedAuthBackends, supportedTypes } from 'vault/helpers/supported-auth-backends';

/**
 * @module Auth::FormTemplate
 * */

export default class AuthFormTemplate extends Component {
  @service flags;
  @service store;
  @service version;

  // form display logic
  @tracked authTabs = null; // listing visibility or backup types
  @tracked selectedTabIndex = 0;
  @tracked showOtherMethods = true;

  // auth login variables
  @tracked selectedAuthMethod = 'token';
  @tracked errorMessage = null;

  displayName = (type) => allSupportedAuthBackends().find((t) => t.type === type).typeDisplay;

  constructor() {
    super(...arguments);
    // TODO fetch login customization config here
    this.fetchMounts.perform();
  }

  get availableMethodTypes() {
    return supportedTypes(this.version.isEnterprise);
  }

  get formComponent() {
    const { selectedAuthMethod } = this;
    const isSupported = this.availableMethodTypes.includes(selectedAuthMethod);
    const formFile = () => (['oidc', 'jwt'].includes(selectedAuthMethod) ? 'oidc-jwt' : selectedAuthMethod);
    const component = isSupported ? formFile() : 'base';

    // an Auth::Form::<Type> component exists for each type in supported-auth-backends
    // eventually "base" component could be leveraged for rendering custom auth plugins
    return `auth/form/${component}`;
  }

  get namespaceInput() {
    const namespaceQP = this.args.namespace;
    if (this.flags.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/` for their nested namespace
      const split = namespaceQP.split('/');
      if (split.length > 1) {
        split.shift();
        return `/${split.join('/')}`;
      }
      return '';
    }
    return namespaceQP;
  }

  get showTabs() {
    switch (true) {
      // listing visibility is set and user has NOT clicked "Sign in with other"
      case this.authTabs && !this.showOtherMethods:
        return true;
      // TODO add case(s) for login customization (ent only)
      // case this.defaultConfigured && this.showOtherMethods:
      //   return true;
      // if only backups set (maybe not allowed)
      // case this.backupsConfigured && !this.showOtherMethods:
      //   return true;
      default:
        return false;
    }
  }

  @action
  handleAuthSelect(element, event, idx) {
    if (element === 'tab') {
      this.setAuthTypeFromTab(idx);
    } else {
      this.selectedAuthMethod = event.target.value;
    }
  }

  setAuthTypeFromTab(idx) {
    this.selectedAuthMethod = Object.keys(this.authTabs)[idx];
    this.selectedTabIndex = idx;
  }

  @action
  handleError(message) {
    this.errorMessage = message;
  }

  @action
  toggleView() {
    this.showOtherMethods = !this.showOtherMethods;

    if (this.showTabs) {
      // reset selected auth method to first tab
      this.handleAuthSelect('tab', null, 0);
    } else {
      // all methods render, reset dropdown
      this.selectedAuthMethod = 'token';
    }
  }

  // this will be SO MUCH NICER with the auth updates that remove ember data
  fetchMounts = task(
    waitFor(async () => {
      try {
        const unauthMounts = await this.store.findAll('auth-method', {
          adapterOptions: {
            unauthenticated: true,
          },
        });

        if (unauthMounts.length !== 0) {
          this.authTabs = unauthMounts.reduce((obj, m) => {
            // serialize the ember data model into a regular ol' object
            const mountData = m.serialize();
            const methodType = mountData.type;
            if (!Object.keys(obj).includes(methodType)) {
              // create a new empty array for that type
              obj[methodType] = [];
            }
            // push mount data into corresponding type's array
            obj[methodType].push(mountData);
            return obj;
          }, {});

          // set tracked selected auth type to first tab
          this.setAuthTypeFromTab(0);
          // hide other methods to prioritize tabs (visible mounts)
          this.showOtherMethods = false;
        }
      } catch (e) {
        // if for some reason there's an error fetching mounts, swallow and just show standard form
      }
    })
  );
}
