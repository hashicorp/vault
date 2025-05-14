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

  get authViewMode() {
    const { directLinkData, loginSettings, visibleMountsByType } = this.args;

    const hasBackupMethods = !!loginSettings?.backupTypes?.length;
    const hasDefaultType = !!loginSettings?.defaultType;
    const hasVisibleMounts = !!visibleMountsByType;
    const isDirectLink = !!directLinkData?.isVisibleMount;

    // Rendering alternate view
    if (this.showOtherMethods) {
      if (!directLinkData && hasBackupMethods) return 'backupTabs';
      return 'dropdown';
    }

    // Rendering default view
    if (isDirectLink) return 'directLink';
    // Login settings only render without a direct link
    if (!directLinkData && hasDefaultType) return 'loginSettingsDefault';
    if (!directLinkData && hasBackupMethods) return 'tabs';
    if (hasVisibleMounts) return 'tabs';
    return 'dropdown';
  }

  get canToggleViews() {
    return this.hasAlternateView() && !this.showOtherMethods;
  }

  // Reveals custom path input
  get showAdvancedSettings() {
    // token does not support custom paths
    if (this.selectedAuthMethod === 'token') return false;

    const { directLinkData, loginSettings } = this.args;
    // in most cases, if the selected method has visible mount(s) the UI should prefer those and hide "Advanced settings"
    const hasMounts = !!this.selectedAuthHasMounts;

    // showOtherMethods is the fallback view in most cases, so we always want to show the advanced settings toggle.
    // UNLESS viewing backup methods then just rely on mount status
    if (directLinkData) {
      // always show advanced settings if showOtherMethods is true, otherwise hide/show depending on mount visibility
      return this.showOtherMethods ? true : !hasMounts;
    }
    return loginSettings?.backupTypes.length ? !hasMounts : !hasMounts || this.showOtherMethods;
  }

  get backupMethodTabs() {
    const { loginSettings, visibleMountsByType } = this.args;

    if (!loginSettings?.backupTypes?.length) return null;

    const tabs: UnauthMountsByType = {};
    for (const type of loginSettings.backupTypes) {
      // adds visible mounts for each type, if they exist
      tabs[type] = visibleMountsByType?.[type] || null;
    }
    return tabs;
  }

  get visibleMountTypes() {
    return Object.keys(this.args.visibleMountsByType || {});
  }

  get selectedAuthHasMounts() {
    return this.visibleMountTypes.includes(this.selectedAuthMethod);
  }

  @action
  handleError(message: string) {
    this.errorMessage = message;
  }

  @action
  handleNamespaceUpdate(event: HTMLElementEvent<HTMLInputElement>) {
    this.args.handleNamespaceUpdate(event.target.value);
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

    const type = this.determineAuthType(this.showOtherMethods);
    this.setAuthType(type);
  }

  @action
  initializeState() {
    // First, set auth type
    const type = this.determineAuthType();
    this.setAuthType(type);

    // Depending on which method is selected, determine which view renders
    // (if alternate views exist)
    this.showOtherMethods = this.determineShowOtherMethods();
  }

  private determineAuthType(showOtherMethods = false): string {
    const { canceledMfaAuth, directLinkData, loginSettings } = this.args;

    if (showOtherMethods) {
      // If "other methods" view is shown, prioritize backup types
      return loginSettings?.backupTypes?.[0] || this.auth.getAuthType() || 'token';
    }

    return (
      // Prioritize canceledMfaAuth since it's triggered by user interaction.
      canceledMfaAuth ||
      // Next, check type from directLinkData as it's specified by the URL.
      directLinkData?.type ||
      // Then, if there are login settings, prioritize those.
      loginSettings?.defaultType ||
      loginSettings?.backupTypes?.[0] ||
      // If listing_visibility is configured, the first tab should be selected
      this.visibleMountTypes?.[0] ||
      // Finally, fall back to the most recently used auth method in localStorage.
      this.auth.getAuthType() ||
      // Token is the default otherwise
      'token'
    );
  }

  private determineShowOtherMethods(): boolean {
    const { loginSettings, directLinkData } = this.args;

    if (directLinkData && !directLinkData?.isVisibleMount) {
      return this.args.visibleMountsByType ? !this.selectedAuthHasMounts : false;
    }

    if (loginSettings) {
      return false;
    }
    if (this.args.visibleMountsByType) {
      // if selectedAuthMethod is a tab (visible mount), then showOtherMethods should be "false"
      return !this.selectedAuthHasMounts;
    }
    return false;
  }

  private hasAlternateView(): boolean {
    const { directLinkData, loginSettings, visibleMountsByType } = this.args;

    if (directLinkData && directLinkData?.isVisibleMount) return true;
    if (directLinkData && !directLinkData?.isVisibleMount && !visibleMountsByType) return false;

    if (loginSettings) {
      const hasDefault = !!loginSettings?.defaultType;
      const hasBackups = !!loginSettings?.backupTypes.length;
      // if login settings exist *both* a default and backups must be set to toggle views
      return hasBackups && hasDefault;
    }

    if (visibleMountsByType) return true;
    return false;
  }
}
