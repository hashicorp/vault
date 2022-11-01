import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/redirect', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.router = this.owner.lookup('service:router');
    this.route = this.owner.lookup('route:vault/cluster/redirect');
    this.auth = this.owner.lookup('service:auth');
    this.originalTransition = this.router.replaceWith;
    this.router.replaceWith = sinon.stub().returns({
      followRedirects: function () {
        return {
          then: function (callback) {
            callback();
          },
        };
      },
    });
    this.setCurrentToken = (token) => {
      this.auth.setCluster(token);
      if (token) {
        this.auth.setTokenData(token, { token: 'foo' });
        this.auth.tokens = [token];
      }
    };
  });

  hooks.afterEach(function () {
    this.router.replaceWith = this.originalTransition;
  });

  test('it calls route', function (assert) {
    assert.ok(this.route);
  });

  test('it redirects to auth when unauthenticated', function (assert) {
    const originalToken = this.auth.currentToken;
    this.setCurrentToken(null);

    this.route.beforeModel({
      to: { queryParams: { redirect_to: 'vault/cluster/tools', namespace: 'admin' } },
    });

    assert.true(
      this.router.replaceWith.calledWithExactly('vault.cluster.auth', {
        queryParams: { namespace: 'admin' },
      }),
      'transitions to auth when not authenticated'
    );
    this.setCurrentToken(originalToken);
  });

  test('it redirects to cluster when authenticated without redirect param', function (assert) {
    const originalToken = this.auth.currentToken;
    this.setCurrentToken('s.xxxxxxxxx');

    this.route.beforeModel({ to: { queryParams: { foo: 'bar' } } });
    assert.true(
      this.router.replaceWith.calledWithExactly('vault.cluster', { queryParams: { foo: 'bar' } }),
      'transitions to cluster when authenticated but no redirect param'
    );
    this.setCurrentToken(originalToken);
  });

  test('it redirects to desired path when authenticated with redirect param', function (assert) {
    const originalToken = this.auth.currentToken;
    this.setCurrentToken('s.xxxxxxxxx');

    this.route.beforeModel({
      to: {
        queryParams: { redirect_to: 'vault/cluster/tools?namespace=admin', namespace: 'ns1', foo: 'bar' },
      },
    });

    assert.true(
      this.router.replaceWith.calledWithExactly('vault/cluster/tools?namespace=admin'),
      'transitions to redirect_to path when authenticated and removes other params'
    );
    this.setCurrentToken(originalToken);
  });
});
