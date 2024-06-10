/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, find, settled, waitUntil, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import { GENERAL } from '../helpers/general-selectors';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { overrideResponse } from 'vault/tests/helpers/stubs';

const AWS_CREDS = {
  configTab: '[data-test-configuration-tab]',
  configure: '[data-test-secret-backend-configure]',
  awsForm: '[data-test-aws-root-creds-form]',
  viewBackend: '[data-test-backend-view-link]',
  createSecret: '[data-test-secret-create]',
  secretHeader: '[data-test-secret-header]',
  secretLink: (name) => (name ? `[data-test-secret-link="${name}"]` : '[data-test-secret-link]'),
  crumb: (path) => `[data-test-secret-breadcrumb="${path}"] a`,
  ttlToggle: '[data-test-ttl-toggle="TTL"]',
  warning: '[data-test-warning]',
  delete: (role) => `[data-test-aws-role-delete="${role}"]`,
  backButton: '[data-test-back-button]',
  generateLink: '[data-test-backend-credentials]',
};
const ROLE_TYPES = [
  {
    credentialType: 'iam_user',
    async fillOutForm(assert) {
      // nothing to fill out
      assert.dom('[data-test-field]').exists({ count: 1 });
    },
    expectedPayload: {},
  },
  {
    credentialType: 'assumed_role',
    async fillOutForm(assert) {
      await click(GENERAL.toggleInput('TTL'));
      assert.dom(GENERAL.toggleInput('TTL')).isNotChecked();
      await fillIn(GENERAL.inputByAttr('roleArn'), 'foobar');
    },
    expectedPayload: {
      role_arn: 'foobar',
    },
  },
  {
    credentialType: 'federation_token',
    async fillOutForm(assert) {
      assert.dom(GENERAL.toggleInput('TTL')).isChecked();
      await fillIn(GENERAL.ttl.input('TTL'), '3');
    },
    expectedPayload: {
      ttl: '10800s',
    },
  },
  {
    credentialType: 'session_token',
    async fillOutForm(assert) {
      await click(GENERAL.toggleInput('TTL'));
      assert.dom(GENERAL.toggleInput('TTL')).isNotChecked();
    },
    expectedPayload: null,
  },
];
module('Acceptance | aws secret backend', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');

    this.uid = uuidv4();
    return authPage.login();
  });

  test('aws backend', async function (assert) {
    const path = `aws-${this.uid}`;
    const roleName = 'awsrole';

    await enablePage.enable('aws', path);
    await settled();
    await click(AWS_CREDS.configTab);

    await click(AWS_CREDS.configure);

    assert.strictEqual(currentURL(), `/vault/settings/secrets/configure/${path}`);

    assert.dom(AWS_CREDS.awsForm).exists();
    assert.dom(GENERAL.tab('access-to-aws')).exists('renders the root creds tab');
    assert.dom(GENERAL.tab('lease')).exists('renders the leases config tab');

    await fillIn(GENERAL.inputByAttr('accessKey'), 'foo');
    await fillIn(GENERAL.inputByAttr('secretKey'), 'bar');

    await click(GENERAL.saveButton);

    assert.true(
      this.flashSuccessSpy.calledWith('The backend configuration saved successfully!'),
      'success flash message is rendered'
    );

    await click(GENERAL.tab('lease'));

    await click(GENERAL.saveButton);

    assert.true(
      this.flashSuccessSpy.calledTwice,
      'a new success flash message is rendered upon saving lease'
    );

    await click(AWS_CREDS.viewBackend);

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`, 'navigates to the roles list');

    await click(AWS_CREDS.createSecret);

    assert.dom(AWS_CREDS.secretHeader).hasText('Create an AWS Role', 'aws: renders the create page');

    await fillIn(GENERAL.inputByAttr('name'), roleName);

    // save the role
    await click(GENERAL.saveButton);
    await waitUntil(() => currentURL() === `/vault/secrets/${path}/show/${roleName}`); // flaky without this
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/show/${roleName}`,
      'aws: navigates to the show page on creation'
    );
    await click(AWS_CREDS.crumb(path));

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/list`);
    assert.dom(AWS_CREDS.secretLink(roleName)).exists();

    //and delete
    await click(`${AWS_CREDS.secretLink(roleName)} [data-test-popup-menu-trigger]`);
    await waitUntil(() => find(AWS_CREDS.delete(roleName))); // flaky without
    await click(AWS_CREDS.delete(roleName));
    await click(GENERAL.confirmButton);
    assert.dom(AWS_CREDS.secretLink(roleName)).doesNotExist('aws: role is no longer in the list');
  });

  ROLE_TYPES.forEach((scenario) => {
    test(`aws credentials - type ${scenario.credentialType}`, async function (assert) {
      const path = `aws-cred-${this.uid}`;
      const roleName = `awsrole-${scenario.credentialType}`;
      this.server.post(`/${path}/creds/${roleName}`, (_, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.deepEqual(payload, scenario.expectedPayload);
        return {
          data: {
            access_key: 'AKIA...',
            secret_key: 'xlCs...',
            security_token: 'some-token',
            arn: 'arn:aws:sts::123456789012:assumed-role/DeveloperRole/some-user-supplied-role-session-name',
          },
        };
      });
      this.server.get(`/${path}/creds/${roleName}`, () => {
        return {
          data: {
            access_key: 'AKIA...',
            secret_key: 'xlCs...',
            security_token: 'some-token',
            arn: 'arn:aws:sts::123456789012:assumed-role/DeveloperRole/some-user-supplied-role-session-name',
          },
        };
      });
      await runCmd(mountEngineCmd('aws', path));

      await visit(`/vault/secrets/${path}/create`);
      assert.dom('h1').hasText('Create an AWS Role');
      await fillIn(GENERAL.inputByAttr('name'), roleName);
      await fillIn(GENERAL.inputByAttr('credentialType'), scenario.credentialType);
      await click(GENERAL.saveButton);
      await waitUntil(() => currentURL() === `/vault/secrets/${path}/show/${roleName}`); // flaky without this
      assert.strictEqual(currentURL(), `/vault/secrets/${path}/show/${roleName}`);
      await click(AWS_CREDS.generateLink);
      assert
        .dom(GENERAL.inputByAttr('credentialType'))
        .hasValue(scenario.credentialType, 'credentialType matches backing role');

      // based on credentialType, fill out form
      await scenario.fillOutForm(assert);

      await click(GENERAL.saveButton);
      assert.dom(AWS_CREDS.warning).exists('Shows access warning after generation');
      assert.dom(GENERAL.infoRowValue('Access key')).exists();
      assert.dom(GENERAL.infoRowValue('Secret key')).exists();
      assert.dom(GENERAL.infoRowValue('Security token')).exists();
      await visit('/vault/dashboard');

      await runCmd(deleteEngineCmd(path));
    });
  });

  test(`aws credentials without role read access`, async function (assert) {
    const path = `aws-cred-${this.uid}`;
    const roleName = `awsrole-noread`;
    this.server.post(`/${path}/creds/${roleName}`, () => {
      return {
        data: {
          access_key: 'AKIA...',
          secret_key: 'xlCs...',
          security_token: 'some-token',
          arn: 'arn:aws:sts::123456789012:assumed-role/DeveloperRole/some-user-supplied-role-session-name',
        },
      };
    });
    this.server.get(`/${path}/roles/${roleName}`, () => overrideResponse(403));
    await runCmd(mountEngineCmd('aws', path));
    await runCmd(`write ${path}/roles/${roleName} credential_type=assumed_role`);

    await visit(`/vault/secrets/${path}/list`);
    assert.dom(AWS_CREDS.secretLink(roleName)).exists();
    await click(AWS_CREDS.secretLink(roleName));

    assert.strictEqual(currentURL(), `/vault/secrets/${path}/credentials/${roleName}`);
    assert
      .dom(GENERAL.inputByAttr('credentialType'))
      .hasValue('iam_user', 'credentialType defaults to first in list due to no role read permissions');

    await fillIn(GENERAL.inputByAttr('credentialType'), 'assumed_role');

    await click(GENERAL.saveButton);
    assert.dom(AWS_CREDS.warning).exists('Shows access warning after generation');
    assert.dom(GENERAL.infoRowValue('Access key')).exists();
    assert.dom(GENERAL.infoRowValue('Secret key')).exists();
    assert.dom(GENERAL.infoRowValue('Security token')).exists();
    await visit('/vault/dashboard');

    await runCmd(deleteEngineCmd(path));
  });
});
