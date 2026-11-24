/**
 * Copyright IBM Corp. 2016, 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL, find, waitUntil, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';

const ROLE_TYPES = [
  {
    credentialType: 'iam_user',
    async fillOutForm(assert) {
      // nothing to fill out
      assert
        .dom(GENERAL.helpText)
        .hasText('For Vault roles of credential type iam_user, there are no inputs, just submit the form.');
    },
    expectedPayload: {},
  },
  {
    credentialType: 'assumed_role',
    async fillOutForm(assert) {
      await click(GENERAL.toggleInput('TTL'));
      assert.dom(GENERAL.toggleInput('TTL')).isNotChecked();
      await fillIn(GENERAL.inputByAttr('role_arn'), 'foobar');
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
    expectedPayload: {},
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
    return login();
  });

  test('it creates role and deletes role', async function (assert) {
    const path = `aws-${this.uid}`;
    const roleName = 'awsrole';
    await enablePage.enable('aws', path);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${path}/list`,
      'After enabling aws secrets engine it navigates to roles list'
    );

    await click(SES.createSecretLink);
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create an AWS Role', 'It renders the create role page');

    await fillIn(GENERAL.inputByAttr('name'), roleName);
    await click(GENERAL.submitButton);
    await waitUntil(() => currentURL() === `/vault/secrets-engines/${path}/show/${roleName}`); // flaky without this
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${path}/show/${roleName}`,
      'aws: navigates to the show page on creation'
    );

    await click(SES.crumb(path));
    assert.strictEqual(currentURL(), `/vault/secrets-engines/${path}/list`);
    assert.dom(SES.secretLink(roleName)).exists();

    // delete role
    await click(`${SES.secretLink(roleName)} [data-test-popup-menu-trigger]`);
    await waitUntil(() => find(SES.aws.deleteRole(roleName))); // flaky without
    await click(SES.aws.deleteRole(roleName));
    await click(GENERAL.confirmButton);
    assert.dom(SES.secretLink(roleName)).doesNotExist('aws: role is no longer in the list');
  });

  test('role edits persist after revisiting the details page', async function (assert) {
    const path = `aws-edit-${this.uid}`;
    const roleName = 'editable-role';
    await runCmd(mountEngineCmd('aws', path));
    await runCmd(`write ${path}/roles/${roleName} credential_type=iam_user`);
    await visit(`/vault/secrets-engines/${path}/edit/${roleName}`);
    assert.dom(GENERAL.inputByAttr('name')).doesNotExist('Role names cannot be changed during editing');
    await fillIn(GENERAL.inputByAttr('credential_type'), 'session_token');
    await click(GENERAL.submitButton);
    assert.strictEqual(currentURL(), `/vault/secrets-engines/${path}/show/${roleName}`);
    await visit('/vault/dashboard');
    await visit(`/vault/secrets-engines/${path}/show/${roleName}`);
    assert.dom(GENERAL.infoRowValue('Credential type')).hasText('session_token');
    await visit('/vault/dashboard');
    await runCmd(deleteEngineCmd(path));
  });

  for (const [credentialType, arnField, arn] of [
    ['iam_user', 'policy_arns', 'arn:aws:iam::aws:policy/ReadOnlyAccess'],
    ['assumed_role', 'role_arns', 'arn:aws:iam::123456789012:role/Developer'],
  ]) {
    test(`${credentialType} role ARN and policy widget changes persist after reload`, async function (assert) {
      const path = `aws-widgets-${this.uid}`;
      const roleName = 'widget-role';
      const policy = {
        Version: '2012-10-17',
        Statement: [{ Effect: 'Allow', Action: 's3:GetObject', Resource: '*' }],
      };
      await runCmd(mountEngineCmd('aws', path));
      try {
        await visit(`/vault/secrets-engines/${path}/create`);
        await fillIn(GENERAL.inputByAttr('name'), roleName);
        await fillIn(GENERAL.inputByAttr('credential_type'), credentialType);
        await fillIn(GENERAL.stringListByIdx(0), arn);
        setCodeEditorValue(codemirror(), JSON.stringify(policy));
        await click(GENERAL.submitButton);
        assert.strictEqual(
          currentURL(),
          `/vault/secrets-engines/${path}/show/${roleName}`,
          'Save navigates to details'
        );
        await visit('/vault/dashboard');
        await visit(`/vault/secrets-engines/${path}/show/${roleName}`);
        assert
          .dom(GENERAL.infoRowValue(arnField === 'policy_arns' ? 'Policy ARNs' : 'Role ARNs'))
          .includesText(arn, 'Reloaded details show the persisted ARN');
        assert.deepEqual(
          JSON.parse(find(GENERAL.infoRowValue('Policy Document')).textContent),
          policy,
          'Reloaded details show the persisted policy'
        );
        await visit(`/vault/secrets-engines/${path}/edit/${roleName}`);
        const updatedArn = arn.replace(/ReadOnlyAccess|Developer/, 'Updated');
        await fillIn(GENERAL.stringListByIdx(0), updatedArn);
        policy.Statement[0].Action = 's3:ListBucket';
        setCodeEditorValue(codemirror(), JSON.stringify(policy));
        await click(GENERAL.submitButton);
        await visit('/vault/dashboard');
        await visit(`/vault/secrets-engines/${path}/show/${roleName}`);
        assert
          .dom(GENERAL.infoRowValue(arnField === 'policy_arns' ? 'Policy ARNs' : 'Role ARNs'))
          .includesText(updatedArn, 'Editing persists the updated ARN');
        assert.deepEqual(
          JSON.parse(find(GENERAL.infoRowValue('Policy Document')).textContent),
          policy,
          'Editing persists the updated policy'
        );
        await click(GENERAL.confirmTrigger);
        await click(GENERAL.confirmButton);
        assert.strictEqual(
          currentURL(),
          `/vault/secrets-engines/${path}/list`,
          'Details deletion returns to the list'
        );
        await visit('/vault/dashboard');
        await visit(`/vault/secrets-engines/${path}/list`);
        assert.dom(SES.secretLink(roleName)).doesNotExist('The role stays deleted after revisiting the list');
      } finally {
        await visit('/vault/dashboard');
        await runCmd(deleteEngineCmd(path));
      }
    });
  }

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

      await visit(`/vault/secrets-engines/${path}/create`);
      assert.dom('h1').hasText('Create an AWS Role');
      await fillIn(GENERAL.inputByAttr('name'), roleName);
      await fillIn(GENERAL.inputByAttr('credential_type'), scenario.credentialType);
      await click(GENERAL.submitButton);
      await waitUntil(() => currentURL() === `/vault/secrets-engines/${path}/show/${roleName}`); // flaky without this
      assert.strictEqual(currentURL(), `/vault/secrets-engines/${path}/show/${roleName}`);
      await click(SES.generateLink);
      assert
        .dom(`[data-test-credential-type=${scenario.credentialType}]`)
        .exists(scenario.credentialType, 'credentialType matches backing role');

      // based on credentialType, fill out form
      await scenario.fillOutForm(assert);

      await click(GENERAL.submitButton);
      assert.dom(SES.warning).exists('Shows access warning after generation');
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

    await visit(`/vault/secrets-engines/${path}/list`);
    assert.dom(SES.secretLink(roleName)).exists();
    await click(SES.secretLink(roleName));

    assert.strictEqual(currentURL(), `/vault/secrets-engines/${path}/credentials/${roleName}`);
    assert
      .dom(GENERAL.inputByAttr('credential_type'))
      .hasValue('iam_user', 'credentialType defaults to first in list due to no role read permissions');
    await fillIn(GENERAL.inputByAttr('credential_type'), 'assumed_role');

    await click(GENERAL.submitButton);
    assert.dom(SES.warning).exists('Shows access warning after generation');
    assert.dom(GENERAL.infoRowValue('Access key')).exists();
    assert.dom(GENERAL.infoRowValue('Secret key')).exists();
    assert.dom(GENERAL.infoRowValue('Security token')).exists();
    await visit('/vault/dashboard');

    await runCmd(deleteEngineCmd(path));
  });
});

module('Acceptance | aws secret backend | restricted tokens', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.path = `aws-acl-${uuidv4()}`;
    this.policyName = `${this.path}-policy`;
    this.roleName = 'restricted-role';
    await login();
    await runCmd([
      mountEngineCmd('aws', this.path),
      `write ${this.path}/roles/${this.roleName} credential_type=assumed_role role_arns=arn:aws:iam::123456789012:role/Allowed`,
    ]);
  });

  hooks.afterEach(async function () {
    await visit('/vault/dashboard');
    await login();
    await runCmd([
      deleteEngineCmd(this.path),
      `delete sys/policies/acl/${this.policyName}`,
      `write auth/token/revoke token=${this.token}`,
    ]);
  });

  test('read-only tokens see role details but no mutation or generation actions', async function (assert) {
    this.token = await runCmd(
      tokenWithPolicyCmd(
        this.policyName,
        `
        path "${this.path}/roles" { capabilities = ["list"] }
        path "${this.path}/roles/*" { capabilities = ["read"] }
      `
      )
    );
    await login(this.token);
    await visit(`/vault/secrets-engines/${this.path}/list`);
    assert.dom(SES.secretLink(this.roleName)).exists('Lists roles with list access');
    await click(`${SES.secretLink(this.roleName)} ${GENERAL.menuTrigger}`);
    assert.dom('[data-test-role-ssh-link="show"]').exists('Read permission exposes the Details action');
    assert.dom('[data-test-role-ssh-link="edit"]').doesNotExist('Update permission is required to edit');
    assert.dom(SES.aws.deleteRole(this.roleName)).doesNotExist('Delete permission is required to delete');
    assert
      .dom('[data-test-role-aws-link="generate"]')
      .doesNotExist('Credential permission is required to generate');
    await click('[data-test-role-ssh-link="show"]');
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.path}/show/${this.roleName}`,
      'Details action opens the role'
    );
    assert
      .dom(GENERAL.infoRowValue('Credential type'))
      .hasText('assumed_role', 'Reads the persisted role as the restricted user');
    assert.dom(GENERAL.confirmTrigger).doesNotExist('Delete is unavailable without delete capability');
    assert.dom(SES.generateLink).doesNotExist('Generate is unavailable without credential capability');
    assert
      .dom(`a[href$="/edit/${this.roleName}"]`)
      .doesNotExist('Edit is unavailable without update capability');
    const api = this.owner.lookup('service:api');
    await assert.rejects(
      api.secrets.awsWriteRole(this.roleName, this.path, { credential_type: 'session_token' }),
      'Vault denies role updates'
    );
    await assert.rejects(api.secrets.awsDeleteRole(this.roleName, this.path), 'Vault denies role deletion');
  });

  test('update-only credential access allows manual type selection without role read', async function (assert) {
    this.token = await runCmd(
      tokenWithPolicyCmd(
        this.policyName,
        `
        path "${this.path}/roles" { capabilities = ["list"] }
        path "${this.path}/creds/*" { capabilities = ["update"] }
      `
      )
    );
    await login(this.token);
    await visit(`/vault/secrets-engines/${this.path}/list`);
    await click(SES.secretLink(this.roleName));
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.path}/credentials/${this.roleName}`,
      'The role link opens generation without read access'
    );
    assert
      .dom(GENERAL.inputByAttr('credential_type'))
      .hasValue('iam_user', 'Read denial leaves type selection available');
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-message-error]')
      .includesText('permission denied', 'Update-only access does not authorize IAM GET requests');
    await fillIn(GENERAL.inputByAttr('credential_type'), 'assumed_role');
    await fillIn(GENERAL.inputByAttr('role_arn'), 'arn:aws:iam::123456789012:role/NotAllowed');
    await click(GENERAL.submitButton);
    assert
      .dom('[data-test-message-error]')
      .includesText(
        'not in allowed role arns',
        'POST reaches AWS backend validation with update-only access'
      );
    assert
      .dom(GENERAL.inputByAttr('role_arn'))
      .hasValue('arn:aws:iam::123456789012:role/NotAllowed', 'Retains inputs after backend validation');
    assert.dom(GENERAL.submitButton).isNotDisabled('Can retry after the backend rejects the ARN');
  });
});
