/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  click,
  fillIn,
  currentURL,
  find,
  settled,
  waitUntil,
  currentRouteName,
  waitFor,
} from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | ssh secret backend', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  const PUB_KEY = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCn9p5dHNr4aU4R2W7ln+efzO5N2Cdv/SXk6zbCcvhWcblWMjkXf802B0PbKvf6cJIzM/Xalb3qz1cK+UUjCSEAQWefk6YmfzbOikfc5EHaSKUqDdE+HlsGPvl42rjCr28qYfuYh031YfwEQGEAIEypo7OyAj+38NLbHAQxDxuaReee1YCOV5rqWGtEgl2VtP5kG+QEBza4ZfeglS85f/GGTvZC4Jq1GX+wgmFxIPnd6/mUXa4ecoR0QMfOAzzvPm4ajcNCQORfHLQKAcmiBYMiyQJoU+fYpi9CJGT1jWTmR99yBkrSg6yitI2qqXyrpwAbhNGrM0Fw0WpWxh66N9Xp meirish@Macintosh-3.local`;

  const ROLES = [
    {
      type: 'ca',
      name: 'carole',
      credsRoute: 'vault.cluster.secrets.backend.sign',
      async fillInCreate() {
        await click('[data-test-input="allowUserCertificates"]');
      },
      async fillInGenerate() {
        await fillIn('[data-test-input="publicKey"]', PUB_KEY);
        await click('[data-test-toggle-button]');

        await click('[data-test-toggle-label="TTL"]');
        await fillIn('[data-test-select="ttl-unit"]', 'm');

        document.querySelector('[data-test-ttl-value="TTL"]').value = 30;
      },
      assertBeforeGenerate(assert) {
        assert.dom('[data-test-form-field-from-model]').exists('renders the FormFieldFromModel');
        const value = document.querySelector('[data-test-ttl-value="TTL"]').value;
        // confirms that the actions are correctly being passed down to the FormFieldFromModel component
        assert.strictEqual(value, '30', 'renders action updateTtl');
      },
      assertAfterGenerate(assert, sshPath) {
        assert.strictEqual(
          currentURL(),
          `/vault/secrets/${sshPath}/sign/${this.name}`,
          'ca sign url is correct'
        );
        assert.dom('[data-test-row-label="Signed key"]').exists({ count: 1 }, 'renders the signed key');
        assert
          .dom('[data-test-row-value="Signed key"]')
          .exists({ count: 1 }, "renders the signed key's value");
        assert.dom('[data-test-row-label="Serial number"]').exists({ count: 1 }, 'renders the serial');
        assert.dom('[data-test-row-value="Serial number"]').exists({ count: 1 }, 'renders the serial value');
      },
    },
    {
      type: 'otp',
      name: 'otprole',
      credsRoute: 'vault.cluster.secrets.backend.credentials',
      async fillInCreate() {
        await fillIn('[data-test-input="defaultUser"]', 'admin');
        await click('[data-test-toggle-group="Options"]');
        await fillIn('[data-test-input="cidrList"]', '1.2.3.4/32');
      },
      async fillInGenerate() {
        await fillIn('[data-test-input="username"]', 'admin');
        await fillIn('[data-test-input="ip"]', '1.2.3.4');
      },
      assertAfterGenerate(assert, sshPath) {
        assert.strictEqual(
          currentURL(),
          `/vault/secrets/${sshPath}/credentials/${this.name}`,
          'otp credential url is correct'
        );
        assert.dom('[data-test-row-label="Key"]').exists({ count: 1 }, 'renders the key');
        assert.dom('[data-test-masked-input]').exists({ count: 1 }, 'renders mask for key value');
        assert.dom('[data-test-row-label="Port"]').exists({ count: 1 }, 'renders the port');
        assert.dom('[data-test-row-value="Port"]').exists({ count: 1 }, "renders the port's value");
      },
    },
  ];
  test('ssh backend', async function (assert) {
    assert.expect(30);
    const sshPath = `ssh-${this.uid}`;

    await enablePage.enable('ssh', sshPath);
    await settled();
    await click('[data-test-configuration-tab]');

    await click('[data-test-secret-backend-configure]');

    assert.strictEqual(currentURL(), `/vault/settings/secrets/configure/${sshPath}`);
    assert.dom('[data-test-ssh-configure-form]').exists('renders the empty configuration form');

    // default has generate CA checked so we just submit the form
    await click('[data-test-ssh-input="configure-submit"]');

    await waitFor('[data-test-ssh-input="public-key"]');
    assert.dom('[data-test-ssh-input="public-key"]').exists();
    await click('[data-test-backend-view-link]');

    assert.strictEqual(currentURL(), `/vault/secrets/${sshPath}/list`, `redirects to ssh index`);

    for (const role of ROLES) {
      // create a role
      await click('[data-test-secret-create]');

      assert
        .dom('[data-test-secret-header]')
        .includesText('SSH Role', `${role.type}: renders the create page`);

      await fillIn('[data-test-input="name"]', role.name);
      await fillIn('[data-test-input="keyType"]', role.type);
      await role.fillInCreate();
      await settled();

      // save the role
      await click('[data-test-role-ssh-create]');
      await waitUntil(() => currentURL() === `/vault/secrets/${sshPath}/show/${role.name}`); // flaky without this
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${sshPath}/show/${role.name}`,
        `${role.type}: navigates to the show page on creation`
      );

      // sign a key with this role
      await click('[data-test-backend-credentials]');
      assert.strictEqual(currentRouteName(), role.credsRoute);
      await role.fillInGenerate();
      if (role.type === 'ca') {
        await settled();
        role.assertBeforeGenerate(assert);
      }

      // generate creds
      await click(GENERAL.saveButton);
      await settled(); // eslint-disable-line
      role.assertAfterGenerate(assert, sshPath);

      // click the "Back" button
      await click('[data-test-back-button]');

      assert.dom('[data-test-secret-generate-form]').exists(`${role.type}: back takes you back to the form`);

      await click(GENERAL.cancelButton);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${sshPath}/list`,
        `${role.type}: cancel takes you to ssh index`
      );
      assert.dom(`[data-test-secret-link="${role.name}"]`).exists(`${role.type}: role shows in the list`);

      //and delete
      await click(`[data-test-secret-link="${role.name}"] [data-test-popup-menu-trigger]`);
      await waitUntil(() => find('[data-test-ssh-role-delete]')); // flaky without
      await click(`[data-test-ssh-role-delete]`);
      await click(`[data-test-confirm-button]`);
      assert
        .dom(`[data-test-secret-link="${role.name}"]`)
        .doesNotExist(`${role.type}: role is no longer in the list`);
    }
  });
});
