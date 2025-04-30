/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';

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
 * @wrappedToken={{this.wrappedToken}}
/>
 *
 * @param {object} cluster - the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {string} namespaceQueryParam - namespace to login with, updated by typing in to the namespace input
 * @param {string} oidcProviderQueryParam - oidc provider query param, set in url as "?o=someprovider"
 * @param {function} onAuthSuccess - callback task in controller that receives the auth response (after MFA, if enabled) when login is successful
 * @param {function} onNamespaceUpdate - callback task that passes user input to the controller to update the login namespace in the url query params
 * @param {string} wrappedToken - passed down to the AuthForm component and can be used to login if added directly to the URL via the "wrapped_token" query param
 * */

export default class AuthPage extends Component {
  @service api;
  @service auth;
  @service flags;

  @tracked preselectedAuthType;

  // application state error handling
  @tracked mfaErrors = '';
  @tracked mfaAuthData;
  @tracked tokenUnwrapError = '';

  constructor() {
    super(...arguments);
    if (this.args.wrappedToken) {
      this.unwrapToken.perform();
    }
  }

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

  get pageError() {
    if (this.mfaErrors) {
      return {
        message:
          'Multi-factor authentication is required, but failed. Go back and try again, or contact your administrator.',
        error: this.mfaErrors,
        action: 'retryMfa',
      };
    }
    if (this.tokenUnwrapError) {
      return {
        message: 'Token unwrap failed',
        error: this.tokenUnwrapError,
        action: 'retryTokenUnwrap',
      };
    }
    return null;
  }

  unwrapToken = task(async () => {
    try {
      const response = await this.api.sys.unwrap(
        {},
        this.api.buildHeaders({ token: this.args.wrappedToken })
      );
      const authResponse = await this.auth.authenticate({
        clusterId: this.args.cluster.id,
        backend: 'token',
        data: { token: response.auth.clientToken },
        selectedAuth: 'token',
      });

      this.onAuthResponse(authResponse);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.tokenUnwrapError = message;
    }
  });

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
  dismissError(action) {
    if (action === 'retryTokenUnwrap') {
      this.tokenUnwrapError = '';
    }
    if (action === 'retryMfa') {
      this.mfaAuthData = null; // resets auth form
      this.mfaErrors = '';
    }
  }
}
