/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import editPage from 'vault/tests/pages/secrets/backend/ssh/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/ssh/show';
import generatePage from 'vault/tests/pages/secrets/backend/ssh/generate-otp';
import listPage from 'vault/tests/pages/secrets/backend/list';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/ssh', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  const mountAndNav = async () => {
    const path = `ssh-${new Date().getTime()}`;
    await enablePage.enable('ssh', path);
    await settled();
    await editPage.visitRoot({ backend: path });
    await settled();
    return path;
  };

  test('it creates a role and redirects', async function (assert) {
    assert.expect(5);
    const path = await mountAndNav(assert);
    await editPage.createOTPRole('role');
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.generateIsPresent, 'shows the generate button');

    await showPage.visit({ backend: path, id: 'role' });
    await settled();
    await showPage.generate();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await listPage.visitRoot({ backend: path });
    await settled();
    assert.strictEqual(listPage.secrets.length, 1, 'shows role in the list');
    const secret = listPage.secrets.objectAt(0);
    await secret.menuToggle();
    assert.ok(listPage.menuItems.length > 0, 'shows links in the menu');
  });

  test('it deletes a role', async function (assert) {
    assert.expect(2);
    const path = await mountAndNav(assert);
    await editPage.createOTPRole('role');
    await settled();
    await showPage.visit({ backend: path, id: 'role' });
    await settled();
    await showPage.deleteRole();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.list-root',
      'redirects to list page'
    );
    assert.ok(listPage.backendIsEmpty, 'no roles listed');
  });

  test('it generates an OTP', async function (assert) {
    assert.expect(6);
    const path = await mountAndNav(assert);
    await editPage.createOTPRole('role');
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.show',
      'redirects to the show page'
    );
    assert.ok(showPage.generateIsPresent, 'shows the generate button');

    await showPage.visit({ backend: path, id: 'role' });
    await settled();
    await showPage.generate();
    await settled();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.credentials',
      'navs to the credentials page'
    );

    await generatePage.generateOTP();
    await settled();
    assert.ok(generatePage.warningIsPresent, 'shows warning');
    await generatePage.back();
    await settled();
    assert.ok(generatePage.userIsPresent, 'clears generate, shows user input');
    assert.ok(generatePage.ipIsPresent, 'clears generate, shows ip input');
  });
});
