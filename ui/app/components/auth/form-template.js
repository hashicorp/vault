/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { restartableTask, timeout } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { ALL_LOGIN_METHODS, supportedTypes } from 'vault/utils/supported-login-methods';

/**
 * @module Auth::FormTemplate
 * This component is responsible for managing the layout and display logic for the auth form. When initialized it fetches
 * the unauthenticated sys/internal/ui/mounts endpoint to check the listing_visibility configuration of available mounts.
 * If mounts have been configured as listing_visibility="unauth" then tabs render for the corresponding method types,
 * otherwise all auth methods display in a dropdown list. The endpoint is re-requested anytime the namespace input is updated.
 *
 * When auth type changes (by selecting a new one from the dropdown or select a tab), the form component updates and
 * dynamically renders the corresponding form.
 *
 * @example
 * <Auth::FormTemplate
 *  @wrappedToken={{@wrappedToken}}
 *  @cluster={{@cluster}}
 *  @handleNamespaceUpdate={{this.handleNamespaceUpdate}}
 *  @namespace={{@namespaceQueryParam}}
 *  @onSuccess={{this.onAuthResponse}}
 *  />
 *
 * @param {string} wrappedToken - Query param value of a wrapped token that can be used to login when added directly to the URL via the "wrapped_token" query param
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} handleNamespaceUpdate - callback task that passes user input to the controller and updates the namespace query param in the url
 * @param {string} namespace - namespace query param from the url
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 *
 * */

export default class AuthFormTemplate extends Component {
  @service flags;
  @service store;
  @service version;

  // form display logic
  @tracked authTabs = null; // renders method types as tabs (default is to list in a dropdown)
  @tracked selectedTabIndex = 0;
  @tracked showOtherMethods = false;

  // auth login variables
  @tracked selectedAuthMethod = 'token';
  @tracked errorMessage = null;

  displayName = (type) => {
    const displayName = ALL_LOGIN_METHODS?.find((t) => t.type === type)?.displayName;
    return displayName || type;
  };

  constructor() {
    super(...arguments);
    this.fetchMounts.perform();
  }

  get availableMethodTypes() {
    return supportedTypes(this.version.isEnterprise);
  }

  get formComponent() {
    const { selectedAuthMethod } = this;
    // isSupported means there is a component file defined for that auth type
    const isSupported = this.availableMethodTypes.includes(selectedAuthMethod);
    const formFile = () => (['oidc', 'jwt'].includes(selectedAuthMethod) ? 'oidc-jwt' : selectedAuthMethod);
    const component = isSupported ? formFile() : 'base';

    // an Auth::Form::<Type> component exists for each method in supported-login-methods
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

  get renderTabs() {
    // renders tabs if listing visibility is set (auth tabs exist)
    // and user has NOT clicked "Sign in with other"
    if (this.authTabs && !this.showOtherMethods) {
      return true;
    }
    return false;
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

    if (this.renderTabs) {
      // reset selected auth method to first tab
      this.handleAuthSelect('tab', null, 0);
    } else {
      // all methods render, reset dropdown
      this.selectedAuthMethod = 'token';
    }
  }

  @action
  async handleNamespaceUpdate(event) {
    // update query param
    this.args.handleNamespaceUpdate(event.target.value);
    // reset tabs
    this.authTabs = null;
    // fetch mounts for that namespace
    this.fetchMounts.perform(500);
  }

  fetchMounts = restartableTask(
    waitFor(async (wait) => {
      // task is `restartable` so if the user starts typing again,
      // it will cancel and restart from the beginning.
      if (wait) await timeout(wait);

      try {
        // clear ember data store before re-requesting.. :(
        this.store.unloadAll('auth-method');

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
        this.authTabs = null;
      }
    })
  );
}
