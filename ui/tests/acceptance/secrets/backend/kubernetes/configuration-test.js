import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import kubernetesScenario from 'vault/mirage/scenarios/kubernetes';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { visit, click, currentRouteName } from '@ember/test-helpers';
import { Response } from 'miragejs';

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
      'vault.cluster.secrets.backend.kubernetes.configuration',
      'Transitions to configuration route on fetch 404'
    );
    assert.dom('[data-test-empty-state-title]').hasText('Kubernetes not configured', 'Config cta renders');
  });
});
