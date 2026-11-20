/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { create } from 'ember-cli-page-object';
import { fillIn, settled, findAll, click } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';

import enablePage from 'vault/tests/pages/settings/auth/enable';
import page from 'vault/tests/pages/settings/auth/configure/section';
import indexPage from 'vault/tests/pages/settings/auth/configure/index';
import consolePanel from 'vault/tests/pages/components/console/ui-panel';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { TTL_PICKER } from 'vault/tests/helpers/components/ttl-picker-selectors';

import { GENERAL } from 'vault/tests/helpers/general-selectors';

const cli = create(consolePanel);

module('Acceptance | settings/auth/configure/section', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return login();
  });

  test('it can save options', async function (assert) {
    assert.expect(6);

    const path = `approle-save-${this.uid}`;
    const type = 'approle';
    const section = 'options';

    this.server.post(`/sys/mounts/auth/${path}/tune`, function (schema, request) {
      const body = JSON.parse(request.requestBody);
      const keys = Object.keys(body);
      assert.strictEqual(body.token_type, 'batch', 'passes new token type');
      assert.true(keys.includes('default_lease_ttl'), 'passes default_lease_ttl on tune');
      assert.true(keys.includes('max_lease_ttl'), 'passes max_lease_ttl on tune');
      assert.true(keys.includes('description'), 'passes updated description on tune');
      return request.passthrough();
    });

    await enablePage.enable(type, path);
    await page.visit({ path, section });
    await fillIn(GENERAL.inputByAttr('description'), 'This is Approle!');
    assert
      .dom(GENERAL.inputByAttr('config.token_type'))
      .hasValue('default-service', 'as default the token type selected is default-service.');
    await fillIn(GENERAL.inputByAttr('config.token_type'), 'batch');
    await click(GENERAL.submitButton);

    assert.strictEqual(
      page.flash.latestMessage,
      `The configuration was saved successfully.`,
      'success flash shows'
    );
  });

  // Runs against a real Vault so a key it does not recognise would show up as a value that is not kept.
  test('it saves lockout settings and shows them when the page is opened again', async function (assert) {
    assert.expect(7);

    const path = `userpass-lockout-${this.uid}`;
    await enablePage.enable('userpass', path);
    await page.visit({ path, section: 'options' });

    await fillIn(GENERAL.inputByAttr('user_lockout_config.lockout_threshold'), '7');
    await click(GENERAL.ttl.toggle('Lockout duration'));
    await fillIn(GENERAL.ttl.input('Lockout duration'), '10');
    await fillIn(TTL_PICKER.unitInputByLabel('Lockout duration'), 'm');
    await click(GENERAL.ttl.toggle('Lockout counter reset'));
    await fillIn(GENERAL.ttl.input('Lockout counter reset'), '30');
    await fillIn(TTL_PICKER.unitInputByLabel('Lockout counter reset'), 'm');
    await click(GENERAL.submitButton);
    assert.strictEqual(page.flash.latestMessage, 'The configuration was saved successfully.');

    await page.visit({ path, section: 'options' });
    assert.dom(GENERAL.inputByAttr('user_lockout_config.lockout_threshold')).hasValue('7');
    assert.dom(GENERAL.ttl.input('Lockout duration')).hasValue('10');
    assert.dom(TTL_PICKER.unitInputByLabel('Lockout duration')).hasValue('m');
    assert.dom(GENERAL.ttl.input('Lockout counter reset')).hasValue('30');
    assert.dom(TTL_PICKER.unitInputByLabel('Lockout counter reset')).hasValue('m');

    await click(GENERAL.submitButton);
    await page.visit({ path, section: 'options' });
    assert.dom(GENERAL.ttl.input('Lockout counter reset')).hasValue('30', 'saving again keeps the values');
  });

  for (const type of ['aws', 'azure', 'gcp', 'github', 'kubernetes']) {
    test(`it shows tabs for auth method: ${type}`, async function (assert) {
      const path = `${type}-showtab-${this.uid}`;
      await cli.toggle();
      await settled();
      await cli.consoleInput(`write sys/auth/${path} type=${type}`);
      await cli.enter();
      await indexPage.visit({ path });
      // aws has 4 tabs, the others will have 'Configuration' and 'Method Options' tabs
      const numTabs = type === 'aws' ? 4 : 2;
      const tabs = findAll(GENERAL.linkTo('auth-tab'));
      assert.strictEqual(tabs.length, numTabs, 'shows correct number of tabs');
    });
  }
});
