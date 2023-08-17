import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { Response } from 'miragejs';
import { ERROR_JWT_LOGIN } from 'vault/components/auth-jwt';

module('Acceptance | jwt auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.stub = sinon.stub();
    this.server.post(
      '/auth/:path/oidc/auth_url',
      () =>
        new Response(
          400,
          { 'Content-Type': 'application/json' },
          JSON.stringify({ errors: [ERROR_JWT_LOGIN] })
        )
    );
    this.server.get('/auth/foo/oidc/callback', () => ({
      auth: { client_token: 'root' },
    }));
  });

  test('it works correctly with default name and no role', async function (assert) {
    assert.expect(6);
    this.server.post('/auth/jwt/login', (schema, req) => {
      const { jwt, role } = JSON.parse(req.requestBody);
      assert.ok(true, 'request made to auth/jwt/login after submit');
      assert.equal(jwt, 'my-test-jwt-token', 'JWT token is sent in body');
      assert.equal(role, undefined, 'role is not sent in body when not filled in');
      req.passthrough();
    });
    await visit('/vault/auth');
    await fillIn('[data-test-select="auth-method"]', 'jwt');
    assert.dom('[data-test-role]').exists({ count: 1 }, 'Role input exists');
    assert.dom('[data-test-jwt]').exists({ count: 1 }, 'JWT input exists');
    await fillIn('[data-test-jwt]', 'my-test-jwt-token');
    await click('[data-test-auth-submit]');
    assert.dom('[data-test-error]').exists('Failed login');
  });

  test('it works correctly with default name and a role', async function (assert) {
    assert.expect(7);
    this.server.post('/auth/jwt/login', (schema, req) => {
      const { jwt, role } = JSON.parse(req.requestBody);
      assert.ok(true, 'request made to auth/jwt/login after login');
      assert.equal(jwt, 'my-test-jwt-token', 'JWT token is sent in body');
      assert.equal(role, 'some-role', 'role is sent in the body when filled in');
      req.passthrough();
    });
    await visit('/vault/auth');
    await fillIn('[data-test-select="auth-method"]', 'jwt');
    assert.dom('[data-test-role]').exists({ count: 1 }, 'Role input exists');
    assert.dom('[data-test-jwt]').exists({ count: 1 }, 'JWT input exists');
    await fillIn('[data-test-role]', 'some-role');
    await fillIn('[data-test-jwt]', 'my-test-jwt-token');
    assert.dom('[data-test-jwt]').exists({ count: 1 }, 'JWT input exists');
    await click('[data-test-auth-submit]');
    assert.dom('[data-test-error]').exists('Failed login');
  });

  test('it works correctly with custom endpoint and a role', async function (assert) {
    assert.expect(6);
    this.server.get('/sys/internal/ui/mounts', () => ({
      data: {
        auth: {
          'test-jwt/': { description: '', options: {}, type: 'jwt' },
        },
      },
    }));
    this.server.post('/auth/test-jwt/login', (schema, req) => {
      const { jwt, role } = JSON.parse(req.requestBody);
      assert.ok(true, 'request made to auth/custom-jwt-login after login');
      assert.equal(jwt, 'my-test-jwt-token', 'JWT token is sent in body');
      assert.equal(role, 'some-role', 'role is sent in body when filled in');
      req.passthrough();
    });
    await visit('/vault/auth');
    await click('[data-test-auth-method-link="jwt"]');
    assert.dom('[data-test-role]').exists({ count: 1 }, 'Role input exists');
    assert.dom('[data-test-jwt]').exists({ count: 1 }, 'JWT input exists');
    await fillIn('[data-test-role]', 'some-role');
    await fillIn('[data-test-jwt]', 'my-test-jwt-token');
    await click('[data-test-auth-submit]');
    assert.dom('[data-test-error]').exists('Failed login');
  });
});
