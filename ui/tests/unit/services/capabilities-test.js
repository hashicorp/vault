/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';

module('Unit | Service | capabilities', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.capabilities = this.owner.lookup('service:capabilities');
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

  test('fetch: it makes request to capabilities-self', async function (assert) {
    const paths = ['/my/api/path', 'another/api/path'];
    const expectedPayload = { paths };

    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse({
        paths,
        capabilities: { '/my/api/path': ['read', 'list'], 'another/api/path': ['read', 'delete'] },
      });
    });

    const actual = await this.capabilities.fetch(paths);
    const expected = {
      '/my/api/path': {
        canCreate: false,
        canDelete: false,
        canList: true,
        canPatch: false,
        canRead: true,
        canSudo: false,
        canUpdate: false,
      },
      'another/api/path': {
        canCreate: false,
        canDelete: true,
        canList: false,
        canPatch: false,
        canRead: true,
        canSudo: false,
        canUpdate: false,
      },
    };
    assert.propEqual(actual, expected, `it returns expected response: ${JSON.stringify(actual)}`);
  });

  test('fetch: it defaults to true if the capabilities request fails', async function (assert) {
    // don't stub endpoint which causes request to fail
    const paths = ['/my/api/path', 'another/api/path'];
    const actual = await this.capabilities.fetch(paths);
    const expected = {
      '/my/api/path': {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
      'another/api/path': {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
    };
    assert.propEqual(actual, expected, `it returns expected response: ${JSON.stringify(actual)}`);
  });

  test('fetch: it defaults to true if no model is found', async function (assert) {
    const paths = ['/my/api/path', 'another/api/path'];
    const expectedPayload = { paths };

    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse({
        paths: ['/my/api/path'],
        capabilities: { '/my/api/path': ['read', 'list'] },
      });
    });

    const actual = await this.capabilities.fetch(paths);
    const expected = {
      '/my/api/path': {
        canCreate: false,
        canDelete: false,
        canList: true,
        canPatch: false,
        canRead: true,
        canSudo: false,
        canUpdate: false,
      },
      'another/api/path': {
        canCreate: true,
        canDelete: true,
        canList: true,
        canPatch: true,
        canRead: true,
        canSudo: true,
        canUpdate: true,
      },
    };
    assert.propEqual(actual, expected, `it returns expected response: ${JSON.stringify(actual)}`);
  });

  test('fetch: it tracks the requested paths', async function (assert) {
    const paths = ['/my/api/path', 'another/api/path'];

    this.server.post('/sys/capabilities-self', () => {
      return this.generateResponse({
        paths,
        capabilities: { '/my/api/path': ['read'], 'another/api/path': ['update'] },
      });
    });

    assert.strictEqual(this.capabilities.requestedPaths.size, 0, 'requestedPaths is empty before fetch');

    await this.capabilities.fetch(paths);

    assert.strictEqual(this.capabilities.requestedPaths.size, 2, 'requestedPaths contains 2 items');
    assert.true(this.capabilities.requestedPaths.has('/my/api/path'), 'contains first path');
    assert.true(this.capabilities.requestedPaths.has('another/api/path'), 'contains second path');
  });

  test('fetch: it replaces requestedPaths on each call', async function (assert) {
    const firstPaths = ['/path/one', '/path/two'];
    const secondPaths = ['/path/three'];

    this.server.post('/sys/capabilities-self', () => {
      return this.generateResponse({
        paths: firstPaths,
        capabilities: { '/path/one': ['read'], '/path/two': ['read'], '/path/three': ['read'] },
      });
    });

    await this.capabilities.fetch(firstPaths);
    assert.strictEqual(this.capabilities.requestedPaths.size, 2, 'initially has 2 paths');

    await this.capabilities.fetch(secondPaths);
    assert.strictEqual(this.capabilities.requestedPaths.size, 1, 'updated to have 1 path');
    assert.true(this.capabilities.requestedPaths.has('/path/three'), 'contains new path');
    assert.false(this.capabilities.requestedPaths.has('/path/one'), 'no longer contains old path');
  });

  test('fetchPathCapabilities: it makes request to capabilities-self and returns capabilities for single path', async function (assert) {
    const path = '/my/api/path';
    const expectedPayload = { paths: [path] };

    this.server.post('/sys/capabilities-self', (schema, req) => {
      const actual = JSON.parse(req.requestBody);
      assert.true(true, 'request made to capabilities-self');
      assert.propEqual(actual, expectedPayload, `request made with path: ${JSON.stringify(actual)}`);
      return this.generateResponse({ path, capabilities: ['read'] });
    });

    const actual = await this.capabilities.fetchPathCapabilities(path);
    const expected = {
      canCreate: false,
      canDelete: false,
      canList: false,
      canPatch: false,
      canRead: true,
      canSudo: false,
      canUpdate: false,
    };
    assert.propEqual(actual, expected, 'returns capabilities for provided path');
  });

  test('pathFor: it should resolve the correct path to fetch capabilities', async function (assert) {
    const syncActivate = this.capabilities.pathFor('syncActivate');
    assert.strictEqual(
      syncActivate,
      'sys/activation-flags/secrets-sync/activate',
      'pathFor returns expected path for syncActivate'
    );

    const syncDestination = this.capabilities.pathFor('syncDestination', { type: 'aws-sm', name: 'foo' });
    assert.strictEqual(
      syncDestination,
      'sys/sync/destinations/aws-sm/foo',
      'pathFor returns expected path for syncDestination'
    );
  });

  test('for: it should fetch capabilities for single path using pathFor and fetchPathCapabilities methods', async function (assert) {
    const pathForStub = sinon.spy(this.capabilities, 'pathFor');
    const fetchPathCapabilitiesStub = sinon.stub(this.capabilities, 'fetchPathCapabilities').resolves();
    await this.capabilities.for('customMessages', { id: 'foo' });
    assert.true(pathForStub.calledWith('customMessages', { id: 'foo' }), 'pathFor called with expected args');
    assert.true(
      fetchPathCapabilitiesStub.calledWith('sys/config/ui/custom-messages/foo'),
      'fetchPathCapabilities called with expected args'
    );
  });

  module('specific methods', function () {
    const path = '/my/api/path';
    [
      {
        capabilities: ['read'],
        expectedRead: true, // expected computed properties based on response
        expectedUpdate: false,
        expectedPatch: false,
      },
      {
        capabilities: ['update'],
        expectedRead: false,
        expectedUpdate: true,
        expectedPatch: false,
      },
      {
        capabilities: ['patch'],
        expectedRead: false,
        expectedUpdate: false,
        expectedPatch: true,
      },
      {
        capabilities: ['deny'],
        expectedRead: false,
        expectedUpdate: false,
        expectedPatch: false,
      },
      {
        capabilities: ['read', 'update'],
        expectedRead: true,
        expectedUpdate: true,
        expectedPatch: false,
      },
    ].forEach(({ capabilities, expectedRead, expectedUpdate, expectedPatch }) => {
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

      test(`canPatch returns expected value for "${capabilities.join(', ')}"`, async function (assert) {
        this.server.post('/sys/capabilities-self', () => {
          return this.generateResponse({ path, capabilities });
        });
        const response = await this.capabilities.canPatch(path);
        assert[expectedPatch](response, `canPatch returns ${expectedPatch}`);
      });
    });
  });

  module('within a namespace', function (hooks) {
    // capabilities within namespaces are queried at the user's root namespace with a path that includes
    // the relative namespace.
    hooks.beforeEach(function () {
      this.nsSvc = this.owner.lookup('service:namespace');
      this.nsSvc.path = 'ns1';
    });

    test('fetchPathCapabilities works as expected', async function (assert) {
      const ns = this.nsSvc.path;
      const path = '/my/api/path';
      const expectedAttrs = {
        canCreate: false,
        canDelete: false,
        canList: true,
        canPatch: false,
        canRead: true,
        canSudo: false,
        canUpdate: false,
      };

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const actual = JSON.parse(req.requestBody);
        assert.strictEqual(req.url, '/v1/sys/capabilities-self', 'request made to capabilities-self');
        assert.propEqual(
          actual.paths,
          [`${ns}/my/api/path`],
          `request made with path: ${JSON.stringify(actual)}`
        );
        return this.generateResponse({
          path: `${ns}${path}`,
          capabilities: ['read', 'list'],
        });
      });

      const actual = await this.capabilities.fetchPathCapabilities(path);

      Object.keys(expectedAttrs).forEach(function (key) {
        assert.strictEqual(
          actual[key],
          expectedAttrs[key],
          `record has expected value for ${key}: ${actual[key]}`
        );
      });
    });

    test('fetch works as expected', async function (assert) {
      const ns = this.nsSvc.path;
      // there was a bug when stripping the relativeNamespace from the key in the response data
      // this would result in a leading slash in the returned key causing a mismatch if the path was provided without a leading slash
      // ensure the provided path is returned by passing at least one path without a leading slash
      const paths = ['my/api/path', '/another/api/path'];
      const expectedPayload = [`${ns}/my/api/path`, `${ns}/another/api/path`];

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const actual = JSON.parse(req.requestBody);
        assert.strictEqual(req.url, '/v1/sys/capabilities-self', 'request made to capabilities-self');
        assert.propEqual(actual.paths, expectedPayload, `request made with paths: ${JSON.stringify(actual)}`);
        return this.generateResponse({
          paths: expectedPayload,
          capabilities: {
            [`${ns}/my/api/path`]: ['read', 'list'],
            [`${ns}/another/api/path`]: ['update', 'patch'],
          },
        });
      });

      const actual = await this.capabilities.fetch(paths);
      const expected = {
        'my/api/path': {
          canCreate: false,
          canDelete: false,
          canList: true,
          canPatch: false,
          canRead: true,
          canSudo: false,
          canUpdate: false,
        },
        '/another/api/path': {
          canCreate: false,
          canDelete: false,
          canList: false,
          canPatch: true,
          canRead: false,
          canSudo: false,
          canUpdate: true,
        },
      };
      assert.propEqual(actual, expected, 'method returns expected response');
    });

    /* 
    The setup in this test simulates a user whose auth method is mounted in the "root" namespace 
    but their policy only grants access to paths in the context of the "ns1" namespace.

    * ~Example policy paths~ *
    # explicitly grants access to read "my-secret" in the kv engine mounted in the "ns1" namespace
    path "ns1/kv/data/my-secret" {
      capabilities = ["read", "delete"]
    }
    
    # alternatively, their policy could grant access to read everything within the "ns1" namespace
    path "ns1/*" {
      capabilities = ["read"]
    }
  */
    test(`if the user's root namespace is "root" and the resource is in a child namespace`, async function (assert) {
      assert.expect(2);
      const ns = this.nsSvc.path;
      const paths = ['my/api/path', '/another/api/path'];
      const expectedPayload = [`${ns}/my/api/path`, `${ns}/another/api/path`];

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const nsHeader = req.requestHeaders['x-vault-namespace'];
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(nsHeader, '', 'request is made in the context of the "root" namespace');
        assert.propEqual(payload.paths, expectedPayload, `paths include the relative namespace`);
        return req.passthrough();
      });
      await this.capabilities.fetch(paths);
    });

    /* 
    The setup in this test simulates a user whose root namespace is "root" and 
    they are accessing a resource at a nested namespace: "ns1/child". 
    */
    test(`if the user's root namespace is "root" and the resource is in a grandchild`, async function (assert) {
      assert.expect(2);
      // the path in the namespace service is always the FULL namespace path of the current context
      this.nsSvc.path = 'ns1/child';

      const paths = ['my/api/path', '/another/api/path'];
      const expectedPaths = ['ns1/child/my/api/path', 'ns1/child/another/api/path'];

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const nsHeader = req.requestHeaders['x-vault-namespace'];
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(nsHeader, '', 'request is made in the context of the "root" namespace');
        assert.propEqual(payload.paths, expectedPaths, `paths include the relative namespace`);
        return req.passthrough();
      });

      await this.capabilities.fetch(paths);
    });

    /* 
    The setup in this test simulates a user whose auth method is mounted in the "ns1" namespace and so cannot log in directly to "root" at all.
    Since this user's context (along with their policy) is exclusively "ns1" the paths do not include the namespace.

    * ~Example policy paths~ *
    path "kv/data/my-secret" {
      capabilities = ["read", "delete"]
    }
    */
    test(`if the user's root namespace is an immediate child of "root" and they are accessing resources in the same namespace context`, async function (assert) {
      assert.expect(2);

      const ns = this.nsSvc.path;
      const authService = this.owner.lookup('service:auth');
      const authStub = sinon.stub(authService, 'authData').value({ userRootNamespace: ns });

      const paths = ['my/api/path', '/another/api/path'];

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const nsHeader = req.requestHeaders['x-vault-namespace'];
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(nsHeader, 'ns1', 'request is made in the context of the "ns1" namespace');
        assert.propEqual(
          payload.paths,
          paths,
          'paths do not include the namespace because request header manages context'
        );
        return req.passthrough();
      });

      await this.capabilities.fetch(paths);
      authStub.restore();
    });

    /* 
    The setup in this test simulates a user whose root namespace is "ns1" and 
    they are accessing a resource at a namespace one level deeper in "ns1/child". 
    */
    test(`if the user's root namespace is a child of "root" and the resource is nested one more level`, async function (assert) {
      assert.expect(2);
      // the path in the namespace service is always the FULL namespace path of the current context
      this.nsSvc.path = 'ns1/child';
      const authService = this.owner.lookup('service:auth');
      const authStub = sinon.stub(authService, 'authData').value({ userRootNamespace: 'ns1' });

      const paths = ['my/api/path', '/another/api/path'];
      const expectedPaths = ['child/my/api/path', 'child/another/api/path'];

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const nsHeader = req.requestHeaders['x-vault-namespace'];
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(nsHeader, 'ns1', 'request is made in the context of the "ns1" namespace');
        assert.propEqual(payload.paths, expectedPaths, 'paths include the relative namespace');
        return req.passthrough();
      });

      await this.capabilities.fetch(paths);
      authStub.restore();
    });

    /* 
    The setup in this test simulates a user whose root namespace is "ns1/child" and 
    they are accessing a resource in the same context. 
    */
    test(`if the user's root namespace is a grandchild of "root" and the resource is in the same context`, async function (assert) {
      assert.expect(2);
      // the path in the namespace service is always the FULL namespace path of the current context
      this.nsSvc.path = 'ns1/child';
      const authService = this.owner.lookup('service:auth');
      const authStub = sinon.stub(authService, 'authData').value({ userRootNamespace: 'ns1/child' });

      const paths = ['my/api/path', '/another/api/path'];

      this.server.post('/sys/capabilities-self', (schema, req) => {
        const nsHeader = req.requestHeaders['x-vault-namespace'];
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(
          nsHeader,
          'ns1/child',
          'request is made in the context of the "ns1/child" namespace'
        );
        assert.propEqual(
          payload.paths,
          paths,
          'paths do not include namespace because header manages namespace context'
        );
        return req.passthrough();
      });

      await this.capabilities.fetch(paths);
      authStub.restore();
    });
  });
});
