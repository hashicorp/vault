/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, visit, currentURL } from '@ember/test-helpers';
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
import {
  createConfig,
  expectedConfigKeys,
  expectedValueOfConfigKeys,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

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

  test('it should prompt configuration after mounting the engine', async function (assert) {
    const path = `aws-${this.uid}`;
    // in this test go through the full mount process. Bypass this step in later tests.
    await visit('/vault/settings/mount-secret-backend');
    await click(SES.mountType('aws'));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(SES.mountSubmit);
    await click(SES.configTab);

    assert.dom(GENERAL.emptyStateTitle).hasText('AWS not configured');
    assert.dom(GENERAL.emptyStateActions).hasText('Configure AWS');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
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
    await enablePage.enable('aws', path);
    await click(SES.configTab);
    await click(SES.configure);
    await click(GENERAL.hdsTab('lease'));
    await click(GENERAL.toggleInput('Lease'));
    await fillIn(GENERAL.ttl.input('Lease'), '55');
    await click(GENERAL.toggleInput('Maximum Lease'));
    await fillIn(GENERAL.ttl.input('Maximum Lease'), '65');
    this.server.post(`${path}/config/lease`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.deepEqual(payload.lease, '55s', 'lease is set to 55s');
      assert.deepEqual(payload.lease_max, '65s', 'maximum_lease is set to 65s');
      return { data: { id: path, type: 'aws', attributes: payload } };
    });

    await click(GENERAL.saveButtonId('lease'));
    assert.true(
      this.flashSuccessSpy.calledWith('The backend configuration saved successfully!'),
      'Success flash message is rendered'
    );
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it show AWS configuration details', async function (assert) {
    // TODO: with WIF project will show Lease details as well.
    assert.expect(12);
    const path = `aws-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    createConfig(this.store, path, type); // create the aws root config in the store
    this.server.get(`${path}/config/root`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.ok(true, 'request made to config/root when navigating to the configuration page.');
      return { data: { id: path, type, attributes: payload } };
    });
    await click(SES.configTab);
    for (const key of expectedConfigKeys(type)) {
      assert.dom(GENERAL.infoRowLabel(key)).exists(`key for ${key} on the ${type} config details exists.`);
      const responseKeyAndValue = expectedValueOfConfigKeys(type, key);
      assert
        .dom(GENERAL.infoRowValue(key))
        .hasText(responseKeyAndValue, `value for ${key} on the ${type} config details exists.`);
    }
    // check mount configuration details is present and accurate.
    await click(SES.configurationToggle);
    assert
      .dom(GENERAL.infoRowValue('Path'))
      .hasText(`${path}/`, 'mount path is displayed in the configuration details');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should update AWS configuration details after editing', async function (assert) {
    // TODO: with WIF project will show Lease details as well.
    assert.expect(4);
    const path = `aws-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    // create accessKey with value foo and confirm it shows up in the details page.
    await click(SES.configTab);
    await click(SES.configure);
    await fillIn(GENERAL.inputByAttr('accessKey'), 'foo');
    await click(GENERAL.saveButtonId('root'));
    await click(SES.viewBackend);
    await click(SES.configTab);
    assert.dom(GENERAL.infoRowValue('Access key')).hasText('foo', 'Access key is foo');
    assert
      .dom(GENERAL.infoRowValue('Region'))
      .doesNotExist('Region has not been added therefor it does not show up on the details view.');
    // edit accessKey and another field and confirm the details page is updated.
    await click(SES.configure);
    await fillIn(GENERAL.inputByAttr('accessKey'), 'hello');
    await click(GENERAL.menuTrigger);
    await fillIn(GENERAL.selectByAttr('region'), 'ca-central-1');
    await click(GENERAL.saveButtonId('root'));
    await click(SES.viewBackend);
    await click(SES.configTab);
    assert.dom(GENERAL.infoRowValue('Access key')).hasText('hello', 'Access key has been updated to hello');
    assert.dom(GENERAL.infoRowValue('Region')).hasText('ca-central-1', 'Region has been added');
    // cleanup
  });

  // TODO once AWS configuration forms have been fixed, assert: transitions to configuration details page on cancel and on save, as well as breadcrumbs.
});
