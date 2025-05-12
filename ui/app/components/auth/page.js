/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module AuthPage
 * The Auth::Page is the route template for the login splash view. It renders the Auth::FormTemplate or MFA component if an
 * mfa validation is returned from the auth request. It also formats mount data to manage what tabs are rendered in Auth::FormTemplate.
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
 * @param {string} directLinkData - type or mount data gleaned from query param
 * @param {object} cluster - the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {string} namespaceQueryParam - namespace to login with, updated by typing in to the namespace input
 * @param {string} oidcProviderQueryParam - oidc provider query param, set in url as "?o=someprovider"
 * @param {function} onAuthSuccess - callback task in controller that receives the auth response (after MFA, if enabled) when login is successful
 * @param {function} onNamespaceUpdate - callback task that passes user input to the controller to update the login namespace in the url query params
 * @param {object} visibleAuthMounts - mount paths with listing_visibility="unauth", keys are the mount path and value is it's mount data such as "type" or "description," if it exists
 * */

export const CSP_ERROR =
  "This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.";

export default class AuthPage extends Component {
  @service('csp-event') csp;

  @tracked canceledMfaAuth = '';
  @tracked mfaAuthData;
  @tracked mfaErrors = '';

  get visibleMountsByType() {
    const visibleAuthMounts = this.args.visibleAuthMounts;
    if (visibleAuthMounts) {
      const authMounts = visibleAuthMounts;
      return Object.entries(authMounts).reduce((obj, [path, mountData]) => {
        const { type } = mountData;
        obj[type] ??= []; // if an array doesn't already exist for that type, create it
        obj[type].push({ path, ...mountData });
        return obj;
      }, {});
    }
    return null;
  }

  get cspError() {
    const isStandby = this.args.cluster.standby;
    const hasConnectionViolations = this.csp.connectionViolations.length;
    return isStandby && hasConnectionViolations ? CSP_ERROR : '';
  }

  @action
  onAuthResponse(authResponse, { selectedAuth, path }) {
    const { mfa_requirement } = authResponse;
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
    this.canceledMfaAuth = this.mfaAuthData.selectedAuth;
    this.mfaAuthData = null;
  }

  @action
  onMfaSuccess(authResponse) {
    // calls authSuccess in auth.js controller
    this.args.onAuthSuccess(authResponse);
  }

  @action
  onMfaErrorDismiss() {
    this.mfaAuthData = null;
    this.mfaErrors = '';
  }
}
