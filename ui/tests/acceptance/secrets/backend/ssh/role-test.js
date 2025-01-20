/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import editPage from 'vault/tests/pages/secrets/backend/ssh/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/ssh/show';
import generatePage from 'vault/tests/pages/secrets/backend/ssh/generate-otp';
import listPage from 'vault/tests/pages/secrets/backend/list';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';

module('Acceptance | secrets/ssh', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  const mountAndNav = async (uid) => {
    const path = `ssh-${uid}`;
    await enablePage.enable('ssh', path);
    await settled();
    await editPage.visitRoot({ backend: path });
    await settled();
    return path;
  };

  test('it creates a role and redirects', async function (assert) {
    assert.expect(5);
    const path = await mountAndNav(this.uid);
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
    assert.dom('.hds-dropdown li').exists({ count: 5 }, 'Renders 5 popup menu items');
  });

  test('it deletes a role', async function (assert) {
    assert.expect(2);
    const path = await mountAndNav(this.uid);
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
    const path = await mountAndNav(this.uid);
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
