/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import ClusterRoute from 'vault/routes/vault/cluster';
import { INIT, UNSEAL, AUTH, CLUSTER, CLUSTER_INDEX, DR_REPLICATION_SECONDARY } from 'vault/lib/route-paths';
import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import sinon from 'sinon';

module('Integration | Route | ClusterRoute', function (hooks) {
  setupTest(hooks);

  function createClusterRoute(
    context,
    clusterModel = {},
    methods = {
      auth: { currentToken: null },
    }
  ) {
    context.owner.register('route:cluster-route-test', ClusterRoute);
    const instance = context.owner.lookup('route:cluster-route-test');
    for (const key of Object.keys(methods)) {
      if (typeof methods[key] === 'function') {
        sinon.stub(instance, key).callsFake(methods[key]);
      }
    }
    if (methods.auth) {
      instance.auth = methods.auth;
    }
    instance.modelFor = () => clusterModel;
    instance.router = { transitionTo: () => {} };
    return instance;
  }

  const INIT_TESTS = [
    {
      clusterState: { needsInit: true },
      expected: INIT,
      description: 'forwards to INIT when cluster needs init',
    },
    {
      clusterState: { needsInit: false, sealed: true },
      expected: UNSEAL,
      description: 'forwards to UNSEAL if sealed and initialized',
    },
    {
      clusterState: { needsInit: false, sealed: false },
      expected: AUTH,
      description: 'forwards to AUTH if unsealed and initialized',
    },
    {
      clusterState: { dr: { isSecondary: true } },
      expected: DR_REPLICATION_SECONDARY,
      description: 'forwards to DR_REPLICATION_SECONDARY if is a dr secondary',
    },
  ];

  for (const { clusterState, expected, description } of INIT_TESTS) {
    test(`#targetRouteName init case: ${expected}`, function (assert) {
      const subject = createClusterRoute(this, clusterState);
      subject.routeName = CLUSTER;
      assert.strictEqual(subject.targetRouteName(), expected, description);
    });
  }

  test('#targetRouteName happy path when not authed forwards to AUTH', function (assert) {
    const subject = createClusterRoute(
      this,
      { needsInit: false, sealed: false, dr: { isSecondary: false } },
      { auth: { currentToken: null } }
    );
    subject.router.currentRouteName = INIT;
    assert.strictEqual(subject.targetRouteName(), AUTH, 'forwards when inited and navigating to INIT');

    subject.router.currentRouteName = UNSEAL;
    assert.strictEqual(subject.targetRouteName(), AUTH, 'forwards when unsealed and navigating to UNSEAL');

    subject.router.currentRouteName = AUTH;
    assert.strictEqual(
      subject.targetRouteName(),
      AUTH,
      'forwards when non-authenticated and navigating to AUTH'
    );

    subject.router.currentRouteName = DR_REPLICATION_SECONDARY;
    assert.strictEqual(
      subject.targetRouteName(),
      AUTH,
      'forwards when not a DR secondary and navigating to DR_REPLICATION_SECONDARY'
    );
  });

  test('#transitionToTargetRoute', function (assert) {
    const redirectRouteURL = '/vault/secrets-engines/secret/create';
    const subject = createClusterRoute(this, { needsInit: false, sealed: false });
    subject.router.currentURL = redirectRouteURL;
    const spy = sinon.stub(subject.router, 'transitionTo');
    subject.transitionToTargetRoute();
    assert.ok(
      spy.calledWithExactly(AUTH, { queryParams: { redirect_to: redirectRouteURL } }),
      'calls transitionTo with the expected args'
    );

    spy.restore();
  });

  test('#transitionToTargetRoute with auth as a target', function (assert) {
    const subject = createClusterRoute(this, { needsInit: false, sealed: false });
    const spy = sinon.stub(subject.router, 'transitionTo');
    // in this case it's already transitioning to the AUTH route so we don't need to call transitionTo again
    subject.transitionToTargetRoute({ targetName: AUTH });
    assert.ok(spy.notCalled, 'transitionTo is not called');
    spy.restore();
  });

  test('#transitionToTargetRoute with auth target, coming from cluster route', function (assert) {
    const subject = createClusterRoute(this, { needsInit: false, sealed: false });
    const spy = sinon.stub(subject.router, 'transitionTo');
    subject.transitionToTargetRoute({ targetName: CLUSTER_INDEX });
    assert.ok(spy.calledWithExactly(AUTH), 'calls transitionTo without redirect_to');
    spy.restore();
  });
});
