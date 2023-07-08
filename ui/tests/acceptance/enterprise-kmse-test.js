/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, fillIn } from '@ember/test-helpers';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Acceptance | Enterprise | keymgmt', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await logout.visit();
    return authPage.login();
  });

  test('it should add new key and distribute to provider', async function (assert) {
    const path = `keymgmt-${Date.now()}`;
    this.server.post(`/${path}/key/test-key`, () => ({}));
    this.server.put(`/${path}/kms/test-keyvault/key/test-key`, () => ({}));

    await mountSecrets.enable('keymgmt', path);
    await click('[data-test-secret-create]');
    await fillIn('[data-test-input="provider"]', 'azurekeyvault');
    await fillIn('[data-test-input="name"]', 'test-keyvault');
    await fillIn('[data-test-input="keyCollection"]', 'test-keycollection');
    await fillIn('[data-test-input="credentials.client_id"]', '123');
    await fillIn('[data-test-input="credentials.client_secret"]', '456');
    await fillIn('[data-test-input="credentials.tenant_id"]', '789');
    await click('[data-test-kms-provider-submit]');
    await click('[data-test-distribute-key]');
    await click('[data-test-component="search-select"] .ember-basic-dropdown-trigger');
    await fillIn('.ember-power-select-search-input', 'test-key');
    await click('.ember-power-select-option');
    await fillIn('[data-test-keymgmt-dist-keytype]', 'rsa-2048');
    await click('[data-test-operation="encrypt"]');
    await fillIn('[data-test-protection="hsm"]', 'hsm');

    this.server.get(`/${path}/kms/test-keyvault/key`, () => ({ data: { keys: ['test-key'] } }));
    await click('[data-test-secret-save]');
    await click('[data-test-kms-provider-tab="keys"] a');
    assert.dom('[data-test-secret-link="test-key"]').exists('Key is listed under keys tab of provider');
  });
});
