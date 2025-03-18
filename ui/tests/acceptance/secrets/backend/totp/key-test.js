/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import listPage from 'vault/tests/pages/secrets/backend/list';
import { click, fillIn, currentURL, settled, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Acceptance | totp key backend', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  const createVaultKey = async (keyName, issuer, accountName) => {
    await fillIn(GENERAL.inputByAttr('name'), keyName);
    await fillIn(GENERAL.inputByAttr('issuer'), issuer);
    await fillIn(GENERAL.inputByAttr('accountName'), accountName);
    await click('[data-test-totp-create]');
  };

  hooks.beforeEach(async function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');

    this.uid = uuidv4();
    this.mountPath = `totp-${this.uid}`;
    this.path = `totp-${this.uid}`;
    this.keyName = 'totp-key';
    this.issuer = 'totp-issuer';
    this.accountName = 'totp-acount';

    await authPage.login();
    // Setup TOTP engine
    await enablePage.enable('totp', this.mountPath);
  });

  test('it views a key via menu option', async function (assert) {
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/list`,
      'After enabling totp secrets engine it navigates to keys list'
    );

    await click(SES.createSecret);
    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await settled();
    await listPage.visit({ backend: this.path, id: this.keyName });
    await click('[data-test-popup-menu-trigger]');
    assert.dom('.hds-dropdown li:nth-of-type(1)').hasText('Details', 'first list item is "details"');
    await click('.hds-dropdown li:nth-of-type(1) a');
    assert.dom('.title').hasText(`TOTP key ${this.keyName}`);
    assert.dom('[data-test-totp-key-details]').exists();
  });

  test('it deletes a key via menu option', async function (assert) {
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/list`,
      'After enabling totp secrets engine it navigates to keys list'
    );

    await click(SES.createSecret);
    await createVaultKey(this.keyName, this.issuer, this.accountName);
    await settled();
    await listPage.visit({ backend: this.path, id: this.keyName });
    await click('[data-test-popup-menu-trigger]');
    assert.dom('.hds-dropdown li:nth-of-type(2)').hasText('Delete', 'first list item is "details"');
    await click('.hds-dropdown li:nth-of-type(2) [data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button]');
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
    await waitUntil(() => currentURL() === `/vault/secrets/${this.path}/show/${this.keyName}`); // flaky without this
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.path}/show/${this.keyName}`,
      'totp: navigates to the show page on creation'
    );
  });
});
