import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

// window.opener = {
//   postMessage: () => {},
//   origin: 'http://localhost:4200',
// };
const origin = 'http://localhost:4200';

module('Unit | Route | vault/cluster/oidc-callback', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.originalOpener = window.opener;
    window.opener = {
      postMessage: () => {},
      origin, // todo need?
    };
    this.router = this.owner.lookup('service:router');
    this.route = this.owner.lookup('route:vault/cluster/oidc-callback');
    this.windowStub = sinon.stub(window.opener, 'postMessage');
    this.parentNs = 'admin';
    this.childNs = 'admin/child-ns';
    this.path = 'oidc';
    this.customPath = 'oidc-dev';
    this.code = 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T';
    this.state = (ns) => {
      ns ? 'st_91ji6vR2sQ2zBiZSQkqJ' + `,ns=${ns}` : 'st_91ji6vR2sQ2zBiZSQkqJ';
    };
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
    window.opener = this.originalOpener;
  });

  test('it calls route', function (assert) {
    assert.ok(this.route);
  });

  test('it uses namespace param from state not namespaceQueryParam from cluster with default path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: this.parentNs };
      return {
        auth_path: this.path,
        state: this.state(this.childNs),
        code: this.code,
      };
    };
    this.route.afterModel();

    assert.ok(this.windowStub.calledOnce, 'it is called');
    assert.propContains(
      this.windowStub.getCall(0).args[0],
      {
        code: 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T',
        namespace: 'admin',
        path: 'oidc',
        source: 'oidc-callback',
      },

      'calls correct'
    );
  });

  test('it uses namespace param from state not namespaceQueryParam from cluster with custom path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: parentNs };
      return {
        auth_path: this.customPath,
        state: this.state(this.childNs),
        code: this.code,
      };
    };
    assert.propContains(
      this.route.afterModel(),
      {
        path: this.customPath,
        namespace: this.childNs,
        state: this.state(),
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
