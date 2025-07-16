/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, fillable, visitable, selectable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/list'),
  visitShow: visitable('/vault/secrets/:backend/show/:id'),
  visitCreate: visitable('/vault/secrets/:backend/create'),
  dbPlugin: selectable('[data-test-input="plugin_name"]'),
  name: fillable('[data-test-input="name"]'),
  toggleVerify: clickable('[data-test-input="verify_connection"]'),
  connectionUrl: fillable('[data-test-input="connection_url"]'),
  url: fillable('[data-test-input="url"]'),
  username: fillable('[data-test-input="username"]'),
  password: fillable('[data-test-input="password"]'),
  save: clickable('[data-test-secret-save]'),
  addRole: clickable('[data-test-add-role]'),
  enable: clickable('[data-test-enable-connection]'),
  edit: clickable('[data-test-edit-link]'),
  delete: clickable('[data-test-database-connection-delete]'),
});
