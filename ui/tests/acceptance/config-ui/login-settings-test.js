/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, currentRouteName, fillIn, waitUntil } from '@ember/test-helpers';
import { login, logout, rootToken } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { overrideResponse } from 'vault/tests/helpers/stubs';

import customLoginScenario from 'vault/mirage/scenarios/custom-login';
import customLoginHandler from 'vault/mirage/handlers/custom-login';
import Sinon from 'sinon';

// read view for custom login settings
module('Acceptance | Enterprise | config-ui/login-settings', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
  });

  test('it renders empty state if no login settings exist', async function (assert) {
    await visit('vault/config-ui/login-settings');
    assert.dom(GENERAL.emptyStateTitle).hasText('No UI login settings yet');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'Login settings can be used to customize which methods display in the web UI login form by setting a default and back up login methods. Available to be created via the CLI or HTTP API.'
      );
  });

  test('it renders error template when permission is denied', async function (assert) {
    this.server.get('/sys/config/ui/login/default-auth', () => overrideResponse(403));
    await visit('vault/config-ui/login-settings');
    assert.dom(GENERAL.pageError.error).hasText('Error permission denied');
  });

  module('list, read and delete', function (hooks) {
    hooks.beforeEach(async function () {
      customLoginScenario(this.server);
      customLoginHandler(this.server);
      this.loginRules = this.server.db.loginRules;

      // Cannot use the login() helper because customLoginHandler returns "token" as the default auth method
      await logout();
      await fillIn(GENERAL.inputByAttr('token'), rootToken);
      await click(GENERAL.submitButton);
      await waitUntil(() => currentRouteName() === 'vault.cluster.dashboard');
      await click(GENERAL.navLink('Operational tools'));
      await click(GENERAL.navLink('UI login settings'));
    });

    test('it renders login rules', async function (assert) {
      assert
        .dom(GENERAL.listItem())
        .exists({ count: this.loginRules.length }, `${this.loginRules.length} rules render`);
      this.loginRules.forEach(({ name, disable_inheritance, namespace_path }) => {
        const inheritance = disable_inheritance ? 'Inheritance disabled' : 'Inheritance enabled';
        assert.dom(GENERAL.listItem(name)).hasText(`${name} ${namespace_path} ${inheritance}`);
      });
    });

    test('it deletes rule from list view', async function (assert) {
      const successFlashSpy = Sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
      const ruleToDelete = this.loginRules[0].name;
      const initialCount = this.loginRules.length; // cache record length so we can confirm delete
      await click(`${GENERAL.listItem(ruleToDelete)} ${GENERAL.menuTrigger}`);
      await click(GENERAL.menuItem('delete-rule'));
      await click(GENERAL.confirmButton);
      const [success] = successFlashSpy.lastCall.args;
      assert.strictEqual(
        success,
        `Successfully deleted rule ${ruleToDelete}.`,
        'it calls flash success with expected message'
      );
      assert.dom(GENERAL.listItem(ruleToDelete)).doesNotExist('the deleted rule does not exist');
      assert.dom(GENERAL.listItem()).exists({ count: initialCount - 1 }, `${initialCount - 1} rules render`);
    });

    test('it navigates to rule details page and renders rule data', async function (assert) {
      const rule = this.server.db.loginRules[0];
      // visit individual rule page
      await click(`${GENERAL.listItem(rule.name)} ${GENERAL.menuTrigger}`);
      await click(GENERAL.menuItem('view-rule'));
      // verify that user is redirected to the rule details page
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.login-settings.rule.details',
        'goes to rule details page'
      );
      // verify fetched rule data is rendered
      assert.dom(GENERAL.infoRowValue('Default method')).hasText(rule.default_auth_type);
      assert.dom(GENERAL.infoRowValue('Namespace the rule applies to')).hasText(rule.namespace_path);
      assert.dom(GENERAL.infoRowValue('Backup methods')).hasText(rule.backup_auth_types.join(','));
      assert.dom(GENERAL.infoRowValue('Inheritance enabled')).hasText('Yes');
    });

    test('it navigates to rule details from linked block', async function (assert) {
      const rule = this.server.db.loginRules[2];
      await visit('vault/config-ui/login-settings');
      await click(GENERAL.listItem(rule.name));
      assert.strictEqual(
        currentRouteName(),
        'vault.cluster.config-ui.login-settings.rule.details',
        'goes to rule details page'
      );
      assert.dom(GENERAL.infoRowValue('Default method')).hasText(rule.default_auth_type);
      assert.dom(GENERAL.infoRowValue('Namespace the rule applies to')).hasText(rule.namespace_path);
      assert.dom(GENERAL.infoRowValue('Backup methods')).hasText(rule.backup_auth_types.join(', '));
      assert.dom(GENERAL.infoRowValue('Inheritance enabled')).hasText('No');
    });
  });
});
