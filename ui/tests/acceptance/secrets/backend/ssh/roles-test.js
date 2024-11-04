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
import listPage from 'vault/tests/pages/secrets/backend/list';
import editPage from 'vault/tests/pages/secrets/backend/ssh/edit-role';
import showPage from 'vault/tests/pages/secrets/backend/ssh/show';
import generatePage from 'vault/tests/pages/secrets/backend/ssh/generate-otp';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

// There is duplication within this test suite. The duplication occurred because two test suites both testing ssh roles were merged.
// refactoring the tests to remove the duplication would be a good next step as well as removing the tests/pages.

module('Acceptance | ssh | roles', function (hooks) {
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
        await click(GENERAL.inputByAttr('allowUserCertificates'));
        await click(GENERAL.toggleGroup('Options'));
        // it's recommended to keep allow_empty_principals false, check for testing so we don't have to input an extra field when signing a key
        await click(GENERAL.inputByAttr('allowEmptyPrincipals'));
      },
      async fillInGenerate() {
        await fillIn(GENERAL.inputByAttr('publicKey'), PUB_KEY);
        await click('[data-test-toggle-button]');

        await click(GENERAL.ttl.toggle('TTL'));
        await fillIn(GENERAL.selectByAttr('ttl-unit'), 'm');

        document.querySelector(GENERAL.ttl.input('TTL')).value = 30;
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
        await fillIn(GENERAL.inputByAttr('defaultUser'), 'admin');
        await click(GENERAL.toggleGroup('Options'));
        await fillIn(GENERAL.inputByAttr('cidrList'), '1.2.3.4/32');
      },
      async fillInGenerate() {
        await fillIn(GENERAL.inputByAttr('username'), 'admin');
        await fillIn(GENERAL.inputByAttr('ip'), '1.2.3.4');
      },
      assertAfterGenerate(assert, sshPath) {
        assert.strictEqual(
          currentURL(),
          `/vault/secrets/${sshPath}/credentials/${this.name}`,
          'otp credential url is correct'
        );
        assert.dom(GENERAL.infoRowLabel('Key')).exists({ count: 1 }, 'renders the key');
        assert.dom('[data-test-masked-input]').exists({ count: 1 }, 'renders mask for key value');
        assert.dom(GENERAL.infoRowLabel('Port')).exists({ count: 1 }, 'renders the port');
        assert.dom('[data-test-row-value="Port"]').exists({ count: 1 }, "renders the port's value");
      },
    },
  ];

  test('it creates roles, generates keys and deletes roles', async function (assert) {
    assert.expect(28);
    const sshPath = `ssh-${this.uid}`;
    await enablePage.enable('ssh', sshPath);
    await click(SES.configTab);
    await click(SES.configure);
    // default has generate CA checked so we just submit the form
    await click(SES.ssh.save);
    // There is a delay in the backend for the public key to be generated, wait for it to complete by checking that the public key is displayed
    await waitFor(GENERAL.infoRowLabel('Public key'));
    await click(GENERAL.tab(sshPath));
    for (const role of ROLES) {
      // create a role
      await click(SES.createSecret);
      assert.dom(SES.secretHeader).includesText('SSH Role', `${role.type}: renders the create page`);

      await fillIn(GENERAL.inputByAttr('name'), role.name);
      await fillIn(GENERAL.inputByAttr('keyType'), role.type);
      await role.fillInCreate();
      await settled();

      // save the role
      await click(SES.ssh.createRole);
      await waitUntil(() => currentURL() === `/vault/secrets/${sshPath}/show/${role.name}`); // flaky without this
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${sshPath}/show/${role.name}`,
        `${role.type}: navigates to the show page on creation`
      );

      // sign a key with this role
      await click(SES.generateLink);
      assert.strictEqual(currentRouteName(), role.credsRoute, 'navigates to the credentials page');
      await role.fillInGenerate();
      if (role.type === 'ca') {
        await settled();
        role.assertBeforeGenerate(assert);
      }

      // generate creds
      await click(GENERAL.saveButton);
      await settled(); // eslint-disable-line
      role.assertAfterGenerate(assert, sshPath);

      await click(GENERAL.backButton);
      assert.dom('[data-test-secret-generate-form]').exists(`${role.type}: back takes you back to the form`);

      await click(GENERAL.cancelButton);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${sshPath}/list`,
        `${role.type}: cancel takes you to ssh index`
      );
      assert.dom(SES.secretLink(role.name)).exists(`${role.type}: role shows in the list`);
      const secret = listPage.secrets.objectAt(0);
      await secret.menuToggle();
      assert.dom('.hds-dropdown li').exists({ count: 5 }, 'Renders 5 popup menu items');

      // and delete from the popup list menu
      await waitUntil(() => find(SES.ssh.deleteRole)); // flaky without
      await click(SES.ssh.deleteRole);
      await click(GENERAL.confirmButton);
      assert.dom(SES.secretLink(role.name)).doesNotExist(`${role.type}: role is no longer in the list`);
    }
    // cleanup
    await runCmd(`delete sys/mounts/${sshPath}`);
  });
  module('Acceptance | ssh | otp role', function () {
    const createOTPRole = async (name) => {
      await fillIn(GENERAL.inputByAttr('name'), name);
      await fillIn(GENERAL.inputByAttr('keyType'), name);
      await click(GENERAL.toggleGroup('Options'));
      await fillIn(GENERAL.inputByAttr('keyType'), 'otp');
      await fillIn(GENERAL.inputByAttr('defaultUser'), 'admin');
      await fillIn(GENERAL.inputByAttr('cidrList'), '0.0.0.0/0');
      await click(SES.ssh.createRole);
    };
    test('it deletes a role from list view', async function (assert) {
      assert.expect(2);
      const path = `ssh-${this.uid}`;
      await enablePage.enable('ssh', path);
      await settled();
      await editPage.visitRoot({ backend: path });
      await createOTPRole('role');
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
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });

    test('it generates an OTP', async function (assert) {
      assert.expect(6);
      const path = `ssh-${this.uid}`;
      await enablePage.enable('ssh', path);
      await settled();
      await editPage.visitRoot({ backend: path });
      await createOTPRole('role');
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

      await fillIn(GENERAL.inputByAttr('username'), 'admin');
      await fillIn(GENERAL.inputByAttr('ip'), '192.168.1.1');
      await click(GENERAL.saveButton);
      assert.ok(generatePage.warningIsPresent, 'shows warning');
      await click(GENERAL.backButton);
      assert.ok(generatePage.userIsPresent, 'clears generate, shows user input');
      assert.ok(generatePage.ipIsPresent, 'clears generate, shows ip input');
      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });
  });
});
