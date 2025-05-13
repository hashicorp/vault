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
import type { UnauthMountsByType, AuthTabMountData } from 'vault/vault/auth/form';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module Auth::FormTemplate
 * This component manages the layout and display logic of the login form. When auth type changes the component dynamically renders the corresponding form.
 * The route fetches the unauthenticated sys/internal/ui/mounts endpoint to check if any mounts have `listing_visibility="unauth"`.
 * The endpoint is re-requested anytime the namespace input updates.
 *
 * 🔧 CONFIGURATION OPTIONS:
 * Outlined below are the different configuration options, in some scenarios there are two views to toggle between:
 * - The initial view (shown by default)
 * - The alternate methods view (displayed when the user clicks "Sign in with other methods →")
 *
 * 📋 Standard dropdown view (no form configurations):
 *   ▸ Dropdown lists ALL auth methods supported by the UI
 *   ▸ Alternate view: None
 *
 * 🗂️ Unauth mount tabs:
 *   ▸ Auth mounts (not methods) with `listing_visibility="unauth"` are grouped by type and render as tabs.
 *   ▸ Alternate view: Dropdown with all methods
 *
 * 🔗 Direct link (auth URL contains the `?with=` query param):
 *   ▸ If param references a visible mount, the corresponding method type renders and the mount path is assumed for login
 *      ↳ Alternate view: Dropdown with all methods
 *   ▸ Param references a type (backward compatibility)
 *      ↳ Type is selected in either dropdown or as tab, depending on listing visibility configs
 *
 * 🏢 *Enterprise-only* Login customizations:
 *   ▸ A namespace can have a default auth method and/or preferred (backup) methods set. Preferred methods display as tabs.
 *    ✎ Default + preferred is set
 *        ↳ Default method displays by default
 *        ↳ Alternate view: Preferred methods as tabs
 *    ✎ Only default OR only preferred methods selected
 *        ↳ Alternate view: None
 *   ▸ The "path" input depends on the number of visible mounts for a method:
 *    🚫 No visible mounts
 *        ↳ The UI assumes the default path (which is the auth method type)
 *        ↳ "Advanced settings" toggle reveals an input for an optional custom path
 *    1️⃣ One visible mount
 *        ↳ Path renders in a hidden input and it is assumed for login
 *    🔀 Multiple visible mounts
 *        ↳ Dropdown lists with all paths.
 *
 * @param {string} canceledMfaAuth - saved auth type from a cancelled mfa verification
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {object} directLinkData - mount data built from the "with" query param. If param is a mount path and maps to a visible mount, the login form defaults to this mount. Otherwise the form preselects the passed auth type.
 * @param {function} handleNamespaceUpdate - callback task that passes user input to the controller and updates the namespace query param in the url
 * @param {object} loginSettings -
 * @param {string} namespaceQueryParam - namespace query param from the url
 * @param {string} oidcProviderQueryParam - oidc provider query param, set in url as "?o=someprovider". if present, disables the namespace input
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 * @param {object} visibleMountsByType - auth methods to render as tabs, contains mount data for any mounts with listing_visibility="unauth"
 *
 * */

interface Args {
  canceledMfaAuth: string;
  cluster: ClusterModel;
  directLinkData: (AuthTabMountData & { isVisibleMount: boolean }) | null;
  handleNamespaceUpdate: CallableFunction;
  loginSettings: { defaultType: string; backupTypes: string[] };
  namespaceQueryParam: string;
  oidcProviderQueryParam: string;
  onSuccess: CallableFunction;
  visibleMountsByType: UnauthMountsByType;
}

export default class AuthFormTemplate extends Component<Args> {
  @service declare readonly auth: AuthService;
  @service declare readonly flags: FlagsService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  // true → "Back" button renders, false → "Sign in with other methods→" (only if an alternate view exists)
  @tracked showOtherMethods = false;

  @tracked selectedAuthMethod = ''; // determines which form renders
  @tracked errorMessage = '';

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
    return (
      // Prioritize canceledMfaAuth since it's triggered by user interaction.
      this.args.canceledMfaAuth ||
      // Next, check type from directLinkData as it's specified by the URL.
      this.args.directLinkData?.type ||
      // Then, if there is a default method set we want to prioritize that.
      this.args.loginSettings?.defaultType ||
      // Finally, fall back to the most recently used auth method in localStorage.
      this.auth.getAuthType()
    );
  }

  // The "standard" login view is a dropdown listing all auth methods.
  // This getter determines if an alternate configuration exists and whether it should render.
  get showCustomAuthOptions() {
    const hasNonStandardConfiguration =
      this.args?.directLinkData?.isVisibleMount || this.args.loginSettings || this.args.visibleMountsByType;
    // Show if a configuration exists and the user has NOT clicked "Sign in with other methods →"
    return hasNonStandardConfiguration && !this.showOtherMethods;
  }

  get canToggleViews() {
    const { loginSettings } = this.args;
    if (loginSettings) {
      const hasBackups = loginSettings?.backupTypes.length;
      const hasDefault = loginSettings.defaultType;
      // if login customizations exist, users can only toggle form views if *both* a default and backups have been set
      return hasBackups && hasDefault;
    }
    return true;
  }

  get showAdvancedSettings() {
    // token does not support custom paths
    if (this.selectedAuthMethod === 'token') return false;

    const { loginSettings, directLinkData, visibleMountsByType } = this.args;
    // for the most part, if a method has a visible mount the UI should assume to use that path
    const isVisibleMount = !!visibleMountsByType?.[this.selectedAuthMethod];
    // but if ONLY listing visibility is configured and the user has toggled to view the "other" methods, the UI allows inputting a custom path
    const hasToggledFromVisibleMountTabs =
      !loginSettings && !directLinkData && visibleMountsByType && this.showOtherMethods;
    return !isVisibleMount || hasToggledFromVisibleMountTabs;
  }

  get backupMethodTabs() {
    return this.args.loginSettings.backupTypes?.reduce((obj, type) => {
      const mountData = this.args.visibleMountsByType?.[type];
      obj[type] = mountData || null;
      return obj;
    }, {} as UnauthMountsByType);
  }

  get visibleMountTypes() {
    return Object.keys(this.args.visibleMountsByType || {});
  }

  @action
  initializeState() {
    // SET AUTH TYPE
    const type =
      this.preselectedType ||
      this.args.loginSettings?.defaultType ||
      this.args.loginSettings?.backupTypes?.[0] ||
      this.visibleMountTypes?.[0] ||
      'token';
    this.setAuthType(type);

    // DETERMINES INITIAL RENDER: custom selection (direct link or tabs) vs dropdown
    if (this.args.loginSettings) {
      this.showOtherMethods = false;
    } else if (this.args.visibleMountsByType) {
      // render tabs if selectedAuthMethod is one, otherwise render dropdown (i.e. showOtherMethods = false)
      this.showOtherMethods = this.visibleMountTypes.includes(this.selectedAuthMethod) ? false : true;
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

    const { loginSettings } = this.args;
    const hasDefault = !!loginSettings?.defaultType;

    if (loginSettings && this.showOtherMethods) {
      // user has clicked "Sign in with other", select first tab
      // toggle button is hidden if there are no backups
      // so if we've gotten to this conditional, this form is rendering backup tabs
      const firstTab = this.args.loginSettings.backupTypes[0] as string;
      this.setAuthType(firstTab);
    } else if (loginSettings && !this.showOtherMethods) {
      // this is the initial view of login form
      if (hasDefault) {
        this.setAuthType(loginSettings?.defaultType);
      } else {
        // if no default, but customizations exist then the backup tabs render here
        const firstTab = this.args.loginSettings.backupTypes[0] as string;
        this.setAuthType(firstTab);
      }
    } else if (this.showOtherMethods) {
      // all methods render, reset dropdown
      this.selectedAuthMethod = this.preselectedType || 'token';
    } else {
      const firstTab = this.visibleMountTypes[0] as string;
      this.setAuthType(firstTab);
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
