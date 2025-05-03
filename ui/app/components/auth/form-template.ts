/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
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
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} handleNamespaceUpdate - callback task that passes user input to the controller and updates the namespace query param in the url
 * @param {string} namespaceQueryParam - namespace query param from the url
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 * @param {string} preselectedAuthType - auth type to preselect to in login form, set from either local storage (last method used to log in) or on canceled mfa validation
 * @param {object} visibleAuthMounts - mount data from auth mounts tuned with listing_visibility="unauth"
 *
 * */

interface Args {
  cluster: ClusterModel;
  handleNamespaceUpdate: CallableFunction;
  namespaceQueryParam: string;
  onSuccess: CallableFunction;
  preselectedAuthType: string; // set by local storage or canceled MFA validation
  visibleAuthMounts: VisibleAuthMounts;
}

interface VisibleAuthMounts {
  [key: string]: {
    description: string;
    type: string;
  };
}

interface AuthTabData {
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
  @tracked showOtherMethods = false;

  // auth login variables
  @tracked selectedAuthMethod = '';
  @tracked errorMessage = '';

  displayName = (type: string) => {
    const displayName = ALL_LOGIN_METHODS?.find((t) => t.type === type)?.displayName;
    return displayName || type;
  };

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    // ensures args have settled before setting form state
    setTimeout(() => this.initializeState(), 0);
  }

  initializeState() {
    // SET AUTH TYPE
    if (!this.args.preselectedAuthType) {
      // if nothing has been preselected, select first tab or set to 'token'
      const authType = this.authTabData ? (this.authTabTypes[0] as string) : 'token';
      this.setAuthType(authType);
    } else {
      // there is a preselected type, set is as the selectedAuthType
      this.setAuthType(this.args.preselectedAuthType);
    }

    // INITIALLY RENDER TABS OR DROPDOWN
    // selectedAuthMethod is a tab, render tabs
    // otherwise render dropdown (i.e. showOtherMethods = false)
    if (this.authTabData) {
      this.showOtherMethods = this.authTabTypes.includes(this.selectedAuthMethod) ? false : true;
    } else {
      this.showOtherMethods = false;
    }
  }

  get authTabTypes() {
    return this.authTabData ? Object.keys(this.authTabData) : [];
  }

  get authTabData() {
    if (this.args.visibleAuthMounts) {
      const authMounts = this.args.visibleAuthMounts;
      return Object.entries(authMounts).reduce((obj, [path, mountData]) => {
        const { type } = mountData;
        obj[type] ??= []; // if an array doesn't already exist for that type, create it
        obj[type].push({ path, ...mountData });
        return obj;
      }, {} as AuthTabData);
    }
    return null;
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
    const namespaceQueryParam = this.args.namespaceQueryParam;
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
    if (this.authTabData && !this.showOtherMethods) {
      return true;
    }
    return false;
  }

  @action
  setAuthType(authType: string) {
    this.selectedAuthMethod = authType;
  }

  @action
  setTypeFromDropdown(event: HTMLElementEvent<HTMLInputElement>) {
    this.selectedAuthMethod = event.target.value;
  }

  @action
  toggleView() {
    this.showOtherMethods = !this.showOtherMethods;

    if (this.renderTabs) {
      const firstTab = this.authTabTypes[0] as string;
      this.setAuthType(firstTab);
    } else {
      // all methods render, reset dropdown
      this.selectedAuthMethod = this.args.preselectedAuthType || 'token';
    }
  }

  @action
  handleError(message: string) {
    this.errorMessage = message;
  }

  @action
  async handleNamespaceUpdate(event: HTMLElementEvent<HTMLInputElement>) {
    this.args.handleNamespaceUpdate(event.target.value);
  }
}
