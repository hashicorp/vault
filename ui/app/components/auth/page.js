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
 * The Auth::Page is the route template for the login splash view. It renders the Auth::LoginForm or MFA component if an 
 * mfa validation is returned from the auth request. It also handles display logic if there is an oidc provider query param.
 *
 * @example
 * <Auth::Page
 * @cluster={{this.model}}
 * @namespaceQueryParam={{this.namespaceQueryParam}}
 * @oidcProviderQueryParam={{this.oidcProvider}}
 * @onAuthSuccess={{action "authSuccess"}}
 * @onNamespaceUpdate={{perform this.updateNamespace}}
/>
 *
 * @param {object} cluster - the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {string} namespaceQueryParam - namespace to login with, updated by typing in to the namespace input
 * @param {string} oidcProviderQueryParam - oidc provider query param, set in url as "?o=someprovider"
 * @param {function} onAuthSuccess - callback task in controller that receives the auth response (after MFA, if enabled) when login is successful
 * @param {function} onNamespaceUpdate - callback task that passes user input to the controller to update the login namespace in the url query params
 * */

export default class AuthPage extends Component {
  @service flags;

  @tracked mfaAuthData;
  @tracked mfaErrors = '';
  @tracked preselectedAuthType;

  get namespaceInput() {
    const namespaceQP = this.args.namespaceQueryParam;
    if (this.flags.hvdManagedNamespaceRoot) {
      // When managed, the user isn't allowed to edit the prefix `admin/` for their nested namespace
      const split = namespaceQP.split('/');
      if (split.length > 1) {
        split.shift();
        return `/${split.join('/')}`;
      }
      return '';
    }
    return namespaceQP;
  }

  @action
  handleNamespaceUpdate(event) {
    this.args.onNamespaceUpdate(event.target.value);
  }

  @action
  onAuthResponse(authResponse, backend, data) {
    const { mfa_requirement } = authResponse;
    /*
    Checking for an mfa_requirement happens in two places.
    If doSubmit in <AuthForm> is called directly (by the <form> component) mfa is just handled here.
  
    Login methods submitted using a child form component of <AuthForm> are first checked for mfa 
    in the Auth::LoginForm "authenticate" task, and then that data eventually bubbles up here.
    */
    if (mfa_requirement) {
      // if an mfa requirement exists further action is required
      this.mfaAuthData = { mfa_requirement, backend, data };
    } else {
      // calls authSuccess in auth.js controller
      this.args.onAuthSuccess(authResponse);
    }
  }

  @action
  onCancelMfa() {
    // before resetting mfaAuthData, preserve auth type
    this.preselectedAuthType = this.mfaAuthData.backend;
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
