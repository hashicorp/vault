/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module AuthRoutePage
 * The Auth::RoutePage wraps OktaNumberChallenge and AuthForm to manage the login flow and is responsible for calling the authenticate method
 *
 * @example
 * <Auth::RoutePag @namespaceQueryParam={{this.namespaceQueryParam} @onAuthSuccess={{action "authSuccess"}} @oidcProviderQueryParam={{this.oidcProvider} @cluster={{this.model} @onNamespaceUpdate={{perform this.updateNamespace}} />
 *
 * @param {string} param - info about the param
 * */

export default class AuthRoutePage extends Component {
  @service auth;
  @service flags;
  @service namespaceService;

  @tracked mfaErrors;
  @tracked mfaAuthData;

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
    // if an mfa requirement exists further action is required
    if (mfa_requirement) {
      this.mfaAuthData = { mfa_requirement, backend, data };
    } else {
      this.args.onAuthSuccess(authResponse);
    }
  }

  @action
  onMfaSuccess(authResponse) {
    this.arg.onAuthSuccess(authResponse);
  }

  @action
  onMfaErrorDismiss() {
    this.mfaAuthData = null;
    this.mfaErrors = null;
  }
}
