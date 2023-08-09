/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * @module OidcConsentBlock
 * OidcConsentBlock components are used to show the consent form for the OIDC Authorization Code Flow
 *
 * @example
 * ```js
 * <OidcConsentBlock @redirect="https://example.com/oidc-callback" @code="abcd1234" @state="string-for-state" />
 * ```
 * @param {string} redirect - redirect is the URL where successful consent will redirect to
 * @param {string} code - code is the string required to pass back to redirect on successful OIDC auth
 * @param {string} [state] - state is a string which is required to return on redirect if provided, but optional generally
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

const validParameters = ['code', 'state'];
export default class OidcConsentBlockComponent extends Component {
  @tracked didCancel = false;

  get win() {
    return this.window || window;
  }

  buildUrl(urlString, params) {
    try {
      const url = new URL(urlString);
      Object.keys(params).forEach((key) => {
        if (params[key] && validParameters.includes(key)) {
          url.searchParams.append(key, params[key]);
        }
      });
      return url;
    } catch (e) {
      console.debug('DEBUG: parsing url failed for', urlString); // eslint-disable-line
      throw new Error('Invalid URL');
    }
  }

  @action
  handleSubmit(evt) {
    evt.preventDefault();
    const { redirect, ...params } = this.args;
    const redirectUrl = this.buildUrl(redirect, params);
    if (Ember.testing) {
      this.args.testRedirect(redirectUrl.toString());
    } else {
      this.win.location.replace(redirectUrl);
    }
  }

  @action
  handleCancel(evt) {
    evt.preventDefault();
    this.didCancel = true;
  }
}
