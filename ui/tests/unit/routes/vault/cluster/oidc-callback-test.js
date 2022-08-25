import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/oidc-callback', function (hooks) {
  setupTest(hooks);
  const parentNs = 'admin';
  const childNs = 'admin/child-ns';
  const path = 'oidc';
  const customPath = 'oidc-dev';
  const code = 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T';
  const state = (ns) => {
    ns ? 'st_91ji6vR2sQ2zBiZSQkqJ' + `,ns=${ns}` : 'st_91ji6vR2sQ2zBiZSQkqJ';
  };

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.route = this.owner.lookup('route:vault/cluster/oidc-callback');
    this.windowStub = sinon.stub(window.opener, 'postMessage');
  });

  test('it calls route', function (assert) {
    assert.ok(this.route);
  });

  test('it uses namespace param from state not namespaceQueryParam from cluster with default path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';

    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: parentNs };
      return {
        auth_path: path,
        state: state(childNs),
        code,
      };
    };
    assert.ok(this.windowStub.calledWith, 'test');
    assert.propContains(
      this.route.afterModel(),
      {
        path,
        namespace: childNs,
        state: state(),
      },
      'state and namespace queryParams are correct'
    );
  });

  test('it uses namespace param from state not namespaceQueryParam from cluster with custom path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: parentNs };
      return {
        auth_path: customPath,
        state: state(childNs),
        code,
      };
    };
    assert.propContains(
      this.route.afterModel(),
      {
        path: customPath,
        namespace: childNs,
        state: state(),
      },
      'state ns takes precedence, state no longer has ns query'
    );
  });

  test('it uses namespace from namespaceQueryParam when no ns param from state', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: parentNs };
      return {
        auth_path: path,
        state: state(),
        code,
      };
    };
    assert.propContains(
      this.route.afterModel(),
      {
        path,
        namespace: parentNs,
        state: state(),
      },
      'namespace is from cluster namespaceQueryParam'
    );
  });

  test('it uses ns param from state when no namespaceQueryParam from cluster', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: path,
        state: state('ns1'),
        code,
      };
    };
    assert.propContains(
      this.route.afterModel(),
      {
        path,
        namespace: 'ns1',
        state: state(),
      },
      'it strips ns from state and uses as namespace param'
    );
  });
});
