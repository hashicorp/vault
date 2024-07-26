/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { spy } from 'sinon';

import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

module('Acceptance | aws | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.store = this.owner.lookup('service:store');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashDangerSpy = spy(flash, 'danger');

    this.uid = uuidv4();
    return authPage.login();
  });

  test('it should transition to configure page on Configure click from toolbar', async function (assert) {
    const path = `aws-${this.uid}`;
    await enablePage.enable('aws', path);
    await click(SES.configTab);
    await click(SES.configure);
    assert.strictEqual(currentURL(), `/vault/settings/secrets/configure/${path}`);
    assert.dom(SES.configureTitle('aws')).hasText('Configure AWS');
    assert.dom(SES.aws.rootForm).exists('it lands on the root configuration form.');
    assert.dom(GENERAL.tab('access-to-aws')).exists('renders the root creds tab');
    assert.dom(GENERAL.tab('lease')).exists('renders the leases config tab');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should save root AWS configuration', async function (assert) {
    assert.expect(3);
    const path = `aws-${this.uid}`;
    await enablePage.enable('aws', path);
    await click(SES.configTab);
    await click(SES.configure);
    await fillIn(GENERAL.inputByAttr('accessKey'), 'foo');
    await fillIn(GENERAL.inputByAttr('secretKey'), 'bar');
    this.server.post(`${path}/config/root`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.deepEqual(payload.access_key, 'foo', 'access_key is foo');
      assert.deepEqual(payload.secret_key, 'bar', 'secret_key is foo');
      return { data: { id: path, type: 'aws', attributes: payload } };
    });

    await click(GENERAL.saveButtonId('root'));
    assert.true(
      this.flashSuccessSpy.calledWith('The backend configuration saved successfully!'),
      'Success flash message is rendered'
    );
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should save lease AWS configuration', async function (assert) {
    assert.expect(3);
    const path = `aws-${this.uid}`;
    this.server.post(`${path}/config/lease`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.deepEqual(payload.lease, '55s', 'lease is set to 55s');
      assert.deepEqual(payload.lease_max, '65s', 'maximum_lease is set to 65s');
      return { data: { id: path, type: 'aws', attributes: payload } };
    });
    await enablePage.enable('aws', path);
    await click(SES.configTab);
    await click(SES.configure);
    await click(GENERAL.hdsTab('lease'));
    await click(GENERAL.toggleInput('Lease'));
    await fillIn(GENERAL.ttl.input('Lease'), '55');
    await click(GENERAL.toggleInput('Maximum Lease'));
    await fillIn(GENERAL.ttl.input('Maximum Lease'), '65');
    await click(GENERAL.saveButtonId('lease'));
    assert.true(
      this.flashSuccessSpy.calledWith('The backend configuration saved successfully!'),
      'Success flash message is rendered'
    );
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });
});
