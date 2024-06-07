/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { collection, clickable, fillable, text, value, isPresent } from 'ember-cli-page-object';

export default {
  tabs: collection('[data-test-auth-method]', {
    name: text(),
    link: clickable('[data-test-auth-method-link]'),
  }),
  selectMethod: fillable('[data-test-select=auth-method]'),
  username: fillable('[data-test-username]'),
  token: fillable('[data-test-token]'),
  tokenValue: value('[data-test-token]'),
  password: fillable('[data-test-password]'),
  errorText: text('[data-test-message-error]'),
  errorMessagePresent: isPresent('[data-test-message-error]'),
  descriptionText: text('[data-test-description]'),
  login: clickable('[data-test-auth-submit]'),
  oidcRole: fillable('[data-test-role]'),
  oidcMoreOptions: clickable('[data-test-yield-content] button'),
  oidcMountPath: fillable('#custom-path'),
};
