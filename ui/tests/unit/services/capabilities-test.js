/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Service | capabilities', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.capabilities = this.owner.lookup('service:capabilities');
    this.store = this.owner.lookup('service:store');
    this.generateResponse = (apiPath, perms) => {
      return {
        [apiPath]: perms,
        capabilities: perms,
        request_id: '6cc7a484-921a-a730-179c-eaf6c6fbe97e',
        data: {
          capabilities: perms,
          [apiPath]: perms,
        },
      };
    };
  });

  test('it makes request to capabilities-self', function (assert) {
    const apiPath = '/my/api/path';
    const expectedPayload = {
      paths: [apiPath],
    };
    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse(apiPath, ['read']);
    });
    this.capabilities.request(apiPath);
  });

  const TEST_CASES = [
    {
      capabilities: ['read'],
      canRead: true,
      canUpdate: false,
    },
    {
      capabilities: ['update'],
      canRead: false,
      canUpdate: true,
    },
    {
      capabilities: ['deny'],
      canRead: false,
      canUpdate: false,
    },
    {
      capabilities: ['read', 'update'],
      canRead: true,
      canUpdate: true,
    },
  ];
  TEST_CASES.forEach(({ capabilities, canRead, canUpdate }) => {
    test(`it returns expected boolean for "${capabilities.join(', ')}"`, async function (assert) {
      const apiPath = '/my/api/path';
      this.server.post('/sys/capabilities-self', () => {
        return this.generateResponse(apiPath, capabilities);
      });

      const canReadResponse = await this.capabilities.canRead(apiPath);
      const canUpdateResponse = await this.capabilities.canUpdate(apiPath);
      assert[canRead](canReadResponse, `canRead returns ${canRead}`);
      assert[canUpdate](canUpdateResponse, `canUpdate returns ${canRead}`);
    });
  });
});
