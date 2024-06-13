/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Base from 'ember-simple-auth/authenticators/base';

/**
 * This authenticator is for using ember-simple-auth within the current Auth flow.
 * It will eventually be replaced by a more robust authenticator which handles the
 * authentication business logic which currently lives in services/auth.js
 */
export default class BasicAuthenticator extends Base {
  /**
   * restores the auth state when the page is refreshed or entered for the first time
   * @param {object} data restored from store (eg localStorage)
   */
  restore(data) {
    return data?.token ? Promise.resolve(data) : Promise.reject();
  }

  authenticate(payload) {
    return Promise.resolve(payload);
  }

  invalidate() {
    return Promise.resolve();
  }
}
