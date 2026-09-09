/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { stub } from 'sinon';

module('Unit | Route | vault/cluster/settings/auth/configure', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.api = this.owner.lookup('service:api');
    this.route = this.owner.lookup('route:vault/cluster/settings/auth/configure');
    stub(this.api.sys, 'authListEnabledMethods').resolves({
      data: {
        'jwt/': {
          type: 'jwt',
          accessor: 'auth_jwt_abc123',
          config: {
            listing_visibility: 'hidden',
            default_lease_ttl: 0,
            max_lease_ttl: 0,
            token_type: 'default-service',
            force_no_cache: false,
          },
          description: '',
          external_entropy_access: false,
          local: false,
          options: {},
          plugin_version: '',
          running_plugin_version: 'v0.21.0+builtin',
          running_sha256: '',
          seal_wrap: false,
          uuid: 'test-uuid-1234',
        },
      },
    });
  });

  test('it loads `methodOptions` from `authListEnabledMethods`', async function (assert) {
    const result = await this.route.model({ method: 'jwt' });

    assert.true(this.api.sys.authListEnabledMethods.calledOnce, 'authListEnabledMethods was called');
    assert.deepEqual(
      {
        path: result.methodOptions.path,
        type: result.methodOptions.type,
        listing_visibility: result.methodOptions.config.listing_visibility,
        id: result.method.id,
      },
      {
        path: 'jwt/',
        type: 'jwt',
        listing_visibility: 'hidden',
        id: 'jwt',
      },
      'model returns correct methodOptions shape and method.id strips trailing slash'
    );
  });

  test('it throws a 404 when the path does not match any method', async function (assert) {
    await assert.rejects(
      this.route.model({ method: 'nonexistent' }),
      (err) => err.httpStatus === 404 && err.path === 'nonexistent'
    );
  });
});
