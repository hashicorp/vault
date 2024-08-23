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
    this.generateResponse = ({ path, paths, capabilities }) => {
      if (path) {
        // "capabilities" is an array
        return {
          request_id: '6cc7a484-921a-a730-179c-eaf6c6fbe97e',
          data: {
            capabilities,
            [path]: capabilities,
          },
        };
      }
      if (paths) {
        // "capabilities" is an object, paths are keys and values are array of capabilities
        const data = paths.reduce((obj, path) => {
          obj[path] = capabilities[path];
          return obj;
        }, {});
        return {
          request_id: '6cc7a484-921a-a730-179c-eaf6c6fbe97e',
          data,
        };
      }
    };
  });

  module('general methods', function () {
    test('request: it makes request to capabilities-self with path param', function (assert) {
      const path = '/my/api/path';
      const expectedPayload = { paths: [path] };
      this.server.post('/sys/capabilities-self', (schema, req) => {
        const actual = JSON.parse(req.requestBody);
        assert.true(true, 'request made to capabilities-self');
        assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
        return this.generateResponse({ path, capabilities: ['read'] });
      });
      this.capabilities.request({ path });
    });

    test('request: it makes request to capabilities-self with paths param', function (assert) {
      const paths = ['/my/api/path', 'another/api/path'];
      const expectedPayload = { paths };
      this.server.post('/sys/capabilities-self', (schema, req) => {
        const actual = JSON.parse(req.requestBody);
        assert.true(true, 'request made to capabilities-self');
        assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
        return this.generateResponse({
          paths,
          capabilities: { '/my/api/path': ['read'], 'another/api/path': ['read'] },
        });
      });
      this.capabilities.request({ paths });
    });
  });

  test('fetchPathCapabilities: it makes request to capabilities-self with path param', function (assert) {
    const path = '/my/api/path';
    const expectedPayload = { paths: [path] };
    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse({ path, capabilities: ['read'] });
    });
    this.capabilities.fetchPathCapabilities(path);
  });

  test('fetchMultiplePaths: it makes request to capabilities-self with paths param', function (assert) {
    const paths = ['/my/api/path', 'another/api/path'];
    const expectedPayload = { paths };
    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse({
        paths,
        capabilities: { '/my/api/path': ['read'], 'another/api/path': ['read'] },
      });
    });
    this.capabilities.fetchMultiplePaths(paths);
  });

  module('specific methods', function () {
    const path = '/my/api/path';
    [
      {
        capabilities: ['read'],
        expectedRead: true, // expected computed properties based on response
        expectedUpdate: false,
      },
      {
        capabilities: ['update'],
        expectedRead: false,
        expectedUpdate: true,
      },
      {
        capabilities: ['deny'],
        expectedRead: false,
        expectedUpdate: false,
      },
      {
        capabilities: ['read', 'update'],
        expectedRead: true,
        expectedUpdate: true,
      },
    ].forEach(({ capabilities, expectedRead, expectedUpdate }) => {
      test(`canRead returns expected value for "${capabilities.join(', ')}"`, async function (assert) {
        this.server.post('/sys/capabilities-self', () => {
          return this.generateResponse({ path, capabilities });
        });

        const response = await this.capabilities.canRead(path);
        assert[expectedRead](response, `canRead returns ${expectedRead}`);
      });

      test(`canUpdate returns expected value for "${capabilities.join(', ')}"`, async function (assert) {
        this.server.post('/sys/capabilities-self', () => {
          return this.generateResponse({ path, capabilities });
        });
        const response = await this.capabilities.canUpdate(path);
        assert[expectedUpdate](response, `canUpdate returns ${expectedUpdate}`);
      });
    });
  });
});
