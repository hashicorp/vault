/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { supportedTypes } from 'vault/utils/supported-login-methods';

import type VersionService from 'vault/services/version';
import type ClusterModel from 'vault/models/cluster';
import type { UnauthMountsByType } from 'vault/vault/auth/form';
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
 * @param {object | null} alternateView - if an alternate view exists, this is the `FormView` (see interface below) data to render that view.
 * @param {string} canceledMfaAuth - saved auth type from a cancelled mfa verification
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {object} defaultView - The `FormView` (see the interface below) data to render the initial view.
 * @param {object} initialFormState - sets selectedAuthMethod and showAlternateView based on the login form configuration computed in parent component
 * @param {function} onSuccess - callback after the initial authentication request, if an mfa_requirement exists the parent renders the mfa form otherwise it fires the authSuccess action in the auth controller and handles transitioning to the app
 * @param {array} visibleMountTypes - array of auth method types that have mounts with listing_visibility="unauth"
 * */

interface Args {
  alternateView: FormView | null;
  cluster: ClusterModel;
  defaultView: FormView;
  initialFormState: { initialAuthType: string; showAlternate: boolean };
  onSuccess: CallableFunction;
  visibleMountTypes: string[];
}

interface FormView {
  view: string; // "dropdown" or "tabs"
  tabData: UnauthMountsByType | null; // tabs to render if view = "tabs"
}

export default class AuthFormTemplate extends Component<Args> {
  @service declare readonly version: VersionService;

  supportedAuthTypes: string[];

  @tracked errorMessage = '';
  @tracked selectedAuthMethod = '';
  // true → "Back" button renders, false → "Sign in with other methods→" renders if an alternate view exists
  @tracked showAlternateView = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { initialAuthType, showAlternate } = this.args.initialFormState;
    this.selectedAuthMethod = initialAuthType;
    this.showAlternateView = showAlternate;
    this.supportedAuthTypes = supportedTypes(this.version.isEnterprise);
  }

  get tabData() {
    if (this.showAlternateView) return this.args?.alternateView?.tabData;
    return this.args?.defaultView?.tabData;
  }

  get formComponent() {
    const { selectedAuthMethod } = this;
    // isSupported means there is a component file defined for that auth type
    const isSupported = this.supportedAuthTypes.includes(selectedAuthMethod);
    const formFile = () => (['oidc', 'jwt'].includes(selectedAuthMethod) ? 'oidc-jwt' : selectedAuthMethod);
    const component = isSupported ? formFile() : 'base';

    // an Auth::Form::<Type> component exists for each method in supported-login-methods
    return `auth/form/${component}`;
  }

  get hideAdvancedSettings() {
    // Token does not support custom paths
    if (this.selectedAuthMethod === 'token') return true;

    // Always show for dropdown mode
    if (!this.tabData) return false;

    // For remaining scenarios, hide "Advanced settings" if the selected method has visible mount(s)
    return this.args.visibleMountTypes?.includes(this.selectedAuthMethod);
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
    this.showAlternateView = !this.showAlternateView;
    const firstAuthTab = Object.keys(this.tabData || {})[0];
    const type = firstAuthTab || this.args.initialFormState.initialAuthType;
    this.setAuthType(type);
  }

  @action
  handleError(message: string) {
    this.errorMessage = message;
  }
}
