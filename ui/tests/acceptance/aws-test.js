/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, find, settled, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import { GENERAL } from '../helpers/general-selectors';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | aws secret backend', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');

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

    assert.dom('[data-test-aws-root-creds-form]').exists('renders the empty root creds form');
    assert.dom(GENERAL.tab('access-to-aws')).exists('renders the root creds tab');
    assert.dom(GENERAL.tab('lease')).exists('renders the leases config tab');

    await fillIn('[data-test-aws-input="accessKey"]', 'foo');
    await fillIn('[data-test-aws-input="secretKey"]', 'bar');

    await click('[data-test-aws-input="root-save"]');

    assert.true(
      this.flashSuccessSpy.calledWith('The backend configuration saved successfully!'),
      'success flash message is rendered'
    );

    await click(GENERAL.tab('lease'));

    await click('[data-test-aws-input="lease-save"]');

    assert.true(
      this.flashSuccessSpy.calledTwice,
      'a new success flash message is rendered upon saving lease'
    );

    await click('[data-test-backend-view-link]');

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`, 'navigates to the roles list');

    await click('[data-test-secret-create]');

    assert.dom('[data-test-secret-header]').hasText('Create an AWS Role', 'aws: renders the create page');

    await fillIn('[data-test-input="name"]', roleName);

    // save the role
    await click('[data-test-role-aws-create]');
    await waitUntil(() => currentURL() === `/vault/secrets/${path}/show/${roleName}`); // flaky without this
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      'aws: navigates to the show page on creation'
    );
    await click(`[data-test-secret-breadcrumb="${path}"] a`);

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`);
    assert.dom(`[data-test-secret-link="${roleName}"]`).exists('aws: role shows in the list');

    //and delete
    await click(`[data-test-secret-link="${roleName}"] [data-test-popup-menu-trigger]`);
    await waitUntil(() => find(`[data-test-aws-role-delete="${roleName}"]`)); // flaky without
    await click(`[data-test-aws-role-delete="${roleName}"]`);
    await click(GENERAL.confirmButton);
    assert.dom(`[data-test-secret-link="${roleName}"]`).doesNotExist('aws: role is no longer in the list');
  });
});
