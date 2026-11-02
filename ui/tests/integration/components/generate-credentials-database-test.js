/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { ResponseError } from '@hashicorp/vault-client-typescript';
import credentialsTemplate from 'vault/templates/vault/cluster/secrets/backend/credentials';

module('Integration | Component | GenerateCredentialsDatabase', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.roleName = 'my-role';
    this.backendPath = 'database';
  });

  test('dynamic role with password renders username and password rows and omits RSA Private Key row', async function (assert) {
    this.roleType = 'dynamic';
    this.model = {
      username: 'dynamic-user',
      password: 'dynamic-pass',
      rsaPrivateKey: null,
      lease_id: 'database/creds/my-role/abcd',
      lease_duration: 3600,
    };

    await render(hbs`
      <GenerateCredentialsDatabase
        @roleName={{this.roleName}}
        @backendPath={{this.backendPath}}
        @roleType={{this.roleType}}
        @model={{this.model}}
      />
    `);

    assert
      .dom(GENERAL.infoRowLabel('Username'))
      .hasText('Username', 'Username row renders for dynamic role with password');
    assert
      .dom(GENERAL.infoRowLabel('Password'))
      .hasText('Password', 'Password row renders when password credential is present');
    assert
      .dom(GENERAL.infoRowLabel('RSA Private Key'))
      .doesNotExist('RSA Private Key row should not render when credential type is password');
  });

  test('dynamic role with rsa_private_key renders username and RSA Private Key rows and omits Password row', async function (assert) {
    this.roleType = 'dynamic';
    this.model = {
      username: 'dynamic-user',
      password: null,
      rsaPrivateKey: '-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----',
      lease_id: 'database/creds/my-role/abcd',
      lease_duration: 3600,
    };

    await render(hbs`
      <GenerateCredentialsDatabase
        @roleName={{this.roleName}}
        @backendPath={{this.backendPath}}
        @roleType={{this.roleType}}
        @model={{this.model}}
      />
    `);

    assert
      .dom(GENERAL.infoRowLabel('Username'))
      .hasText('Username', 'Username row renders for dynamic role with RSA key');
    assert
      .dom(GENERAL.infoRowLabel('RSA Private Key'))
      .hasText('RSA Private Key', 'RSA Private Key row renders when rsa_private_key credential is present');
    assert
      .dom(GENERAL.infoRowLabel('Password'))
      .doesNotExist('Password row should not render when credential type is rsa_private_key');
  });

  test('static role with password renders password row and omits RSA Private Key row', async function (assert) {
    this.roleType = 'static';
    this.model = {
      username: 'static-user',
      password: 'static-pass',
      rsaPrivateKey: null,
      last_vault_rotation: '2025-01-01T00:00:00Z',
      rotation_period: 86400,
      ttl: 3600,
    };

    await render(hbs`
      <GenerateCredentialsDatabase
        @roleName={{this.roleName}}
        @backendPath={{this.backendPath}}
        @roleType={{this.roleType}}
        @model={{this.model}}
      />
    `);

    assert
      .dom(GENERAL.infoRowLabel('Password'))
      .hasText('Password', 'Password row renders for static role with password credential');
    assert
      .dom(GENERAL.infoRowLabel('Username'))
      .hasText('Username', 'Username row renders for static role with password credential');
    assert
      .dom(GENERAL.infoRowLabel('RSA Private Key'))
      .doesNotExist('RSA Private Key row should not render when credential type is password');
  });

  test('static role with rsa_private_key renders RSA Private Key row and omits Password row', async function (assert) {
    this.roleType = 'static';
    this.model = {
      username: 'static-user',
      password: null,
      rsaPrivateKey: '-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----',
      last_vault_rotation: '2025-01-01T00:00:00Z',
      rotation_period: 86400,
      ttl: 3600,
    };

    await render(hbs`
      <GenerateCredentialsDatabase
        @roleName={{this.roleName}}
        @backendPath={{this.backendPath}}
        @roleType={{this.roleType}}
        @model={{this.model}}
      />
    `);

    assert
      .dom(GENERAL.infoRowLabel('RSA Private Key'))
      .hasText(
        'RSA Private Key',
        'RSA Private Key row renders for static role with rsa_private_key credential'
      );
    assert
      .dom(GENERAL.infoRowLabel('Username'))
      .hasText('Username', 'Username row renders for static role with rsa_private_key credential');
    assert
      .dom(GENERAL.infoRowLabel('Password'))
      .doesNotExist('Password row should not render when credential type is rsa_private_key');
  });
});

module('Integration | database credentials route and template', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.route = this.owner.lookup('route:vault/cluster/secrets/backend/credentials');
    sinon.stub(this.route, 'modelFor').returns({ id: 'database', type: 'database' });
    const api = this.owner.lookup('service:api');
    this.staticRequest = sinon.stub(api.secrets, 'databaseReadStaticRoleCredentials');
    this.dynamicRequest = sinon.stub(api.secrets, 'databaseGenerateCredentials');
  });

  for (const roleType of ['static', 'dynamic']) {
    test(`Overview credentials display for a ${roleType} role without a roleType query parameter`, async function (assert) {
      const successful = roleType === 'static' ? this.staticRequest : this.dynamicRequest;
      const unsuccessful = roleType === 'static' ? this.dynamicRequest : this.staticRequest;
      successful.resolves({
        data: {
          username: 'test-user',
          password: 'test-password',
          last_vault_rotation: '2026-01-01T00:00:00Z',
          rotation_period: 86400,
          ttl: 3600,
        },
        lease_id: 'database/creds/my-role/test-lease',
        lease_duration: 3600,
      });
      unsuccessful.rejects(
        new ResponseError(new Response(JSON.stringify({ errors: ['role not found'] }), { status: 404 }))
      );

      this.set('model', await this.route.model({ secret: 'my-role', roleType: '' }));
      await render(credentialsTemplate);

      assert.strictEqual(
        this.model.dbCred.roleType,
        roleType,
        'The successful endpoint resolves the role type'
      );
      assert.true(
        this.staticRequest.calledOnceWithExactly('my-role', 'database'),
        'Static endpoint tried once'
      );
      assert.true(
        this.dynamicRequest.calledOnceWithExactly('my-role', 'database'),
        'Dynamic endpoint tried once'
      );
      assert.dom(GENERAL.infoRowLabel('Username')).hasText('Username', 'Username row is displayed');
      assert.dom(GENERAL.infoRowLabel('Password')).hasText('Password', 'Password row is displayed');
      if (roleType === 'dynamic') {
        assert.dom(GENERAL.infoRowValue('Lease ID')).hasText('database/creds/my-role/test-lease');
        assert.dom(GENERAL.infoRowValue('Lease Duration')).hasText('1 hour');
      } else {
        assert.dom(GENERAL.infoRowValue('Rotation Period')).hasText('1 day');
        assert.dom(GENERAL.infoRowValue('Time Remaining')).hasText('1 hour');
      }
    });

    test(`an explicit ${roleType} role only calls its matching endpoint and displays RSA credentials`, async function (assert) {
      const successful = roleType === 'static' ? this.staticRequest : this.dynamicRequest;
      const unused = roleType === 'static' ? this.dynamicRequest : this.staticRequest;
      successful.resolves({
        data: {
          username: 'test-user',
          rsa_private_key: 'test-private-key',
          last_vault_rotation: '2026-01-01T00:00:00Z',
        },
      });

      this.set('model', await this.route.model({ secret: 'my-role', roleType }));
      await render(credentialsTemplate);

      assert.true(successful.calledOnceWithExactly('my-role', 'database'), 'Matching endpoint called once');
      assert.true(unused.notCalled, 'The other endpoint is not called');
      assert.strictEqual(this.model.dbCred.rsaPrivateKey, 'test-private-key', 'RSA key is normalized');
      assert.dom(GENERAL.infoRowLabel('RSA Private Key')).hasText('RSA Private Key');
      assert.dom(GENERAL.infoRowLabel('Password')).doesNotExist('No password row for key-pair credentials');
    });
  }

  test('a successful static response takes precedence if both endpoints succeed', async function (assert) {
    this.staticRequest.resolves({ data: { username: 'static-user' } });
    this.dynamicRequest.resolves({ data: { username: 'dynamic-user' } });

    const credential = await this.route.getDatabaseCredential('database', 'my-role');

    assert.strictEqual(
      credential.roleType,
      'static',
      'Static success is preferred, matching the old adapter'
    );
    assert.strictEqual(credential.username, 'static-user', 'Static credentials are returned');
  });

  for (const errorType of ['static', 'dynamic']) {
    test(`when both requests fail, the higher-status ${errorType} error is displayed`, async function (assert) {
      const higher = errorType === 'static' ? this.staticRequest : this.dynamicRequest;
      const lower = errorType === 'static' ? this.dynamicRequest : this.staticRequest;
      higher.rejects(
        new ResponseError(new Response(JSON.stringify({ errors: ['database unavailable'] }), { status: 500 }))
      );
      lower.rejects(new ResponseError(new Response(JSON.stringify({ errors: [] }), { status: 404 })));

      this.set('model', await this.route.model({ secret: 'my-role', roleType: '' }));
      await render(credentialsTemplate);

      assert.strictEqual(this.model.dbCred.status, 500, 'The more informative error is selected');
      assert.dom().includesText('database unavailable', 'API error message is rendered');
      assert.dom(GENERAL.infoRowLabel('Password')).doesNotExist('Errors are not rendered as credentials');
    });

    test(`${errorType} control-group errors are propagated when neither request succeeds`, async function (assert) {
      const locked = errorType === 'static' ? this.staticRequest : this.dynamicRequest;
      const other = errorType === 'static' ? this.dynamicRequest : this.staticRequest;
      locked.rejects(
        new ResponseError(
          new Response(JSON.stringify({ isControlGroupError: true, message: 'Control Group encountered' }), {
            status: 403,
          })
        )
      );
      other.rejects(new ResponseError(new Response(JSON.stringify({ errors: [] }), { status: 404 })));

      await assert.rejects(
        this.route.getDatabaseCredential('database', 'my-role'),
        (error) => error.isControlGroupError === true,
        'Control-group authorization is not swallowed or replaced by the other endpoint error'
      );
    });
  }
});
