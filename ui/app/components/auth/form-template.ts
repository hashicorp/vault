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
import type VersionService from 'vault/services/version';
import type ClusterModel from 'vault/models/cluster';
import type { UnauthMountsByType } from 'vault/vault/auth/form';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module Auth::FormTemplate
 * This component manages which auth type should initially be selected and dynamically rendering the corresponding form.
 * If an `alternateView` exists, it also handles toggling between the two views and updating the selectedAuthMethod.
 *
 * @param {object} [alternateView] - If it exists, "Sign in with other methods" renders and displays either a dropdown or tab data
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {object} defaultView - The initial view of the login form, either tabs with the auth type or a dropdown with all methods
 * @param {string} preselectedType - saved auth type from a cancelled mfa verification
 * @param {function} handleNamespaceUpdate - callback task that passes user input to the controller and updates the namespace query param in the url
 * @param {string} namespaceQueryParam - namespace query param from the url
 * @param {string} oidcProviderQueryParam - oidc provider query param, set in url as "?o=someprovider". if present, disables the namespace input
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 *
 * */

interface Args {
  alternateView: FormState | null;
  cluster: ClusterModel;
  defaultView: FormState;
  handleNamespaceUpdate: CallableFunction;
  namespaceQueryParam: string;
  oidcProviderQueryParam: string;
  onSuccess: CallableFunction;
  initialFormState: () => { initialAuthType: string; showAlternate: boolean };
  preselectedType: string;
  visibleMountTypes: string[];
}

interface FormState {
  view: string;
  tabData: UnauthMountsByType | null;
}

export default class AuthFormTemplate extends Component<Args> {
  @service declare readonly auth: AuthService;
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  @tracked selectedAuth = ''; // determines which form renders
  // true → "Back" button renders, false → "Sign in with other methods→" renders (if an alternate view exists)
  @tracked showAlternateLoginView = false;
  @tracked errorMessage = '';

  supportedAuthTypes: string[];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.supportedAuthTypes = supportedTypes(this.version.isEnterprise);
    const { initialAuthType, showAlternate } = this.args.initialFormState();
    this.showAlternateLoginView = showAlternate;
    this.selectedAuth = initialAuthType;
  }

  get firstAuthTab() {
    const tabs = Object.keys(this.tabData || {});
    return tabs?.[0] ?? '';
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

    // Always show for dropdown mode
    if (!this.tabData) return false;

    // For remaining scenarios, hide "Advanced settings" if the selected method has visible mount(s)
    return this.args.visibleMountTypes?.includes(this.selectedAuth);
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

  get tabData() {
    if (this.showAlternateLoginView) return this.args?.alternateView?.tabData;
    return this.args?.defaultView?.tabData;
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

    const type = this.firstAuthTab || this.determineAuthType();

    this.setAuthType(type);
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
    const { preselectedType } = this.args;

    if (this.showAlternateLoginView) {
      // If rendering alternate view, prioritize backup types
      return this.firstAuthTab || this.auth.getAuthType() || 'token';
    }

    return (
      // Prioritize preselectedType since it's set by canceled MFA validation or from the URL
      preselectedType ||
      // Then, set as first tab which is either the default type or first unauth mount
      this.firstAuthTab ||
      // Finally, fall back to the most recently used auth method in localStorage.
      this.auth.getAuthType() ||
      // Token is the default otherwise
      'token'
    );
  }
}
