/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, skip, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { getParamsForCallback } from 'vault/routes/vault/cluster/oidc-callback';

module('Unit | Route | vault/cluster/oidc-callback', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.originalOpener = window.opener;
    window.opener = {
      postMessage: () => {},
    };
    this.route = this.owner.lookup('route:vault/cluster/oidc-callback');
    this.windowStub = sinon.stub(window.opener, 'postMessage');
    this.state = 'st_yOarDguU848w5YZuotLs';
    this.path = 'oidc';
    this.code = 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: this.path,
        code: this.code,
      };
    };
    this.callbackUrlQueryParams = (stateParam) => {
      switch (stateParam) {
        case '':
          window.history.pushState({}, '');
          break;
        case 'stateless':
          window.history.pushState({}, '', '?' + `code=${this.code}`);
          break;
        default:
          window.history.pushState({}, '', '?' + `code=${this.code}&state=${stateParam}`);
          break;
      }
    };
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
    window.opener = this.originalOpener;
    this.callbackUrlQueryParams('');
  });

  test('it calls route', function (assert) {
    assert.ok(this.route);
  });

  skip('it uses namespace param from state instead of cluster, with custom oidc path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams(encodeURIComponent(`${this.state},ns=test-ns`));
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: 'oidc-dev',
        code: this.code,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: 'oidc-dev',
        namespace: 'test-ns',
        state: this.state,
        source: 'oidc-callback',
      },
      'ns from state not cluster'
    );
  });

  skip('it uses namespace from cluster when state does not include ns param', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams(encodeURIComponent(this.state));
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: this.path,
        code: this.code,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: this.path,
        namespace: 'admin',
        state: this.state,
        source: 'oidc-callback',
      },
      `namespace is from cluster's namespaceQueryParam`
    );
  });

  skip('it correctly parses encoded, nested ns param from state', function (assert) {
    this.callbackUrlQueryParams(encodeURIComponent(`${this.state},ns=parent-ns/child-ns`));
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: this.path,
        namespace: 'parent-ns/child-ns',
        state: this.state,
        source: 'oidc-callback',
      },
      'it has correct nested ns from state and sets as namespace param'
    );
  });

  skip('the afterModel hook returns when both cluster and route params are empty strings', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams('');
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: '',
        code: '',
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        path: '',
        state: '',
        code: '',
        source: 'oidc-callback',
      },
      'model hook returns with empty params'
    );
  });

  skip('the afterModel hook returns when state param does not exist', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams('stateless');
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: 'oidc',
        state: '',
        source: 'oidc-callback',
      },
      'model hook returns empty string when state param nonexistent'
    );
  });

  skip('the afterModel hook returns when cluster ns exists and all route params are empty strings', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams('');
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'ns1' };
      return {
        auth_path: '',
        code: '',
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: '',
        namespace: 'ns1',
        path: '',
        source: 'oidc-callback',
        state: '',
      },
      'model hook returns with empty parameters'
    );
  });

  /*
  If authenticating to a namespace, most SSO providers return a callback url
  with a 'state' query param that includes a URI encoded namespace, example:
  '?code=BZBDVPMz0By2JTqulEMWX5-6rflW3A20UAusJYHEeFygJ&state=sst_yOarDguU848w5YZuotLs%2Cns%3Dadmin'

  Active Directory Federation Service (AD FS), instead, decodes the namespace portion:
  '?code=BZBDVPMz0By2JTqulEMWX5-6rflW3A20UAusJYHEeFygJ&state=st_yOarDguU848w5YZuotLs,ns=admin'

  'ns' isn't recognized as a separate param because there is no ampersand, so using this.paramsFor() returns
  a namespace-less state and authentication fails
  { state: 'st_yOarDguU848w5YZuotLs,ns' }
  */
  skip('it uses namespace when state param is not uri encoded', async function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.callbackUrlQueryParams(`${this.state},ns=admin`);
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        namespace: 'admin',
        path: this.path,
        source: 'oidc-callback',
        state: this.state,
      },
      'namespace is parsed correctly'
    );
  });

  module('getParamsForCallback helper fn', function () {
    test('it parses params correctly with regular inputs and no namespace', function (assert) {
      const qp = {
        state: 'my-state',
        code: 'my-code',
        path: 'oidc-path',
      };
      const searchString = `?code=my-code&state=my-state`;
      const results = getParamsForCallback(qp, searchString);
      assert.deepEqual(results, { source: 'oidc-callback', ...qp });
    });

    test('it parses params correctly regular inputs and namespace param', function (assert) {
      const params = {
        state: 'my-state',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'my-namespace',
      };
      const results = getParamsForCallback(params, '?code=my-code&state=my-state&namespace=my-namespace');
      assert.deepEqual(results, { source: 'oidc-callback', ...params });
    });

    test('it parses params correctly with regular inputs and namespace in state (unencoded)', function (assert) {
      const searchString = '?code=my-code&state=my-state,ns=foo/bar';
      const params = {
        state: 'my-state,ns', // Ember parses the QP incorrectly if unencoded
        code: 'my-code',
        path: 'oidc-path',
      };
      const results = getParamsForCallback(params, searchString);
      assert.deepEqual(results, {
        source: 'oidc-callback',
        ...params,
        state: 'my-state',
        namespace: 'foo/bar',
      });
    });

    test('it parses params correctly with regular inputs and namespace in state (encoded)', function (assert) {
      const qp = {
        state: 'my-state,ns=foo/bar', // Ember parses the QP correctly when encoded
        code: 'my-code',
        path: 'oidc-path',
      };
      const searchString = `?code=my-code&state=${encodeURIComponent(qp.state)}`;
      const results = getParamsForCallback(qp, searchString);
      assert.deepEqual(results, { source: 'oidc-callback', ...qp, state: 'my-state', namespace: 'foo/bar' });
    });

    test('namespace in state takes precedence over namespace in route (encoded)', function (assert) {
      const qp = {
        state: 'my-state,ns=foo/bar',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'other/ns',
      };
      const searchString = `?code=my-code&state=${encodeURIComponent(
        qp.state
      )}&namespace=${encodeURIComponent(qp.namespace)}`;
      const results = getParamsForCallback(qp, searchString);
      assert.deepEqual(results, {
        source: 'oidc-callback',
        ...qp,
        state: 'my-state',
        namespace: 'foo/bar',
      });
    });

    test('namespace in state takes precedence over namespace in route (unencoded)', function (assert) {
      const qp = {
        state: 'my-state,ns',
        code: 'my-code',
        path: 'oidc-path',
        namespace: 'other/ns',
      };
      const searchString = `?code=${qp.code}&state=my-state,ns=foo/bar&namespace=${qp.namespace}`;
      const results = getParamsForCallback(qp, searchString);
      assert.deepEqual(results, {
        source: 'oidc-callback',
        ...qp,
        state: 'my-state',
        namespace: 'foo/bar',
      });
    });

    test('it parses params correctly when window.location.search is empty (HCP scenario)', function (assert) {
      const params = {
        state: 'some-state,ns=admin/child-ns',
        code: 'my-code',
        namespace: 'admin',
        path: 'oidc-path',
      };
      const results = getParamsForCallback(params, '');
      assert.deepEqual(results, {
        source: 'oidc-callback',
        code: 'my-code',
        path: 'oidc-path',
        state: 'some-state',
        namespace: 'admin/child-ns',
      });
    });
  });
});
