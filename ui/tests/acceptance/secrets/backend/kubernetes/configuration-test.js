/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import kubernetesScenario from 'vault/mirage/scenarios/kubernetes';
import kubernetesHandlers from 'vault/mirage/handlers/kubernetes';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { visit, click, currentRouteName } from '@ember/test-helpers';
import { Response } from 'miragejs';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

module('Acceptance | kubernetes | configuration', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    kubernetesHandlers(this.server);
    kubernetesScenario(this.server);
    this.visitConfiguration = () => {
      return visit('/vault/secrets-engines/kubernetes/kubernetes/configuration');
    };
    this.validateRoute = (assert, route, message) => {
      assert.strictEqual(currentRouteName(), `vault.cluster.secrets.backend.kubernetes.${route}`, message);
    };
    return login();
  });

  test('it should transition to configure page on Edit Configuration click from toolbar', async function (assert) {
    assert.expect(1);
    await this.visitConfiguration();
    await click(SES.configure);
    this.validateRoute(assert, 'configure', 'Transitions to Configure route on click');
  });

  test('it should transition to the configuration page on Save click in Configure', async function (assert) {
    assert.expect(1);
    await this.visitConfiguration();
    await click(SES.configure);
    await click('[data-test-config-save]');
    await click('[data-test-config-confirm]');
    this.validateRoute(assert, 'configuration', 'Transitions to Configuration route on click');
  });

  test('it should transition to the configuration page on Cancel click in Configure', async function (assert) {
    assert.expect(1);
    await this.visitConfiguration();
    await click(SES.configure);
    await click('[data-test-config-cancel]');
    this.validateRoute(assert, 'configuration', 'Transitions to Configuration route on click');
  });

  test('it should transition to error route on config fetch error other than 404', async function (assert) {
    this.server.get('/kubernetes/config', () => new Response(403));
    await this.visitConfiguration();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.kubernetes.error',
      'Transitions to error route on config fetch error'
    );
  });

  test('it should not transition to error route on config fetch 404', async function (assert) {
    this.server.get('/kubernetes/config', () => new Response(404));
    await this.visitConfiguration();
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.kubernetes.configure',
      'Transitions to configure route on fetch 404'
    );
  });
});
