import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/oidc-callback', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.originalOpener = window.opener;
    window.opener = {
      postMessage: () => {},
    };
    this.route = this.owner.lookup('route:vault/cluster/oidc-callback');
    this.windowStub = sinon.stub(window.opener, 'postMessage');
    this.path = 'oidc';
    this.code = 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T';
    this.state = (ns) => {
      return ns ? 'st_91ji6vR2sQ2zBiZSQkqJ' + `,ns=${ns}` : 'st_91ji6vR2sQ2zBiZSQkqJ';
    };
    this.pushQueryParam = (queryString) => {
      window.history.pushState({}, '', '?' + queryString);
    };
  });

  hooks.afterEach(function () {
    this.windowStub.restore();
    window.opener = this.originalOpener;
    this.pushQueryParam('');
  });

  test('it calls route', function (assert) {
    assert.ok(this.route);
  });

  test('it uses namespace param from state not namespaceQueryParam from cluster with default path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: this.path,
        state: this.state('admin/child-ns'),
        code: this.code,
      };
    };
    this.route.afterModel();

    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: 'lTazRXEwKfyGKBUCo5TyLJzdIt39YniBJOXPABiRMkL0T',
        namespace: 'admin/child-ns',
        path: 'oidc',
        source: 'oidc-callback',
        state: this.state(),
      },
      'namespace param is from state, ns=admin/child-ns'
    );
  });

  test('it uses namespace param from state not namespaceQueryParam from cluster with custom path', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: 'oidc-dev',
        state: this.state('admin/child-ns'),
        code: this.code,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: 'oidc-dev',
        namespace: 'admin/child-ns',
        state: this.state(),
        source: 'oidc-callback',
      },
      'state ns takes precedence, state no longer has ns query'
    );
  });

  test(`it uses namespace from namespaceQueryParam when state does not include: ',ns=some-namespace'`, function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'admin' };
      return {
        auth_path: this.path,
        state: this.state(),
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
        state: this.state(),
        source: 'oidc-callback',
      },
      'namespace is from cluster namespaceQueryParam'
    );
  });

  test('it uses ns param from state when no namespaceQueryParam from cluster', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: this.path,
        state: this.state('ns1'),
        code: this.code,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        path: this.path,
        namespace: 'ns1',
        state: this.state(),
        source: 'oidc-callback',
      },
      'it strips ns from state and uses as namespace param'
    );
  });

  test('the afterModel hook returns when both cluster and route params are empty strings', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: '',
        state: '',
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

  test('the afterModel hook returns when state param does not exist', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: this.path,
      };
    };
    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: '',
        path: 'oidc',
        state: '',
        source: 'oidc-callback',
      },
      'model hook returns empty string when state param nonexistent'
    );
  });

  test('the afterModel hook returns when cluster namespaceQueryParam exists and all route params are empty strings', function (assert) {
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: 'ns1' };
      return {
        auth_path: '',
        state: '',
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
  '?code=BZBDVPMz0By2JTqulEMWX5-6rflW3A20UAusJYHEeFygJ&state=st_EC8PbzZ7XUQ0ClEgssS9%2Cns%3Dadmin'    

  Active Directory Federation Service (AD FS), instead, decodes the namespace portion:
  '?code=BZBDVPMz0By2JTqulEMWX5-6rflW3A20UAusJYHEeFygJ&state=st_gVRGT4TJe2RpvHNX5HV0,ns=admin'
  
  'ns' isn't recognized as a separate param because there is no ampersand, so using this.paramsFor() returns
  a namespace-less state and authentication fails
  { state: 'st_91ji6vR2sQ2zBiZSQkqJ,ns' }
  */
  test('is parses the namespace from a uri with decoded state param', async function (assert) {
    this.pushQueryParam(`?code=${this.code}&state=st_gVRGT4TJe2RpvHNX5HV0,ns=admin`);
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: this.path,
        state: 'st_91ji6vR2sQ2zBiZSQkqJ,ns',
        code: this.code,
      };
    };

    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        namespace: 'admin',
        path: this.path,
        source: 'oidc-callback',
        state: 'st_gVRGT4TJe2RpvHNX5HV0',
      },
      'namespace is passed to window queryParams'
    );
  });

  test('is parses the namespace from a uri with encoded state param', async function (assert) {
    this.pushQueryParam(`?code=${this.code}&state=st_EC8PbzZ7XUQ0ClEgssS9%2Cns%3Dadmin`);
    this.routeName = 'vault.cluster.oidc-callback';
    this.route.paramsFor = (path) => {
      if (path === 'vault.cluster') return { namespaceQueryParam: '' };
      return {
        auth_path: this.path,
        state: 'st_EC8PbzZ7XUQ0ClEgssS9,ns=admin',
        code: this.code,
      };
    };

    this.route.afterModel();
    assert.propEqual(
      this.windowStub.lastCall.args[0],
      {
        code: this.code,
        namespace: 'admin',
        path: this.path,
        source: 'oidc-callback',
        state: 'st_EC8PbzZ7XUQ0ClEgssS9',
      },
      'namespace is passed to window queryParams'
    );
  });
});
