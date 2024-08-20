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
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import {
  createConfig,
  expectedConfigKeys,
  expectedValueOfConfigKeys,
  configUrl,
  fillInAwsConfig,
} from 'vault/tests/helpers/secret-engine/secret-engine-helpers';

module('Acceptance | aws | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const flash = this.owner.lookup('service:flash-messages');
    this.store = this.owner.lookup('service:store');
    this.flashSuccessSpy = spy(flash, 'success');
    this.flashInfoSpy = spy(flash, 'info');

    this.uid = uuidv4();
    return authPage.login();
  });

  test('it should prompt configuration after mounting the aws engine', async function (assert) {
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

  test('it should transition to configure page on click "Configure" from toolbar', async function (assert) {
    const path = `aws-${this.uid}`;
    await enablePage.enable('aws', path);
    await click(SES.configTab);
    await click(SES.configure);
    assert.strictEqual(currentURL(), `/vault/secrets/${path}/configuration/edit`);
    assert.dom(SES.configureTitle('aws')).hasText('Configure AWS');
    assert.dom(SES.aws.rootForm).exists('it lands on the root configuration form.');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should show error if old url is entered', async function (assert) {
    // we are intentionally not redirecting from the old url to the new one.
    const path = `aws-${this.uid}`;
    await enablePage.enable('aws', path);
    await click(SES.configTab);
    await visit(`/vault/settings/secrets/configure/${path}`);
    assert.dom(GENERAL.notFound).exists('shows page-error');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should save root AWS configuration', async function (assert) {
    assert.expect(3);
    const path = `aws-${this.uid}`;
    await enablePage.enable('aws', path);

    this.server.post(configUrl('aws-lease', path), () => {
      assert.false(true, 'post request was made to config/lease when no data was changed. test should fail.');
    });

    await click(SES.configTab);
    await click(SES.configure);
    await fillInAwsConfig();
    await click(SES.aws.save);
    assert.true(
      this.flashSuccessSpy.calledWith(`Successfully saved ${path}'s root configuration.`),
      'Success flash message is rendered showing the root configuration was saved.'
    );
    assert.dom(GENERAL.infoRowValue('Access key')).hasText('foo', 'Access Key has been set.');
    assert
      .dom(GENERAL.infoRowValue('Secret key'))
      .doesNotExist('Secret key is not shown because it does not get returned by the api.');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should save lease AWS configuration', async function (assert) {
    assert.expect(3);
    const path = `aws-${this.uid}`;
    await enablePage.enable('aws', path);

    this.server.post(configUrl('aws', path), () => {
      assert.false(true, 'post request was made to config/root when no data was changed. test should fail.');
    });
    await click(SES.configTab);
    await click(SES.configure);
    await fillInAwsConfig(false, false, true); // only fills in lease config with defaults
    await click(SES.aws.save);
    assert.true(
      this.flashSuccessSpy.calledWith(`Successfully saved ${path}'s lease configuration.`),
      'Success flash message is rendered showing the lease configuration was saved.'
    );

    assert.dom(GENERAL.infoRowValue('Default Lease TTL')).hasText('33s', `Default TTL has been set.`);
    assert.dom(GENERAL.infoRowValue('Max Lease TTL')).hasText('44s', `Max lease TTL has been set.`);
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it shows AWS mount configuration details', async function (assert) {
    assert.expect(12);
    const path = `aws-${this.uid}`;
    const type = 'aws';
    this.server.get(`${path}/config/root`, (schema, req) => {
      const payload = JSON.parse(req.requestBody);
      assert.ok(true, 'request made to config/root when navigating to the configuration page.');
      return { data: { id: path, type, attributes: payload } };
    });
    await enablePage.enable(type, path);
    createConfig(this.store, path, type); // create the aws root config in the store
    await click(SES.configTab);
    for (const key of expectedConfigKeys(type)) {
      assert.dom(GENERAL.infoRowLabel(key)).exists(`${key} on the ${type} config details exists.`);
      const responseKeyAndValue = expectedValueOfConfigKeys(type, key);
      assert
        .dom(GENERAL.infoRowValue(key))
        .hasText(responseKeyAndValue, `value for ${key} on the ${type} config details exists.`);
    }
    // check mount configuration details are present and accurate.
    await click(SES.configurationToggle);
    assert
      .dom(GENERAL.infoRowValue('Path'))
      .hasText(`${path}/`, 'mount path is displayed in the configuration details');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should update AWS configuration details after editing', async function (assert) {
    const path = `aws-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    // create accessKey with value foo and confirm it shows up in the details page.
    await click(SES.configTab);
    await click(SES.configure);
    await fillInAwsConfig();
    await click(SES.aws.save);
    assert.dom(GENERAL.infoRowValue('Access key')).hasText('foo', 'Access key is foo');
    assert
      .dom(GENERAL.infoRowValue('Region'))
      .doesNotExist('Region has not been added therefor it does not show up on the details view.');
    // edit root config details and lease config details and confirm the configuration.index page is updated.
    await click(SES.configure);
    // edit root config details
    await fillIn(GENERAL.inputByAttr('accessKey'), 'not-foo');
    await click(GENERAL.toggleGroup('Root config options'));
    await fillIn(GENERAL.inputByAttr('region'), 'ap-southeast-2');
    // add lease config details
    await fillInAwsConfig(false, false, true); // only fills in lease config with defaults
    await click(SES.aws.save);
    assert
      .dom(GENERAL.infoRowValue('Access key'))
      .hasText('not-foo', 'Access key has been updated to not-foo');
    assert.dom(GENERAL.infoRowValue('Region')).hasText('ap-southeast-2', 'Region has been added');
    assert.dom(GENERAL.infoRowValue('Default Lease TTL')).hasText('33s', 'Default Lease TTL has been added');
    assert.dom(GENERAL.infoRowValue('Max Lease TTL')).hasText('44s', 'Max Lease TTL has been added');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should show API error when AWS configuration read fails', async function (assert) {
    assert.expect(1);
    const path = `aws-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    // interrupt get and return API error
    this.server.get(configUrl(type, path), () => {
      return overrideResponse(400, { errors: ['bad request'] });
    });
    await click(SES.configTab);
    assert.dom(SES.error.title).hasText('Error', 'shows the secrets backend error route');
  });

  test('it should not make a post request if lease or root data was unchanged', async function (assert) {
    assert.expect(3);
    const path = `aws-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);

    this.server.post(configUrl(type, path), () => {
      assert.false(true, 'post request was made to config/root when no data was changed. test should fail.');
    });
    this.server.post(configUrl('aws-lease', path), () => {
      assert.false(true, 'post request was made to config/lease when no data was changed. test should fail.');
    });

    await click(SES.configTab);
    await click(SES.configure);
    await click(SES.aws.save);
    assert.true(
      this.flashInfoSpy.calledWith('No changes detected.'),
      'Flash message shows no changes detected.'
    );
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${path}/configuration`,
      'navigates back to the configuration index view'
    );
    assert.dom(GENERAL.emptyStateTitle).hasText('AWS not configured');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });

  test('it should reset models after saving', async function (assert) {
    const path = `aws-${this.uid}`;
    const type = 'aws';
    await enablePage.enable(type, path);
    await click(SES.configTab);
    await click(SES.configure);
    await fillInAwsConfig(true);
    //  the way to tell if a record has been unloaded is if the private key is not saved in the store (the API does not return it, but if the record was not unloaded it would have stayed.)
    await click(SES.aws.save); // save the configuration
    await click(SES.configure);
    const privateKeyExists = this.store.peekRecord('aws/root-config', path).privateKey ? true : false;
    assert.false(
      privateKeyExists,
      'private key is not on the store record, meaning it was unloaded after save. This new record without the key comes from the API.'
    );
    assert
      .dom(GENERAL.enableField('secretKey'))
      .exists('secret key field is wrapped inside an enableInput component');
    // cleanup
    await runCmd(`delete sys/mounts/${path}`);
  });
});
