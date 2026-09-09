/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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
