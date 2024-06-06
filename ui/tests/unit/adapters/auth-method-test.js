/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | auth method', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.mockResponse = {
      data: {
        auth: {
          'approle/': {
            accessor: 'auth_approle_43e5a627',
            config: {
              default_lease_ttl: 2764800,
              force_no_cache: false,
              listing_visibility: 'hidden',
              max_lease_ttl: 2764800,
              token_type: 'default-service',
            },
            uuid: '7a8bc146-76d0-3a9c-9feb-47a6713a85b3',
          },
        },
      },
    };
  });

  test('findAll makes request to correct endpoint with no adapterOptions', async function (assert) {
    assert.expect(1);

    this.server.get('sys/auth', () => {
      assert.ok(true, 'request made to sys/auth when no options are passed to findAll');
      return { data: this.mockResponse.data.auth };
    });

    await this.store.findAll('auth-method');
  });

  test('findAll makes request to correct endpoint when unauthenticated is true', async function (assert) {
    assert.expect(1);

    this.server.get('sys/internal/ui/mounts', () => {
      assert.ok(true, 'request made to correct endpoint when unauthenticated');
      return this.mockResponse;
    });

    await this.store.findAll('auth-method', { adapterOptions: { unauthenticated: true } });
  });

  test('findAll makes request to correct endpoint when useMountsEndpoint is true', async function (assert) {
    assert.expect(1);

    this.server.get('sys/internal/ui/mounts', () => {
      assert.ok(true, 'request made to correct endpoint when useMountsEndpoint');
      return this.mockResponse;
    });

    await this.store.findAll('auth-method', { adapterOptions: { useMountsEndpoint: true } });
  });
});
