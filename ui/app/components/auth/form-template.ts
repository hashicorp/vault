/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { supportedTypes } from 'vault/utils/supported-login-methods';
import { getRelativePath } from 'core/utils/sanitize-path';

import type AuthService from 'vault/vault/services/auth';
import type FlagsService from 'vault/services/flags';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type ClusterModel from 'vault/models/cluster';
import type { AuthTabData, AuthTabMountData } from 'vault/vault/auth/form';
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
 * @param {object} authTabData - auth methods to render as tabs, contains mount data for any mounts with listing_visibility="unauth"
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} handleNamespaceUpdate - callback task that passes user input to the controller and updates the namespace query param in the url
 * @param {boolean} hasVisibleAuthMounts - whether or not any mounts have been tuned with listing_visibility="unauth"
 * @param {string} namespaceQueryParam - namespace query param from the url
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 * @param {string} canceledMfaAuth - saved auth type from a cancelled mfa verification
 * @param {object} directLinkData - mount data built from the "with" query param. If param is a mount path and maps to a visible mount, the login form defaults to this mount. Otherwise the form preselects the passed auth type.
 *
 * */

interface Args {
  authTabData: AuthTabData;
  cluster: ClusterModel;
  handleNamespaceUpdate: CallableFunction;
  hasVisibleAuthMounts: boolean;
  namespaceQueryParam: string;
  onSuccess: CallableFunction;
  canceledMfaAuth: string;
  directLinkData: (AuthTabMountData & { hasMountData: boolean }) | null;
}

export default class AuthFormTemplate extends Component<Args> {
  @service declare readonly auth: AuthService;
  @service declare readonly flags: FlagsService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  // true → "Back" button renders, false → "Sign in with other methods→" renders if customizations exist
  @tracked showOtherMethods = false;

  // auth login variables
  @tracked selectedAuthMethod = '';
  @tracked errorMessage = '';

  get authTabTypes() {
    const visibleMounts = this.args.authTabData;
    return visibleMounts ? Object.keys(visibleMounts) : [];
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

  get preselectedType() {
    // Prioritize canceledMfaAuth since it's triggered by user interaction.
    // Next, check type from directLinkData as it's specified by the URL.
    // Finally, fall back to the most recently used auth method in localStorage.
    return this.args.canceledMfaAuth || this.args.directLinkData?.type || this.auth.getAuthType();
  }

  // The "standard" selection is a dropdown listing all auth methods.
  // This getter determines whether to render an alternative view (e.g., tabs or a preferred mount).
  // If `true`, the "Sign in with other methods →" link is shown.
  get showCustomAuthOptions() {
    const hasLoginCustomization = this.args?.directLinkData?.hasMountData || this.args.hasVisibleAuthMounts;
    // Show if customization exists and the user has NOT clicked "Sign in with other methods →"
    return hasLoginCustomization && !this.showOtherMethods;
  }

  @action
  initializeState() {
    // SET AUTH TYPE
    if (this.preselectedType) {
      this.setAuthType(this.preselectedType);
    } else {
      // if nothing has been preselected, select first tab or set to 'token'
      const authType = this.args.hasVisibleAuthMounts ? (this.authTabTypes[0] as string) : 'token';
      this.setAuthType(authType);
    }

    // DETERMINES INITIAL RENDER: custom selection (direct link or tabs) vs dropdown
    if (this.args.hasVisibleAuthMounts) {
      // render tabs if selectedAuthMethod is one, otherwise render dropdown (i.e. showOtherMethods = false)
      this.showOtherMethods = this.authTabTypes.includes(this.selectedAuthMethod) ? false : true;
    } else {
      this.showOtherMethods = false;
    }
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

    if (this.showCustomAuthOptions) {
      const firstTab = this.authTabTypes[0] as string;
      this.setAuthType(firstTab);
    } else {
      // all methods render, reset dropdown
      this.selectedAuthMethod = this.preselectedType || 'token';
    }
  }

  @action
  handleError(message: string) {
    this.errorMessage = message;
  }

  @action
  handleNamespaceUpdate(event: HTMLElementEvent<HTMLInputElement>) {
    this.args.handleNamespaceUpdate(event.target.value);
  }
}
