/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import kubernetesScenario from 'vault/mirage/scenarios/kubernetes';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { visit, click, currentRouteName } from '@ember/test-helpers';

module('Acceptance | kubernetes | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'kubernetes';
  });
  hooks.beforeEach(function () {
    kubernetesScenario(this.server);
    this.visitConfiguration = () => {
      return visit('/vault/secrets/kubernetes/kubernetes/configuration');
    };
    this.validateRoute = (assert, route, message) => {
      assert.strictEqual(currentRouteName(), `vault.cluster.secrets.backend.kubernetes.${route}`, message);
    };
    return authPage.login();
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should transition to configure page on Edit Configuration click from toolbar', async function (assert) {
    assert.expect(1);
    await this.visitConfiguration();
    await click('[data-test-toolbar-config-action]');
    this.validateRoute(assert, 'configure', 'Transitions to Configure route on click');
  });
  test('it should transition to the configuration page on Save click in Configure', async function (assert) {
    assert.expect(1);
    await this.visitConfiguration();
    await click('[data-test-toolbar-config-action]');
    await click('[data-test-config-save]');
    await click('[data-test-config-confirm]');
    this.validateRoute(assert, 'configuration', 'Transitions to Configuration route on click');
  });
  test('it should transition to the configuration page on Cancel click in Configure', async function (assert) {
    assert.expect(1);
    await this.visitConfiguration();
    await click('[data-test-toolbar-config-action]');
    await click('[data-test-config-cancel]');
    this.validateRoute(assert, 'configuration', 'Transitions to Configuration route on click');
  });
});
