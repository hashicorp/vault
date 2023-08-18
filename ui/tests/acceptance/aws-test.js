/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, findAll, currentURL, find, settled, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | aws secret backend', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('aws backend', async function (assert) {
    assert.expect(12);
    const path = `aws-${this.uid}`;
    const roleName = 'awsrole';

    await enablePage.enable('aws', path);
    await settled();
    await click('[data-test-configuration-tab]');

    await click('[data-test-secret-backend-configure]');

    assert.strictEqual(currentURL(), `/vault/settings/secrets/configure/${path}`);
    assert.ok(findAll('[data-test-aws-root-creds-form]').length, 'renders the empty root creds form');
    assert.ok(findAll('[data-test-aws-link="root-creds"]').length, 'renders the root creds link');
    assert.ok(findAll('[data-test-aws-link="leases"]').length, 'renders the leases config link');

    await fillIn('[data-test-aws-input="accessKey"]', 'foo');
    await fillIn('[data-test-aws-input="secretKey"]', 'bar');

    await click('[data-test-aws-input="root-save"]');

    assert.ok(
      find('[data-test-flash-message]').textContent.trim(),
      `The backend configuration saved successfully!`
    );

    await click('[data-test-aws-link="leases"]');

    await click('[data-test-aws-input="lease-save"]');

    assert.ok(
      find('[data-test-flash-message]').textContent.trim(),
      `The backend configuration saved successfully!`
    );

    await click('[data-test-backend-view-link]');

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`, `navigates to the roles list`);

    await click('[data-test-secret-create]');

    assert.ok(
      find('[data-test-secret-header]').textContent.includes('AWS Role'),
      `aws: renders the create page`
    );

    await fillIn('[data-test-input="name"]', roleName);

    // save the role
    await click('[data-test-role-aws-create]');
    await waitUntil(() => currentURL() === `/vault/secrets/${path}/show/${roleName}`); // flaky without this
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      `$aws: navigates to the show page on creation`
    );

    await click('[data-test-secret-root-link]');

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`);
    assert.ok(findAll(`[data-test-secret-link="${roleName}"]`).length, `aws: role shows in the list`);

    //and delete
    await click(`[data-test-secret-link="${roleName}"] [data-test-popup-menu-trigger]`);
    await waitUntil(() => find(`[data-test-aws-role-delete="${roleName}"]`)); // flaky without
    await click(`[data-test-aws-role-delete="${roleName}"]`);
    await click(`[data-test-confirm-button]`);
    assert.dom(`[data-test-secret-link="${roleName}"]`).doesNotExist(`aws: role is no longer in the list`);
  });
});
