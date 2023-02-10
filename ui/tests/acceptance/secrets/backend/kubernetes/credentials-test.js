import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import kubernetesScenario from 'vault/mirage/scenarios/kubernetes';
import ENV from 'vault/config/environment';
import authPage from 'vault/tests/pages/auth';
import { typeIn, visit, click, currentRouteName } from '@ember/test-helpers';

module('Acceptance | kubernetes | credentials', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'kubernetes';
  });
  hooks.beforeEach(function () {
    kubernetesScenario(this.server);
    this.visitRoleCredentials = () => {
      return visit('/vault/secrets/kubernetes/kubernetes/roles/role-0/credentials');
    };
    this.validateRoute = (assert, route, message) => {
      assert.strictEqual(currentRouteName(), `vault.cluster.secrets.backend.kubernetes.${route}`, message);
    };
    return authPage.login();
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  test('it should have correct breadcrumb links in credentials view', async function (assert) {
    assert.expect(3);
    await this.visitRoleCredentials();
    await click('[data-test-breadcrumbs] li:nth-child(3) a');
    this.validateRoute(assert, 'roles.role.details', 'Transitions to role details route on breadcrumb click');
    await this.visitRoleCredentials();
    await click('[data-test-breadcrumbs] li:nth-child(2) a');
    this.validateRoute(assert, 'roles.index', 'Transitions to roles route on breadcrumb click');
    await this.visitRoleCredentials();
    await click('[data-test-breadcrumbs] li:nth-child(1) a');
    this.validateRoute(assert, 'overview', 'Transitions to overview route on breadcrumb click');
  });

  test('it should transition to role details view on Back click', async function (assert) {
    assert.expect(1);
    await this.visitRoleCredentials();
    await click('[data-test-generate-credentials-back]');

    await this.validateRoute(assert, 'roles.role.details', 'Transitions to role details on Back click');
  });

  test('it should transition to role details view on Done click', async function (assert) {
    assert.expect(1);
    await this.visitRoleCredentials();
    this.server.post('/kubernetes-test/creds/role-0', () => {
      assert.ok('POST request made to generate credentials');
      return {
        request_id: '58fefc6c-5195-c17a-94f2-8f889f3df57c',
        lease_id: 'kubernetes/creds/default-role/aWczfcfJ7NKUdiirJrPXIs38',
        renewable: false,
        lease_duration: 3600,
        data: {
          service_account_name: 'default',
          service_account_namespace: 'default',
          service_account_token: 'eyJhbGciOiJSUzI1NiIsImtpZCI6Imlr',
        },
      };
    });
    await typeIn('[data-test-kubernetes-namespace]', 'kubernetes-test');
    await click('[data-test-toggle-input]');
    await click('[data-test-toggle-input="Time-to-Live (TTL)"]');
    await typeIn('[data-test-ttl-value="Time-to-Live (TTL)"]', 2);
    await click('[data-test-generate-credentials-button]');
    await click('[data-test-generate-credentials-done]');

    await this.validateRoute(assert, 'roles.role.details', 'Transitions to role details on Done click');
  });
});
