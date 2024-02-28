/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import { click, settled } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';

const { rootToken } = VAULT_KEYS;

export default create({
  visit: visitable('/vault/auth'),
  logout: visitable('/vault/logout'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  usernameInput: fillable('[data-test-username]'),
  passwordInput: fillable('[data-test-password]'),
  namespaceInput: fillable('[data-test-auth-form-ns-input]'),
  optionsToggle: clickable('[data-test-auth-form-options-toggle]'),
  mountPath: fillable('[data-test-auth-form-mount-path]'),
  authType: fillable('[data-test-select="auth-method"]'),

  login: async function (token = rootToken) {
    // make sure we're always logged out and logged back in
    await this.logout();
    // clear session storage to ensure we have a clean state
    window.localStorage.clear();
    await this.visit({ with: 'token' });
    await settled();
    await this.tokenInput(token).submit();
    await settled();
    return;
  },
  loginUsername: async function (username, password, path) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await settled();
    // clear local storage to ensure we have a clean state
    await window.localStorage.clear();
    await this.visit({ with: 'userpass' });
    await settled();
    if (path) {
      await this.optionsToggle();
      await this.mountPath(path);
    }
    await this.usernameInput(username);
    await this.passwordInput(password).submit();
    return;
  },
  loginNs: async function (ns, token = rootToken) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await settled();
    // clear session storage to ensure we have a clean state
    await window.localStorage.clear();
    await this.visit({ with: 'token' });
    await settled();
    await this.namespaceInput(ns);
    await settled();
    await this.tokenInput(token).submit();
    return;
  },
  clickLogout: async function (clearNamespace = false) {
    await click('[data-test-user-menu-trigger]');
    await click('[data-test-user-menu-content] a#logout');
    if (clearNamespace) {
      await this.namespaceInput('');
    }
    await settled();
    return;
  },
});
