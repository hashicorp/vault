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
import { getRelativePath } from 'core/utils/sanitize-path';

import type FlagsService from 'vault/services/flags';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type ClusterModel from 'vault/models/cluster';
import type { HTMLElementEvent } from 'vault/forms';

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
 *
 * @param {string} wrappedToken - Query param value of a wrapped token that can be used to login when added directly to the URL via the "wrapped_token" query param
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} handleNamespaceUpdate - callback task that passes user input to the controller and updates the namespace query param in the url
 * @param {string} namespace - namespace query param from the url
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 *
 * */

interface Args {
  wrappedToken: string;
  cluster: ClusterModel;
  handleNamespaceUpdate: CallableFunction;
  namespace: string;
  onSuccess: CallableFunction;
}

interface AuthTabs {
  // key is the auth method type
  [key: string]: MountData[];
}

interface MountData {
  path: string;
  type: string;
  description?: string;
  config?: object | null;
}

export default class AuthFormTemplate extends Component<Args> {
  @service declare readonly flags: FlagsService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  // form display logic
  @tracked authTabs: AuthTabs | null = null;
  @tracked showOtherMethods = false;

  // auth login variables
  @tracked selectedAuthMethod = 'token';
  @tracked errorMessage = '';

  displayName = (type: string) => {
    const displayName = ALL_LOGIN_METHODS?.find((t) => t.type === type)?.displayName;
    return displayName || type;
  };

  constructor(owner: unknown, args: Args) {
    super(owner, args);
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
    const namespaceQueryParam = this.args.namespace;
    if (this.flags.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/`
      // so prefill just the relative path in the namespace input
      const path = getRelativePath(namespaceQueryParam, this.flags.hvdManagedNamespaceRoot);
      return path ? `/${path}` : '';
    }
    return namespaceQueryParam;
  }

  get renderTabs() {
    // renders tabs if listing visibility is set (auth tabs exist)
    // and user has NOT clicked "Sign in with other"
    if (this.authTabs && !this.showOtherMethods) {
      return true;
    }
    return false;
  }

  get selectedTabIndex() {
    if (this.authTabs) {
      return Object.keys(this.authTabs).indexOf(this.selectedAuthMethod);
    }
    return 0;
  }

  setAuthTypeFromTab(idx: number) {
    const authTypes = this.authTabs ? Object.keys(this.authTabs) : [];
    this.selectedAuthMethod = authTypes[idx] || '';
  }

  @action
  handleAuthSelect(element: string, event: HTMLElementEvent<HTMLInputElement> | null, idx: number) {
    if (element === 'tab') {
      this.setAuthTypeFromTab(idx);
    } else if (event?.target?.value) {
      this.selectedAuthMethod = event.target.value;
    }
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
  handleError(message: string) {
    this.errorMessage = message;
  }

  @action
  handleNamespaceUpdate(event: HTMLElementEvent<HTMLInputElement>) {
    // update query param
    this.args.handleNamespaceUpdate(event.target.value);
    // reset tabs
    this.authTabs = null;
    // fetch mounts for that namespace
    this.fetchMounts.perform(500);
  }

  fetchMounts = restartableTask(
    waitFor(async (wait = 0) => {
      // task is `restartable` so if the user starts typing again,
      // it will cancel and restart from the beginning.
      if (wait) await timeout(wait);

      try {
        // clear ember data store before re-requesting.. :(
        this.store.unloadAll('auth-method');

        // unauthMounts are tuned with listing_visibility="unauth"
        const unauthMounts = await this.store.findAll('auth-method', {
          adapterOptions: {
            unauthenticated: true,
          },
        });

        if (unauthMounts.length !== 0) {
          this.authTabs = unauthMounts.reduce((obj: AuthTabs, m) => {
            // serialize the ember data model into a regular ol' object
            const mountData = m.serialize();
            const methodType = mountData.type;
            if (!Object.keys(obj).includes(methodType)) {
              // create a new empty array for that type
              obj[methodType] = [];
            }

            if (Array.isArray(obj[methodType])) {
              // push mount data into corresponding type's array
              obj[methodType].push(mountData);
            }

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
