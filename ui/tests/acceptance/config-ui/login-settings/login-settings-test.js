/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, visit, currentRouteName } from '@ember/test-helpers';
import syncScenario from 'vault/mirage/scenarios/sync';
import syncHandlers from 'vault/mirage/handlers/sync';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | config-ui/login-settings', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    syncScenario(this.server);
    syncHandlers(this.server);
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const listPayload = {
      keys: ['ns1', 'ns2'],
      keyInfo: {
        ns1: {
          DisableInheritance: true,
          Name: 'ns1',
          Namespace: 'root',
        },
        ns2: {
          DisableInheritance: true,
          Name: 'ns2',
          Namespace: 'root',
        },
      },
    };

    const rulePayload = {
      data: {
        backup_auth_types: ['okta', 'userpass'],
        default_auth_type: 'ldap',
        disable_inheritance: true,
        namespace: 'root/',
      },
    };

    const api = this.owner.lookup('service:api');
    this.ruleListStub = sinon.stub(api.sys, 'uiLoginDefaultAuthList').resolves(listPayload);
    this.singleRuleStub = sinon.stub(api.sys, 'uiLoginDefaultAuthReadConfiguration').resolves(rulePayload);
    this.deleteStub = sinon.stub(api.sys, 'uiLoginDefaultAuthDeleteConfiguration').resolves({});

    return login();
  });

  test('fetched login rule list renders', async function (assert) {
    // Visit the login settings list index page
    await visit('vault/config-ui/login-settings');

    // verify fetched rules are rendered in list
    assert.dom('.linked-block-item').exists({ count: 2 });
  });

  test('delete rule from list view', async function (assert) {
    // Visit the login settings list index page
    await visit('vault/config-ui/login-settings');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete-rule'));

    assert.dom(GENERAL.confirmationModal).exists();
    await click(GENERAL.confirmButton);

    // verify success message from deletion
    assert.dom(GENERAL.latestFlashContent).includesText('Successfully deleted rule ns1.');
  });

  test('navigate to rule details page and renders rule data', async function (assert) {
    // visit individual rule page
    await visit('vault/config-ui/login-settings/');

    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('view-rule'));

    // verify that user is redirected to the rule details page
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.config-ui.login-settings.rule.details',
      'goes to rule details page'
    );

    // verify fetched rule data is rendered
    assert.dom(GENERAL.infoRowLabel('Name')).exists();
    assert.dom(GENERAL.infoRowValue('Name')).hasText('ns1');
    assert.dom(GENERAL.infoRowLabel('Namespace')).exists();
    assert.dom(GENERAL.infoRowValue('Namespace')).hasText('root/');
  });

  test('delete rule from rule details page', async function (assert) {
    // visit individual rule page
    await visit('vault/config-ui/login-settings/ns1');

    // click delete button & confirm from modal
    await click('[data-test-rule-delete]');

    assert.dom(GENERAL.confirmationModal).exists();

    await click(GENERAL.confirmButton);

    // verify that user is redirected to the list page after deletion
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.config-ui.login-settings.index',
      'goes back to login rule list page'
    );
  });
});
