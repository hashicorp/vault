/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, find, settled, waitUntil } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Acceptance | aws secret backend', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('aws backend', async function (assert) {
    const path = `aws-${this.uid}`;
    const roleName = 'awsrole';
    this.server.post(`/${path}/creds/${roleName}`, (_, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.deepEqual(payload, { role_arn: 'foobar' }, 'does not send TTL when unchecked');
      return {};
    });

    await enablePage.enable('aws', path);
    await settled();
    await click('[data-test-configuration-tab]');

    await click('[data-test-secret-backend-configure]');

    assert.strictEqual(currentURL(), `/vault/settings/secrets/configure/${path}`);
    assert.dom('[data-test-aws-root-creds-form]').exists();
    assert.dom('[data-test-aws-link="root-creds"]').exists();
    assert.dom('[data-test-aws-link="leases"]').exists();

    await fillIn('[data-test-aws-input="accessKey"]', 'foo');
    await fillIn('[data-test-aws-input="secretKey"]', 'bar');

    await click('[data-test-aws-input="root-save"]');

    assert
      .dom('[data-test-flash-message]:last-of-type [data-test-flash-message-body]')
      .includesText(`The backend configuration saved successfully!`);

    await click('[data-test-aws-link="leases"]');

    await click('[data-test-aws-input="lease-save"]');
    assert
      .dom('[data-test-flash-message]:last-of-type [data-test-flash-message-body]')
      .includesText(`The backend configuration saved successfully!`);

    await click('[data-test-backend-view-link]');

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`, `navigates to the roles list`);

    await click('[data-test-secret-create]');

    assert.dom('[data-test-secret-header]').includesText('AWS Role');

    await fillIn('[data-test-input="name"]', roleName);

    // save the role
    await click('[data-test-role-aws-create]');
    await waitUntil(() => currentURL() === `/vault/secrets/${path}/show/${roleName}`); // flaky without this
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      `$aws: navigates to the show page on creation`
    );
    await click(`[data-test-secret-breadcrumb="${path}"] a`);

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`);
    assert.dom(`[data-test-secret-link="${roleName}"]`).exists();

    // check that generates credentials flow is correct
    await click(`[data-test-secret-link="${roleName}"]`);
    assert.dom('h1').hasText('Generate AWS Credentials');
    assert.dom('[data-test-input="credentialType"]').hasValue('iam_user');
    await fillIn('[data-test-input="credentialType"]', 'assumed_role');
    await click('[data-test-ttl-toggle="TTL"]');
    assert.dom('[data-test-ttl-toggle="TTL"]').isNotChecked();
    await fillIn('[data-test-input="roleArn"]', 'foobar');
    await click('[data-test-secret-generate]');
    assert.dom('[data-test-warning]').exists('Shows access warning after generation');
    await click('[data-test-secret-generate-back]');

    //and delete
    await click(`[data-test-secret-link="${roleName}"] [data-test-popup-menu-trigger]`);
    await waitUntil(() => find(`[data-test-aws-role-delete="${roleName}"]`)); // flaky without
    await click(`[data-test-aws-role-delete="${roleName}"]`);
    await click(`[data-test-confirm-button]`);
    assert.dom(`[data-test-secret-link="${roleName}"]`).doesNotExist(`aws: role is no longer in the list`);
  });
});
