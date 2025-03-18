/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import authPage from 'vault/tests/pages/auth';
import { click, fillIn, currentURL, visit, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { v4 as uuidv4 } from 'uuid';

const SELECTORS = {
  menuItem: (item) => `[data-test-popup-menu="${item}"]`,
};

module('Acceptance | totp key backend', function (hooks) {
  setupApplicationTest(hooks);

  const createVaultKey = async (keyName, issuer, accountName) => {
    await fillIn(GENERAL.inputByAttr('name'), keyName);
    await fillIn(GENERAL.inputByAttr('issuer'), issuer);
    await fillIn(GENERAL.inputByAttr('accountName'), accountName);
    await click('[data-test-totp-create]');
    await click(GENERAL.backButton);
  };

  const createNonVaultKey = async (keyName, issuer, accountName, url) => {
    await click('[data-test-radio="Other service"]');
    await fillIn(GENERAL.inputByAttr('name'), keyName);
    await fillIn(GENERAL.inputByAttr('issuer'), issuer);
    await fillIn(GENERAL.inputByAttr('accountName'), accountName);
    await fillIn(GENERAL.inputByAttr('url'), url);
    await click('[data-test-totp-create]');
    await click(GENERAL.backButton);
  };

  hooks.beforeEach(async function () {
    this.uid = uuidv4();
    this.mountPath = `totp-${this.uid}`;
    this.path = `totp-${this.uid}`;
    this.keyName = 'totp-key';
    this.issuer = 'totp-issuer';
    this.accountName = 'totp-acount';
    this.url =
      'otpauth://totp/test-issuer:my-account?algorithm=SHA1&digits=6&issuer=test-issuer&period=30&secret=HICPOBIMFO4YYHFYX3QPVYUL2YEPVJKU';

    await authPage.login();
    // Setup TOTP engine
    await visit('/vault/settings/mount-secret-backend');
    await mountBackend('totp', this.mountPath);
  });

  test('it views a key via menu option', async function (assert) {
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/list`,
      'After enabling totp secrets engine it navigates to keys list'
    );

    await click(SES.createSecret);
    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await visit(`/vault/secrets/${this.path}`);
    await click(GENERAL.menuTrigger);
    await click(`${SELECTORS.menuItem('details')}`);

    assert.dom('.title').hasText(`TOTP key ${this.keyName}`);
    assert.dom('[data-test-totp-key-details]').exists();

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'After clicking details option it navigates to key detail view'
    );
  });

  test('it deletes a key via menu option', async function (assert) {
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/list`,
      'After enabling totp secrets engine it navigates to keys list'
    );

    await click(SES.createSecret);
    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    await visit(`/vault/secrets/${this.path}`);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.dom(SES.secretLink(this.keyName)).doesNotExist(`${this.keyName}: key is no longer in the list`);
  });

  test('it creates a key with Vault as the provider', async function (assert) {
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/list`,
      'After enabling totp secrets engine it navigates to keys list'
    );

    await click(SES.createSecret);
    assert.dom(SES.secretHeader).hasText('Create a TOTP key', 'It renders the create key page');

    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'totp: navigates to the show page on creation'
    );
  });

  test('it creates a key with another service as the provider', async function (assert) {
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/list`,
      'After enabling totp secrets engine it navigates to keys list'
    );

    await click(SES.createSecret);
    assert.dom(SES.secretHeader).hasText('Create a TOTP key', 'It renders the create key page');
    await createNonVaultKey(this.keyName, this.issuer, this.accountName, this.url);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'totp: navigates to the show page on creation'
    );
  });
});
