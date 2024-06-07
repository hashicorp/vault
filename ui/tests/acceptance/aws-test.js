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
import { setupMirage } from 'ember-cli-mirage/test-support';

const AWS_CREDS = {
  configTab: '[data-test-configuration-tab]',
  configure: '[data-test-secret-backend-configure]',
  awsForm: '[data-test-aws-root-creds-form]',
  viewBackend: '[data-test-backend-view-link]',
  createSecret: '[data-test-secret-create]',
  secretHeader: '[data-test-secret-header]',
  secretLink: (name) => `[data-test-secret-link="${name}"]`,
  crumb: (path) => `[data-test-secret-breadcrumb="${path}"] a`,
  ttlToggle: '[data-test-ttl-toggle="TTL"]',
  warning: '[data-test-warning]',
  delete: (role) => `[data-test-aws-role-delete="${role}"]`,
  backButton: '[data-test-back-button]',
};
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
    this.server.post(`/${path}/creds/${roleName}`, (_, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.deepEqual(payload, { role_arn: 'foobar' }, 'does not send TTL when unchecked');
      return {
        data: {
          access_key: 'AKIA...',
          secret_key: 'xlCs...',
          security_token: 'some-token',
          arn: 'arn:aws:sts::123456789012:assumed-role/DeveloperRole/some-user-supplied-role-session-name',
        },
      };
    });

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

    // check that generates credentials flow is correct
    await click(AWS_CREDS.secretLink(roleName));
    assert.dom('h1').hasText('Generate AWS Credentials');
    assert.dom(GENERAL.inputByAttr('credentialType')).hasValue('iam_user');
    await fillIn(GENERAL.inputByAttr('credentialType'), 'assumed_role');
    await click(AWS_CREDS.ttlToggle);
    assert.dom(AWS_CREDS.ttlToggle).isNotChecked();
    await fillIn(GENERAL.inputByAttr('roleArn'), 'foobar');
    await click(GENERAL.saveButton);
    assert.dom(AWS_CREDS.warning).exists('Shows access warning after generation');
    assert.dom(GENERAL.infoRowValue('Access key')).exists();
    assert.dom(GENERAL.infoRowValue('Secret key')).exists();
    assert.dom(GENERAL.infoRowValue('Security token')).exists();
    await click(AWS_CREDS.backButton);

    //and delete
    await click(`${AWS_CREDS.secretLink(roleName)} [data-test-popup-menu-trigger]`);
    await waitUntil(() => find(AWS_CREDS.delete(roleName))); // flaky without
    await click(AWS_CREDS.delete(roleName));
    await click(GENERAL.confirmButton);
    assert.dom(AWS_CREDS.secretLink(roleName)).doesNotExist('aws: role is no longer in the list');
  });
});
