/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import kubernetesScenario from 'vault/mirage/scenarios/kubernetes';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { fillIn, visit, currentURL, click, currentRouteName } from '@ember/test-helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Acceptance | kubernetes | roles', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'kubernetes';
  });
  hooks.beforeEach(function () {
    kubernetesScenario(this.server);
    this.visitRoles = () => {
      return visit('/vault/secrets/kubernetes/kubernetes/roles');
    };
    this.validateRoute = (assert, route, message) => {
      assert.strictEqual(currentRouteName(), `vault.cluster.secrets.backend.kubernetes.${route}`, message);
    };
    return authPage.login();
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should filter roles', async function (assert) {
    await this.visitRoles();
    assert.dom('[data-test-list-item-link]').exists({ count: 3 }, 'Roles list renders');
    await fillIn('[data-test-component="navigate-input"]', '1');
    assert.dom('[data-test-list-item-link]').exists({ count: 1 }, 'Filtered roles list renders');
    assert.ok(currentURL().includes('pageFilter=1'), 'pageFilter query param value is set');
  });

  test('it should link to role details on list item click', async function (assert) {
    assert.expect(1);
    await this.visitRoles();
    await click('[data-test-list-item-link]');
    this.validateRoute(assert, 'roles.role.details', 'Transitions to details route on list item click');
  });

  test('it should have correct breadcrumb links in role details view', async function (assert) {
    assert.expect(2);
    await this.visitRoles();
    await click('[data-test-list-item-link]');
    await click('[data-test-breadcrumbs] li:nth-child(2) a');
    this.validateRoute(assert, 'roles.index', 'Transitions to roles route on breadcrumb click');
    await click('[data-test-list-item-link]');
    await click('[data-test-breadcrumbs] li:nth-child(1) a');
    this.validateRoute(assert, 'overview', 'Transitions to overview route on breadcrumb click');
  });

  test('it should have functional list item menu', async function (assert) {
    // Popup menu causes flakiness
    setRunOptions({
      rules: {
        'color-contrast': { enabled: false },
      },
    });
    assert.expect(3);
    await this.visitRoles();
    for (const action of ['details', 'edit', 'delete']) {
      await click('[data-test-list-item-popup] button');
      await click(`[data-test-${action}]`);
      if (action === 'delete') {
        await click('[data-test-confirm-button]');
        assert.dom('[data-test-list-item-link]').exists({ count: 2 }, 'Deleted role removed from list');
      } else {
        this.validateRoute(
          assert,
          `roles.role.${action}`,
          `Transitions to ${action} route on menu action click`
        );
        const selector =
          action === 'details' ? '[data-test-breadcrumbs] li:nth-child(2) a' : '[data-test-cancel]';
        await click(selector);
      }
    }
  });

  test('it should create role', async function (assert) {
    assert.expect(2);
    await this.visitRoles();
    await click('[data-test-toolbar-roles-action]');
    await click('[data-test-radio-card="basic"]');
    await fillIn('[data-test-input="name"]', 'new-test-role');
    await fillIn('[data-test-input="serviceAccountName"]', 'default');
    await fillIn('[data-test-input="allowedKubernetesNamespaces"]', '*');
    await click('[data-test-save]');
    this.validateRoute(assert, 'roles.role.details', 'Transitions to details route on save success');
    await click('[data-test-breadcrumbs] li:nth-child(2) a');
    assert.dom('[data-test-role="new-test-role"]').exists('New role renders in list');
  });

  test('it should have functional toolbar actions in details view', async function (assert) {
    assert.expect(3);
    await this.visitRoles();
    await click('[data-test-list-item-link]');
    await click('[data-test-generate-credentials]');
    this.validateRoute(assert, 'roles.role.credentials', 'Transitions to credentials route');
    await click('[data-test-breadcrumbs] li:nth-child(3) a');
    await click('[data-test-edit]');
    this.validateRoute(assert, 'roles.role.edit', 'Transitions to edit route');
    await click('[data-test-cancel]');
    await click('[data-test-list-item-link]');
    await click('[data-test-delete]');
    await click('[data-test-confirm-button]');
    assert
      .dom('[data-test-list-item-link]')
      .exists({ count: 2 }, 'Transitions to roles route and deleted role removed from list');
  });

  test('it should generate credentials for role', async function (assert) {
    assert.expect(1);
    await this.visitRoles();
    await click('[data-test-list-item-link]');
    await click('[data-test-generate-credentials]');
    await fillIn('[data-test-kubernetes-namespace]', 'test-namespace');
    await click('[data-test-generate-credentials-button]');
    await click('[data-test-generate-credentials-done]');
    this.validateRoute(
      assert,
      'roles.role.details',
      'Transitions to details route when done generating credentials'
    );
  });
});
