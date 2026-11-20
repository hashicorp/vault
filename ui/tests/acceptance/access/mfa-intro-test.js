/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, currentRouteName, visit } from '@ember/test-helpers';
import { Response } from 'miragejs';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import mfaConfigHandler from 'vault/mirage/handlers/mfa-config';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

const ROUTE = '/vault/access/mfa';

module('Acceptance | mfa | intro page', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    mfaConfigHandler(this.server);
    return login();
  });

  hooks.afterEach(function () {
    this.owner.lookup('service:wizard').reset(WIZARD_ID_MAP.mfa);
  });

  module('with no MFA methods', function (hooks) {
    hooks.beforeEach(function () {
      const notFound = () => new Response(404, {}, { errors: [] });
      this.server.get('/identity/mfa/method/', notFound);
      this.server.get('/identity/mfa/method', notFound);
    });

    test('Mode 1 (Full-Screen Card): shows intro, dismisses on Skip, and reveals header CTA', async function (assert) {
      await visit(ROUTE);

      assert.dom(GENERAL.wizardIntro).exists('full-screen intro is shown on initial visit');
      assert.dom(GENERAL.button('Skip')).exists('"Skip" button is present');
      assert.dom(GENERAL.button('intro')).doesNotExist('"New to MFA?" is hidden during full-screen intro');
      assert.dom('[data-test-mfa-configure]').doesNotExist('page "Configure MFA" is hidden during intro');

      await click(GENERAL.button('Skip'));

      assert.dom(GENERAL.wizardIntro).doesNotExist('intro is hidden after dismissal');
      assert.dom(GENERAL.button('intro')).exists('"New to MFA?" is shown after skip');
      assert.dom('[data-test-mfa-configure]').exists('"Configure MFA" is shown after skip');
    });

    test('Mode 2 (Modal View): re-opens intro via header CTA and closes via Close', async function (assert) {
      this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.mfa);
      await visit(ROUTE);

      assert.dom(GENERAL.wizardIntro).doesNotExist('intro is initially dismissed');
      await click(GENERAL.button('intro'));

      assert.dom(GENERAL.modal.container('Require a second factor at login')).exists('modal is rendered');
      assert.dom(GENERAL.wizardIntro).exists('intro is visible inside modal');

      await click(GENERAL.button('Close'));
      assert.dom(GENERAL.modal.container('Require a second factor at login')).doesNotExist('modal closed');
      assert.dom(GENERAL.wizardIntro).doesNotExist('intro is dismissed after closing');
    });

    test('"Configure MFA" in the intro navigates to create method', async function (assert) {
      await visit(ROUTE);
      await click(GENERAL.button('Configure MFA'));
      assert.strictEqual(currentRouteName(), 'vault.cluster.access.mfa.methods.create');
    });
  });

  test('does not show intro or header CTA when methods exist', async function (assert) {
    await visit(ROUTE);

    assert.strictEqual(currentRouteName(), 'vault.cluster.access.mfa.methods.index');
    assert.dom(GENERAL.wizardIntro).doesNotExist('intro is not shown when methods exist');
    assert.dom(GENERAL.button('intro')).doesNotExist('"New to MFA?" is not shown when methods exist');
  });
});
