/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { setupMirage } from 'ember-cli-mirage/test-support';

import authPage from 'vault/tests/pages/auth';
import assertSecretWrap from 'vault/tests/helpers/secret-engines';

const SELECTORS = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  createSecret: '[data-test-secret-create]',
  editSecret: '[data-test-secret-edit]',
  saveBtn: '[data-test-secret-save]',
  keyInput: (idx = 0) => `[data-test-kv-key="${idx}"]`,
  maskedValueInput: (idx = 0) => `[data-test-kv-value="${idx}"] [data-test-textarea]`,
  label: '[data-test-label-div]',
};
module('Acceptance | secrets/cubbyhole/create', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it creates and can view a secret with the cubbyhole backend', async function (assert) {
    assert.expect(6);
    const kvPath = `cubbyhole-kv-${this.uid}`;
    const requestPath = `cubbyhole/${kvPath}`;
    await visit('/vault/secrets/cubbyhole/list');

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'navigates to the list page'
    );

    await click(SELECTORS.createSecret);
    await fillIn(SELECTORS.inputByAttr('path'), kvPath);
    await fillIn(SELECTORS.keyInput(), 'foo');
    await fillIn(SELECTORS.maskedValueInput(), 'bar');
    await click(SELECTORS.saveBtn);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.dom(SELECTORS.label).hasText('foo');
    assert.dom('[data-test-created-time]').hasText('', 'it does not render created time if blank');
    await assertSecretWrap(assert, this.server, requestPath);

    await click(SELECTORS.editSecret);
    assert.dom(SELECTORS.inputByAttr('path')).doesNotExist('does not render path on edit');
    assert.dom(SELECTORS.keyInput()).hasValue('foo');
  });
});
