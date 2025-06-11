/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, visit, fillIn, waitFor } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { ERROR_JWT_LOGIN } from 'vault/components/auth/form/oidc-jwt';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | jwt auth method', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    localStorage.clear(); // ensure that a token isn't stored otherwise visit('/vault/auth') will redirect to secrets
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
      assert.true(true, 'request made to auth/jwt/login after submit');
      assert.strictEqual(jwt, 'my-test-jwt-token', 'JWT token is sent in body');
      assert.strictEqual(role, undefined, 'role is not sent in body when not filled in');
      return overrideResponse(403);
    });
    await visit('/vault/auth');
    await fillIn(AUTH_FORM.selectMethod, 'jwt');
    assert.dom(GENERAL.inputByAttr('role')).exists({ count: 1 }, 'Role input exists');
    assert.dom(GENERAL.inputByAttr('jwt')).exists({ count: 1 }, 'JWT input exists');
    await fillIn(GENERAL.inputByAttr('jwt'), 'my-test-jwt-token');
    await click(GENERAL.submitButton);
    await waitFor(GENERAL.messageError);
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: permission denied');
  });

  test('it works correctly with default name and a role', async function (assert) {
    assert.expect(7);
    this.server.post('/auth/jwt/login', (schema, req) => {
      const { jwt, role } = JSON.parse(req.requestBody);
      assert.ok(true, 'request made to auth/jwt/login after login');
      assert.strictEqual(jwt, 'my-test-jwt-token', 'JWT token is sent in body');
      assert.strictEqual(role, 'some-role', 'role is sent in the body when filled in');
      return overrideResponse(403);
    });
    await visit('/vault/auth');
    await fillIn(AUTH_FORM.selectMethod, 'jwt');
    assert.dom(GENERAL.inputByAttr('role')).exists({ count: 1 }, 'Role input exists');
    assert.dom(GENERAL.inputByAttr('jwt')).exists({ count: 1 }, 'JWT input exists');
    await fillIn(GENERAL.inputByAttr('role'), 'some-role');
    await fillIn(GENERAL.inputByAttr('jwt'), 'my-test-jwt-token');
    assert.dom(GENERAL.inputByAttr('jwt')).exists({ count: 1 }, 'JWT input exists');
    await click(GENERAL.submitButton);
    await waitFor(GENERAL.messageError);
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: permission denied');
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
      assert.strictEqual(jwt, 'my-test-jwt-token', 'JWT token is sent in body');
      assert.strictEqual(role, 'some-role', 'role is sent in body when filled in');
      return overrideResponse(403);
    });
    await visit('/vault/auth');
    await click(AUTH_FORM.tabBtn('jwt'));
    assert.dom(GENERAL.inputByAttr('role')).exists({ count: 1 }, 'Role input exists');
    assert.dom(GENERAL.inputByAttr('jwt')).exists({ count: 1 }, 'JWT input exists');
    await fillIn(GENERAL.inputByAttr('role'), 'some-role');
    await fillIn(GENERAL.inputByAttr('jwt'), 'my-test-jwt-token');
    await click(GENERAL.submitButton);
    await waitFor(GENERAL.messageError);
    assert.dom(GENERAL.messageError).hasText('Error Authentication failed: permission denied');
  });
});
