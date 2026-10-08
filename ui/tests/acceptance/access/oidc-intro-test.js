/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

const ROUTE = '/vault/access/oidc';

module('Acceptance | oidc provider | intro page', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    return login();
  });

  test('Mode 1 (Full-Screen Card): shows intro card in initial state, dismisses on Skip, and reveals header CTA', async function (assert) {
    await visit(ROUTE);

    // 1. Full-screen intro card is shown; regular page header and outlet are hidden
    assert.dom(GENERAL.wizardIntro).exists('full-screen wizard intro card is shown on initial visit');
    assert.dom(GENERAL.button('Skip')).exists('"Skip" button is present in full-screen intro card');
    assert
      .dom(GENERAL.button('intro'))
      .doesNotExist('"New to OIDC Provider?" button is not visible during full-screen intro');

    // 2. Dismiss via Skip
    await click(GENERAL.button('Skip'));

    // 3. Wizard is hidden and the page header now shows the "New to OIDC Provider?" button
    assert.dom(GENERAL.wizardIntro).doesNotExist('wizard intro is hidden after dismissal');
    assert
      .dom(GENERAL.button('intro'))
      .exists('"New to OIDC Provider?" button is present in page header after skip');
    assert
      .dom(GENERAL.emptyStateTitle)
      .hasText('No OIDC Provider configurations', 'empty state is shown after skip');
  });

  test('Mode 2 (Modal View): re-opens intro as modal via header CTA and closes via Close button', async function (assert) {
    // Start with wizard already dismissed
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.oidcProvider);
    await visit(ROUTE);

    assert.dom(GENERAL.wizardIntro).doesNotExist('intro is initially dismissed');
    assert.dom(GENERAL.button('intro')).exists('"New to OIDC Provider?" button is visible');
    assert.dom('[data-test-oidc-configure]').exists('"Create your first app" button is visible');

    // 1. Click header CTA to re-open intro as modal
    await click(GENERAL.button('intro'));

    // 2. Verify modal presentation
    assert.dom(GENERAL.modal.container('OIDC Provider')).exists('modal dialog is rendered');
    assert.dom(GENERAL.wizardIntro).exists('wizard intro is visible inside modal');

    // 3. Close the modal
    await click(GENERAL.button('Close'));
    assert.dom(GENERAL.modal.container('OIDC Provider')).doesNotExist('modal dialog is closed');
    assert.dom(GENERAL.wizardIntro).doesNotExist('intro is dismissed after closing modal');
  });

  test('does not show intro card or header CTA when clients exist', async function (assert) {
    // Seed a client so the index route redirects to the clients list view
    this.server.get('/identity/oidc/client', () => ({
      data: { keys: ['my-app'], key_info: { 'my-app': {} } },
    }));

    await visit(ROUTE);

    assert.dom(GENERAL.wizardIntro).doesNotExist('wizard intro is not shown when clients exist');
    assert
      .dom(GENERAL.button('intro'))
      .doesNotExist('"New to OIDC Provider?" button is not shown when clients exist');
    assert
      .dom('[data-test-oidc-configure]')
      .doesNotExist('"Create your first app" button is not shown when clients exist');
  });
});
