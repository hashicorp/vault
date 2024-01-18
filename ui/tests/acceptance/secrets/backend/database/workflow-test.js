/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { Response } from 'miragejs';
import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { create } from 'ember-cli-page-object';

import ENV from 'vault/config/environment';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import flashMessage from 'vault/tests/pages/components/flash-message';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';

const flash = create(flashMessage);

const PAGE = {
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateAction: '[data-test-secret-create="connections"]',
  rotateModal: '[data-test-db-connection-modal-title]',
  confirmRotate: '[data-test-enable-rotate-connection]',
  skipRotate: '[data-test-enable-connection]',
  infoRow: '[data-test-component="info-table-row"]',
  infoRowLabel: (label) => `[data-test-row-label="${label}"]`,
  infoRowValue: (label) => `[data-test-row-value="${label}"]`,
};

const FORM = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  saveBtn: '[data-test-secret-save]',
};

async function fillOutConnection(name) {
  await fillIn(FORM.inputByAttr('name'), name);
  await fillIn(FORM.inputByAttr('plugin_name'), 'mysql-database-plugin');
  await fillIn(FORM.inputByAttr('connection_url'), '{{username}}:{{password}}@tcp(127.0.0.1:33060)/');
  await fillIn(FORM.inputByAttr('username'), 'admin');
  await fillIn(FORM.inputByAttr('password'), 'very-secure');
}

/**
 * This test set is for testing the flow for database secrets engine.
 */
module('Acceptance | database workflow', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'database';
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  hooks.beforeEach(async function () {
    this.backend = `db-workflow-${uuidv4()}`;
    this.store = this.owner.lookup('service:store');
    await authPage.login();
    await runCmd(mountEngineCmd('database', this.backend), false);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('connections', function (hooks) {
    hooks.beforeEach(function () {
      this.expectedRows = [
        { label: 'Database plugin', value: 'mysql-database-plugin' },
        { label: 'Connection name', value: `connect-${this.backend}` },
        { label: 'Use custom password policy', value: 'Default' },
        { label: 'Connection URL', value: '{{username}}:{{password}}@tcp(127.0.0.1:33060)/' },
        { label: 'Max open connections', value: '4' },
        { label: 'Max idle connections', value: '0' },
        { label: 'Max connection lifetime', value: '0s' },
        { label: 'Username template', value: 'Default' },
        {
          label: 'Root rotation statements',
          value: `SELECT user from mysql.user,GRANT ALL PRIVILEGES ON *.* to 'sudo'@'%'`,
        },
      ];
    });
    test('create with rotate', async function (assert) {
      assert.expect(24);
      this.server.post('/:backend/rotate-root/:name', () => {
        assert.ok(true, 'rotate root called');
        new Response(204);
      });
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(PAGE.emptyStateAction);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.backend}/create`, 'Takes you to create page');

      // fill in connection details
      await fillOutConnection(`connect-${this.backend}`);
      await click(FORM.saveBtn);

      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.confirmRotate);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/connect-${this.backend}`,
        'Takes you to details page for connection'
      );
      assert.dom(PAGE.infoRow).exists({ count: this.expectedRows.length }, 'correct number of rows');
      this.expectedRows.forEach(({ label, value }) => {
        assert.dom(PAGE.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(PAGE.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
    });
    test('create without rotate', async function (assert) {
      assert.expect(23);
      this.server.post('/:backend/rotate-root/:name', () => {
        assert.notOk(true, 'rotate root called when it should not have been');
        new Response(204);
      });
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(PAGE.emptyStateAction);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.backend}/create`, 'Takes you to create page');

      // fill in connection details
      await fillOutConnection(`connect-${this.backend}`);
      await click(FORM.saveBtn);

      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.skipRotate);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/connect-${this.backend}`,
        'Takes you to details page for connection'
      );
      assert.dom(PAGE.infoRow).exists({ count: this.expectedRows.length }, 'correct number of rows');
      this.expectedRows.forEach(({ label, value }) => {
        assert.dom(PAGE.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(PAGE.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
    });
    test('create failure', async function (assert) {
      assert.expect(25);
      this.server.post('/:backend/rotate-root/:name', (schema, req) => {
        const okay = req.params.name !== 'bad-connection';
        assert.ok(okay, 'rotate root called but not for bad-connection');
        new Response(204);
      });
      await visit(`/vault/secrets/${this.backend}/overview`);
      assert.dom(PAGE.emptyStateTitle).hasText('Connect a database', 'empty state title is correct');
      await click(PAGE.emptyStateAction);
      assert.strictEqual(currentURL(), `/vault/secrets/${this.backend}/create`, 'Takes you to create page');

      // fill in connection details
      await fillOutConnection(`bad-connection`);
      await click(FORM.saveBtn);
      assert.strictEqual(
        flash.latestMessage,
        `error creating database object: error verifying - ping: Error 1045 (28000): Access denied for user 'admin'@'192.168.65.1' (using password: YES)`,
        'shows the error message from API'
      );
      await fillIn(FORM.inputByAttr('name'), `connect-${this.backend}`);
      await click(FORM.saveBtn);
      assert.dom(PAGE.rotateModal).hasText('Rotate your root credentials?', 'rotate modal is shown');
      await click(PAGE.confirmRotate);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${this.backend}/show/connect-${this.backend}`,
        'Takes you to details page for connection'
      );
      assert.dom(PAGE.infoRow).exists({ count: this.expectedRows.length }, 'correct number of rows');
      this.expectedRows.forEach(({ label, value }) => {
        assert.dom(PAGE.infoRowLabel(label)).hasText(label, `Label for ${label} is correct`);
        assert.dom(PAGE.infoRowValue(label)).hasText(value, `Value for ${label} is correct`);
      });
    });
  });
});
