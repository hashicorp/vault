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
    this.generateResponse = ({ path, capabilities }) => {
      if (path) {
        return {
          request_id: '6cc7a484-921a-a730-179c-eaf6c6fbe97e',
          data: {
            capabilities: capabilities,
            [path]: capabilities,
          },
        };
      }
    };
  });

  test('it makes request to capabilities-self', function (assert) {
    const path = '/my/api/path';
    const expectedPayload = {
      paths: [path],
    };
    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse({ path, capabilities: ['read'] });
    });
    this.capabilities.request({ path });
  });

  const SINGULAR_PATH = [
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
  SINGULAR_PATH.forEach(({ capabilities, canRead, canUpdate }) => {
    const path = '/my/api/path';
    test(`singular path returns expected boolean for "${capabilities.join(', ')}"`, async function (assert) {
      this.server.post('/sys/capabilities-self', () => {
        return this.generateResponse({ path, capabilities });
      });

      const canReadResponse = await this.capabilities.canRead(path);
      const canUpdateResponse = await this.capabilities.canUpdate(path);
      assert[canRead](canReadResponse, `canRead returns ${canRead}`);
      assert[canUpdate](canUpdateResponse, `canUpdate returns ${canRead}`);
    });
  });
});
