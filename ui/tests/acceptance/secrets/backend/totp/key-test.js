/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, visit, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { v4 as uuidv4 } from 'uuid';
import sinon from 'sinon';

module('Acceptance | totp key backend', function (hooks) {
  setupApplicationTest(hooks);

  const createVaultKey = async (keyName, issuer, accountName, exported = true, qrSize = 200) => {
    await fillIn(GENERAL.inputByAttr('name'), keyName);
    await fillIn(GENERAL.inputByAttr('issuer'), issuer);
    await fillIn(GENERAL.inputByAttr('accountName'), accountName);
    if (!exported) {
      await click(GENERAL.toggleInput('toggle-exported'));
    }
    if (qrSize !== 200) {
      await click(GENERAL.button('Provider Options'));
      await fillIn(GENERAL.inputByAttr('qrSize'), qrSize);
    }
    await click(GENERAL.submitButton);
  };

  const createNonVaultKey = async (keyName, issuer, accountName, url, key) => {
    await click(GENERAL.radioByAttr('Other service'));
    await fillIn(GENERAL.inputByAttr('name'), keyName);
    await fillIn(GENERAL.inputByAttr('issuer'), issuer);
    await fillIn(GENERAL.inputByAttr('accountName'), accountName);
    if (url) await fillIn(GENERAL.inputByAttr('url'), url);
    if (key) await fillIn(GENERAL.inputByAttr('key'), key);
    await click(GENERAL.submitButton);
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
    this.key = 'VCUDXBWFQXUEAIJYXH5YB62D5WFNQXFA';

    const flash = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = sinon.spy(flash, 'success');

    await login();
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

    await click(SES.createSecretLink);
    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await click(GENERAL.backButton);
    await visit(`/vault/secrets/${this.path}`);
    await click(GENERAL.menuTrigger);
    await click(`${GENERAL.menuItem('details')}`);

    assert.dom('.title').hasText(`TOTP key ${this.keyName}`);
    assert.dom('[data-test-totp-key-details]').exists();

    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'After clicking details option it navigates to key detail view'
    );
  });

  test('it deletes a key via menu option', async function (assert) {
    await click(SES.createSecretLink);
    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await click(GENERAL.backButton);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    await visit(`/vault/secrets/${this.path}`);
    await click(GENERAL.menuTrigger);
    await click(GENERAL.confirmTrigger);
    await click(GENERAL.confirmButton);
    assert.dom(SES.secretLink(this.keyName)).doesNotExist(`${this.keyName}: key is no longer in the list`);

    const [flashMessage] = this.flashSuccessSpy.lastCall.args;
    assert.strictEqual(flashMessage, `${this.keyName} was successfully deleted.`);
  });

  test('it creates a key with Vault as the provider', async function (assert) {
    await click(SES.createSecretLink);
    assert.dom(SES.secretHeader).hasText('Create a TOTP key', 'It renders the create key page');

    await createVaultKey(this.keyName, this.issuer, this.accountName);
    assert.dom('[data-test-qrcode]').exists('QR code exists');
    assert.dom(GENERAL.infoRowLabel('URL')).exists('URL exists');

    await click(GENERAL.backButton);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'totp: navigates to the show page on creation'
    );

    const [flashMessage] = this.flashSuccessSpy.lastCall.args;
    assert.strictEqual(flashMessage, 'Successfully created key.');
  });

  test('it creates a key with another service as the provider with URL', async function (assert) {
    await click(SES.createSecretLink);
    assert.dom(SES.secretHeader).hasText('Create a TOTP key', 'It renders the create key page');
    await createNonVaultKey(this.keyName, this.issuer, this.accountName, this.url);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'totp: navigates to the show page on creation'
    );

    const [flashMessage] = this.flashSuccessSpy.lastCall.args;
    assert.strictEqual(flashMessage, 'Successfully created key.');
  });

  test('it creates a key with another service as the provider with key', async function (assert) {
    await click(SES.createSecretLink);
    assert.dom(SES.secretHeader).hasText('Create a TOTP key', 'It renders the create key page');
    await createNonVaultKey(this.keyName, this.issuer, this.accountName, undefined, this.key);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'totp: navigates to the show page on creation'
    );

    const [flashMessage] = this.flashSuccessSpy.lastCall.args;
    assert.strictEqual(flashMessage, 'Successfully created key.');
  });

  test('it does not render QR code or URL when exported is false', async function (assert) {
    await click(SES.createSecretLink);
    await createVaultKey(this.keyName, this.issuer, this.accountName, false);
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`);
    assert.dom('[data-test-qrcode]').doesNotExist('QR code is not displayed');
    assert.dom(GENERAL.infoRowLabel('URL')).doesNotExist('URl is not displayed');
  });

  test('it does not render QR code when QR size is 0', async function (assert) {
    await click(SES.createSecretLink);
    await createVaultKey(this.keyName, this.issuer, this.accountName, true, 0);
    assert.dom('[data-test-qrcode]').doesNotExist('QR code is not displayed');
    assert.dom(GENERAL.infoRowLabel('URL')).exists('URl is displayed');
  });
});
