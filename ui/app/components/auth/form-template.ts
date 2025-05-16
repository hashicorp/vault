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
 * The route fetches the unauthenticated sys/internal/ui/mounts endpoint to check for visible mounts and re-requests it when the namespace input updates.
 *
 * 🔧 CONFIGURATION OVERVIEW:
 * Each view mode (see `FormView` enum below) has specific layout configurations. In some scenarios, the component supports toggling between a default view and an alternate view.
 *
 * 📋 [DROPDOWN] (default view)
 *   ▸ All supported auth methods show in a dropdown.
 *   ▸ No alternate view.
 *
 * 🗂️ [TABS] (unauth mount tabs)
 *   ▸ Groups visible mounts (`listing_visibility="unauth"`) by type and displays as tabs.
 *   ▸ Alternate view: full dropdown of all methods.
 *
 * 🔗 [DIRECT_LINK] (via `?with=` query param)
 *   ▸ If the param references a visible mount, that method renders by default and the mount path is assumed.
 *     ↳ Alternate view: full dropdown.
 *   ▸ If the param references a method type (legacy behavior), the method is preselected in the dropdown or its tab is selected.
 *     ↳ Alternate view: if other methods have visible mounts, the form can toggle between tabs and dropdown. The initial view depends on whether the chosen type is a tab.
 *
 * 🏢 *Enterprise-only login customizations*
 *   ▸ A namespace can define a default method [LOGIN_SETTINGS_DEFAULT] and/or preferred methods (i.e. "backups") [LOGIN_SETTINGS_TABS].
 *     ✎ Both set:
 *       ▸ Default method shown initially.
 *       ▸ Alternate view: preferred methods in tab layout.
 *     ✎ Only one set:
 *       ▸ No alternate view.
 *
 * 🔁 Advanced settings toggle reveals the custom path input:
 *   🚫 No visible mounts:
 *     ▸ UI defaults to method type as path.
 *     ▸ "Advanced settings" shows a path input.
 *   1️⃣ One visible mount:
 *     ▸ Path is assumed and hidden.
 *   🔀 Multiple visible mounts:
 *     ▸ Path dropdown is shown.
 *
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
  loginSettings: { defaultType: string; backupTypes: string[] | null }; // enterprise only
  namespaceQueryParam: string;
  oidcProviderQueryParam: string;
  onSuccess: CallableFunction;
  visibleMountsByType: UnauthMountsByType;
}

enum FormView {
  DIRECT_LINK = 'directLink',
  DROPDOWN = 'dropdown',
  LOGIN_SETTINGS_DEFAULT = 'loginSettingsDefault',
  LOGIN_SETTINGS_TABS = 'loginSettingsTabs',
  TABS = 'tabs',
}

export default class AuthFormTemplate extends Component<Args> {
  view = FormView;

  @service declare readonly auth: AuthService;
  @service declare readonly flags: FlagsService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  // true → "Back" button renders, false → "Sign in with other methods→" (only if an alternate view exists)
  @tracked showAlternateLoginView = false;

  @tracked selectedAuth = ''; // determines which form renders
  @tracked errorMessage = '';

  supportedAuthTypes: string[];
  backupMethodTabs: UnauthMountsByType | null;
  namespaceInput: string;
  visibleMountTypes: string[];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.backupMethodTabs = this.setBackupTabs();
    this.visibleMountTypes = this.setMountTypes();
    this.namespaceInput = this.setNamespaceInput();
    this.supportedAuthTypes = supportedTypes(this.version.isEnterprise);
  }

  get currentFormView() {
    if (this.showAlternateLoginView) return this.alternateView;

    if (this.isDirectLink()) return this.view.DIRECT_LINK;

    // Direct link overrides login settings so these checks come after
    if (this.hasLoginSettingsDefault()) return this.view.LOGIN_SETTINGS_DEFAULT;
    if (this.hasLoginSettingsTabs()) return this.view.LOGIN_SETTINGS_TABS;
    if (this.hasVisibleMounts()) return this.view.TABS;

    return this.view.DROPDOWN;
  }

  get alternateView() {
    const { directLinkData, loginSettings } = this.args;

    if (directLinkData) {
      return this.methodHasMounts(this.selectedAuth) || this.hasVisibleMounts() ? this.view.DROPDOWN : null;
    }

    if (loginSettings) {
      const hasDefault = !!loginSettings?.defaultType;
      const hasBackups = !!loginSettings?.backupTypes;
      // Both default and backups must be set for an alternate view to exist
      return hasBackups && hasDefault ? this.view.LOGIN_SETTINGS_TABS : null;
    }

    if (this.hasVisibleMounts()) {
      return this.view.DROPDOWN;
    }

    return null;
  }

  get formComponent() {
    const { selectedAuth } = this;
    // isSupported means there is a component file defined for that auth type
    const isSupported = this.supportedAuthTypes.includes(selectedAuth);
    const formFile = () => (['oidc', 'jwt'].includes(selectedAuth) ? 'oidc-jwt' : selectedAuth);
    const component = isSupported ? formFile() : 'base';

    // an Auth::Form::<Type> component exists for each method in supported-login-methods
    return `auth/form/${component}`;
  }

  get hideAdvancedSettings() {
    // Token does not support custom paths
    if (this.selectedAuth === 'token') return true;

    switch (this.currentFormView) {
      case this.view.DROPDOWN:
        // Always show for dropdown mode
        return false;
      case this.view.DIRECT_LINK:
        // For direct links, always show advanced settings when rendering the "other" view.
        // Otherwise hide/show depending on mount visibility
        return !this.showAlternateLoginView || this.methodHasMounts(this.selectedAuth);
      default:
        // For remaining scenarios, hide "Advanced settings" if the selected method has visible mount(s)
        return this.methodHasMounts(this.selectedAuth);
    }
  }

  @action
  initializeState() {
    // First, set auth type
    const type = this.determineAuthType();
    this.setAuthType(type);

    // Depending on which method is selected, determine which view renders
    // (if alternate views exist)
    this.showAlternateLoginView = this.determineShowOtherMethods();
  }

  @action
  setAuthType(authType: string) {
    this.selectedAuth = authType;
  }

  @action
  setTypeFromDropdown(event: HTMLElementEvent<HTMLInputElement>) {
    this.selectedAuth = event.target.value;
  }

  @action
  toggleView() {
    this.showAlternateLoginView = !this.showAlternateLoginView;

    const type = this.defaultTypeByView();
    this.setAuthType(type);
  }

  private defaultTypeByView() {
    switch (this.currentFormView) {
      case this.view.LOGIN_SETTINGS_TABS:
        return this.args.loginSettings.backupTypes?.[0] as string;
      case this.view.TABS:
        return this.visibleMountTypes[0] as string;
      default:
        return this.determineAuthType();
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

  private determineAuthType(): string {
    const { canceledMfaAuth, directLinkData, loginSettings } = this.args;

    if (this.showAlternateLoginView) {
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
    const { loginSettings, directLinkData, visibleMountsByType } = this.args;
    if (!this.alternateView) return false;

    if (directLinkData && !directLinkData?.isVisibleMount) {
      return !this.methodHasMounts(this.selectedAuth);
    }

    if (loginSettings) return false;

    // if selectedAuth is a tab (visible mount), then showAlternateLoginView should be "false"
    return visibleMountsByType ? !this.methodHasMounts(this.selectedAuth) : false;
  }

  private hasVisibleMounts(): boolean {
    return !!this.args.visibleMountsByType;
  }

  private isDirectLink(): boolean {
    return !!this.args.directLinkData?.isVisibleMount;
  }

  private hasLoginSettingsDefault(): boolean {
    return !this.args.directLinkData && !!this.args.loginSettings?.defaultType;
  }

  private hasLoginSettingsTabs(): boolean {
    return !this.args.directLinkData && !!this.args.loginSettings?.backupTypes;
  }

  private setBackupTabs() {
    const { loginSettings, visibleMountsByType } = this.args;

    if (!loginSettings?.backupTypes) return null;

    const tabs: UnauthMountsByType = {};
    for (const type of loginSettings.backupTypes) {
      // adds visible mounts for each type, if they exist
      tabs[type] = visibleMountsByType?.[type] || null;
    }
    return tabs;
  }

  private methodHasMounts(authMethod: string) {
    return this.visibleMountTypes.includes(authMethod);
  }

  private setMountTypes() {
    return Object.keys(this.args.visibleMountsByType || {});
  }

  private setNamespaceInput() {
    const namespaceQueryParam = this.args.namespaceQueryParam;
    if (this.flags.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/`
      // so prefill just the relative path in the namespace input
      const path = getRelativePath(namespaceQueryParam, this.flags.hvdManagedNamespaceRoot);
      return path ? `/${path}` : '';
    }
    return namespaceQueryParam;
  }
}
