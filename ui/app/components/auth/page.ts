/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import type { AuthResponse, AuthResponseWithMfa } from 'vault/vault/services/auth';
import type { UnauthMountsByType, UnauthMountsResponse } from 'vault/vault/auth/form';
import type ClusterModel from 'vault/models/cluster';
import type CspEventService from 'vault/services/csp-event';

/**
 * @module AuthPage
 * The Auth::Page receives configuration settings from the route's model hook and determines the possible form states. It also renders the Auth::FormTemplate or MFA component if an
 * mfa validation is returned from the auth request. The model hook refreshes when the namespace input updates.
 *
 * 🔧 CONFIGURATION OVERVIEW:
 * Each view mode (see `FormView` enum below) has specific layout configurations. In some scenarios, the component supports toggling between a default view and an alternate view.
 *
 * 📋 [DROPDOWN] (default view)
 *   ▸ All supported auth methods show in a dropdown.
 *   ▸ No alternate view.
 *
 * 🗂️ [TABS] (visible (unauth) mount tabs)
 *   ▸ Groups visible mounts (`listing_visibility="unauth"`) by type and displays as tabs.
 *   ▸ Alternate view: full dropdown of all methods.
 *
 * 🔗 [DIRECT_LINK] (via `?with=` query param)
 *   ▸ If the param references a visible mount, that method renders by default and the mount path is assumed.
 *     ↳ Alternate view: full dropdown.
 *   ▸ If the param references a method type (legacy behavior), the method is preselected in the dropdown or its tab is selected.
 *     ↳ Alternate view: if other methods have visible mounts, the form can toggle between tabs and dropdown. The initial view depends on whether the chosen type is a tab.
 *
 * 🏢 *Enterprise-only login settings*
 *   ▸ A namespace can define a default method and/or preferred methods (i.e. "backups") and enable child namespaces to inherit these preferences.
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
 * @example
 * <Auth::Page
 *  @cluster={{this.model.clusterModel}}
 *  @namespaceQueryParam={{this.namespaceQueryParam}}
 *  @oidcProviderQueryParam={{this.oidcProvider}}
 *  @onAuthSuccess={{action "authSuccess"}}
 *  @onNamespaceUpdate={{perform this.updateNamespace}}
 *  @visibleAuthMounts={{this.model.visibleAuthMounts}}
 *  @directLinkData={{this.model.directLinkData}}
 * />
 *
 * @param {object} cluster - the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {object} directLinkData - mount data built from the "with" query param. If param is a mount path and maps to a visible mount, the login form defaults to this mount. Otherwise the form preselects the passed auth type.
 * @param {object} loginSettings - * enterprise only * login settings configured for the namespace
 * @param {string} namespaceQueryParam - namespace to login with, updated by typing in to the namespace input
 * @param {string} oidcProviderQueryParam - oidc provider query param, set in url as "?o=someprovider"
 * @param {function} onAuthSuccess - callback task in controller that receives the auth response (after MFA, if enabled) when login is successful
 * @param {function} onNamespaceUpdate - callback task that passes user input to the controller to update the login namespace in the url query params
 * @param {object} visibleAuthMounts - response from unauthenticated request to sys/internal/ui/mounts which returns mount paths tuned with `listing_visibility="unauth"`. keys are the mount path, values are mount data such as "type" or "description," if it exists
 * */

export const CSP_ERROR =
  "This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.";

interface Args {
  cluster: ClusterModel;
  directLinkData: { type: string; path?: string } | null; // if "path" key is present then mount data is visible
  loginSettings: { defaultType: string; backupTypes: string[] | null }; // enterprise only
  onAuthSuccess: CallableFunction;
  visibleAuthMounts: UnauthMountsResponse;
}

interface MfaAuthData {
  mfa_requirement: object;
  path: string;
  selectedAuth: string;
}

enum FormView {
  DROPDOWN = 'dropdown',
  TABS = 'tabs',
}

export default class AuthPage extends Component<Args> {
  @service('csp-event') declare readonly csp: CspEventService;

  @tracked canceledMfaAuth = '';
  @tracked mfaAuthData: MfaAuthData | null = null;
  @tracked mfaErrors = '';

  get cspError() {
    const isStandby = this.args.cluster.standby;
    const hasConnectionViolations = this.csp.connectionViolations.length;
    return isStandby && hasConnectionViolations ? CSP_ERROR : '';
  }

  get formStates() {
    const { directLinkData, loginSettings } = this.args;

    if (directLinkData) {
      return this.directLinkViews;
    }

    if (loginSettings) {
      return this.loginSettingsViews;
    }

    if (this.visibleMountsByType) {
      return this.visibleMountViews;
    }

    // If none of the above, the UI renders the standard dropdown
    return null;
  }

  get preselectedType() {
    return this.canceledMfaAuth || this.args.directLinkData?.type || '';
  }

  get visibleMountsByType() {
    const visibleAuthMounts = this.args.visibleAuthMounts;
    if (visibleAuthMounts) {
      const authMounts = visibleAuthMounts;
      return Object.entries(authMounts).reduce((obj, [path, mountData]) => {
        const { type } = mountData;
        obj[type] ??= []; // if an array doesn't already exist for that type, create it
        obj[type].push({ path, ...mountData });
        return obj;
      }, {} as UnauthMountsByType);
    }
    return null;
  }

  get visibleMountTypes(): string[] {
    return Object.keys(this.visibleMountsByType || {});
  }

  // Form State Getters
  get directLinkViews() {
    const { directLinkData } = this.args;

    // If "path" key exists we know the "with" query param references a mount with listing_visibility="unauth"
    // Treat it as a preferred method and hide all other tabs.
    if (directLinkData?.path) {
      const tabData = this.filterVisibleMountsByType([directLinkData.type]);
      const defaultView = this.constructViews(FormView.TABS, tabData);
      const alternateView = this.constructViews(FormView.DROPDOWN, null);

      return { defaultView, alternateView };
    }

    // Otherwise, directLinkData just has a "type" key.
    // Render either the dropdown or visibleMountViews, with that type preselected
    return this.visibleMountsByType ? this.visibleMountViews : this.constructViews(FormView.DROPDOWN, null);
  }

  get loginSettingsViews() {
    const { loginSettings } = this.args;
    const defaultType = loginSettings?.defaultType;
    const backupTypes = loginSettings?.backupTypes;

    // If a default is not set, render backup methods as the initial view
    const preferredTypes = defaultType ? [defaultType] : backupTypes;
    let defaultView;
    if (preferredTypes) {
      const tabData = this.filterVisibleMountsByType(preferredTypes);
      defaultView = this.constructViews(FormView.TABS, tabData);
    }

    // Both default and backups must be set for an alternate view to exist
    let alternateView = null;
    if (defaultType && backupTypes) {
      const tabData = this.filterVisibleMountsByType(backupTypes);
      alternateView = this.constructViews(FormView.TABS, tabData);
    }

    return { defaultView, alternateView };
  }

  get visibleMountViews() {
    const defaultView = this.constructViews(FormView.TABS, this.visibleMountsByType);
    const alternateView = this.constructViews(FormView.DROPDOWN, null);
    return { defaultView, alternateView };
  }

  @action
  onAuthResponse(authResponse: AuthResponse | AuthResponseWithMfa, { selectedAuth = '', path = '' }) {
    const mfa_requirement = 'mfa_requirement' in authResponse ? authResponse.mfa_requirement : undefined;
    /*
    Checking for an mfa_requirement happens in two places.
    If doSubmit in <AuthForm> is called directly (by the <form> component) mfa is just handled here.
  
    Login methods submitted using a child form component of <AuthForm> are first checked for mfa 
    in the Auth::LoginForm "authenticate" task, and then that data eventually bubbles up here.
    */
    if (mfa_requirement) {
      // if an mfa requirement exists further action is required
      this.mfaAuthData = { mfa_requirement, selectedAuth, path };
    } else {
      // calls authSuccess in auth.js controller
      this.args.onAuthSuccess(authResponse);
    }
  }

  @action
  onCancelMfa() {
    // before resetting mfaAuthData, preserve auth type
    this.canceledMfaAuth = this.mfaAuthData?.selectedAuth ?? '';
    this.mfaAuthData = null;
  }

  @action
  onMfaSuccess(authResponse: AuthResponse) {
    // calls authSuccess in auth.js controller
    this.args.onAuthSuccess(authResponse);
  }

  @action
  onMfaErrorDismiss() {
    this.mfaAuthData = null;
    this.mfaErrors = '';
  }

  // Helpers
  private filterVisibleMountsByType(authTypes: string[]) {
    const tabs: UnauthMountsByType = {};
    for (const type of authTypes) {
      // adds visible mounts for each type, if they exist
      tabs[type] = this.visibleMountsByType?.[type] || null;
    }
    return tabs;
  }

  private constructViews(view: FormView, tabData: UnauthMountsByType | null) {
    return { view, tabData };
  }
}
