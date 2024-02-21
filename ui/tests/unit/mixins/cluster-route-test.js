/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import ClusterRouteMixin from 'vault/mixins/cluster-route';
import {
  INIT,
  UNSEAL,
  AUTH,
  CLUSTER,
  CLUSTER_INDEX,
  DR_REPLICATION_SECONDARY,
  REDIRECT,
} from 'vault/lib/route-paths';
import { module, test } from 'qunit';
import sinon from 'sinon';

module('Unit | Mixin | cluster route', function () {
  function createClusterRoute(
    clusterModel = {},
    methods = {
      router: { transitionTo: () => {} },
      hasKeyData: () => false,
      authToken: () => null,
      transitionTo: () => {},
    }
  ) {
    const ClusterRouteObject = EmberObject.extend(
      ClusterRouteMixin,
      Object.assign(methods, { clusterModel: () => clusterModel })
    );
    return ClusterRouteObject.create();
  }

  test('#targetRouteName init', function (assert) {
    let subject = createClusterRoute({ needsInit: true });
    subject.routeName = CLUSTER;
    assert.strictEqual(subject.targetRouteName(), INIT, 'forwards to INIT when cluster needs init');

    subject = createClusterRoute({ needsInit: false, sealed: true });
    subject.routeName = CLUSTER;
    assert.strictEqual(subject.targetRouteName(), UNSEAL, 'forwards to UNSEAL if sealed and initialized');

    subject = createClusterRoute({ needsInit: false });
    subject.routeName = CLUSTER;
    assert.strictEqual(subject.targetRouteName(), AUTH, 'forwards to AUTH if unsealed and initialized');

    subject = createClusterRoute({ dr: { isSecondary: true } });
    subject.routeName = CLUSTER;
    assert.strictEqual(
      subject.targetRouteName(),
      DR_REPLICATION_SECONDARY,
      'forwards to DR_REPLICATION_SECONDARY if is a dr secondary'
    );
  });

  test('#targetRouteName when #hasDataKey is true', function (assert) {
    let subject = createClusterRoute(
      { needsInit: false, sealed: true },
      { hasKeyData: () => true, authToken: () => null }
    );

    subject.routeName = CLUSTER;
    assert.strictEqual(
      subject.targetRouteName(),
      INIT,
      'still land on INIT if there are keys on the controller'
    );

    subject.routeName = UNSEAL;
    assert.strictEqual(subject.targetRouteName(), UNSEAL, 'allowed to proceed to unseal');

    subject = createClusterRoute(
      { needsInit: false, sealed: false },
      { hasKeyData: () => true, authToken: () => null }
    );

    subject.routeName = AUTH;
    assert.strictEqual(subject.targetRouteName(), AUTH, 'allowed to proceed to auth');
  });

  test('#targetRouteName happy path forwards to CLUSTER route', function (assert) {
    const subject = createClusterRoute(
      { needsInit: false, sealed: false, dr: { isSecondary: false } },
      { hasKeyData: () => false, authToken: () => 'a token' }
    );
    subject.routeName = INIT;
    assert.strictEqual(subject.targetRouteName(), CLUSTER, 'forwards when inited and navigating to INIT');

    subject.routeName = UNSEAL;
    assert.strictEqual(subject.targetRouteName(), CLUSTER, 'forwards when unsealed and navigating to UNSEAL');

    subject.routeName = AUTH;
    assert.strictEqual(
      subject.targetRouteName(),
      REDIRECT,
      'forwards when authenticated and navigating to AUTH'
    );

    subject.routeName = DR_REPLICATION_SECONDARY;
    assert.strictEqual(
      subject.targetRouteName(),
      CLUSTER,
      'forwards when not a DR secondary and navigating to DR_REPLICATION_SECONDARY'
    );
  });

  test('#targetRouteName happy path when not authed forwards to AUTH', function (assert) {
    const subject = createClusterRoute(
      { needsInit: false, sealed: false, dr: { isSecondary: false } },
      { hasKeyData: () => false, authToken: () => null }
    );
    subject.routeName = INIT;
    assert.strictEqual(subject.targetRouteName(), AUTH, 'forwards when inited and navigating to INIT');

    subject.routeName = UNSEAL;
    assert.strictEqual(subject.targetRouteName(), AUTH, 'forwards when unsealed and navigating to UNSEAL');

    subject.routeName = AUTH;
    assert.strictEqual(
      subject.targetRouteName(),
      AUTH,
      'forwards when non-authenticated and navigating to AUTH'
    );

    subject.routeName = DR_REPLICATION_SECONDARY;
    assert.strictEqual(
      subject.targetRouteName(),
      AUTH,
      'forwards when not a DR secondary and navigating to DR_REPLICATION_SECONDARY'
    );
  });

  test('#transitionToTargetRoute', function (assert) {
    const redirectRouteURL = '/vault/secrets/secret/create';
    const subject = createClusterRoute({ needsInit: false, sealed: false });
    subject.router.currentURL = redirectRouteURL;
    const spy = sinon.spy(subject.router, 'transitionTo');
    subject.transitionToTargetRoute();
    assert.ok(
      spy.calledWithExactly(AUTH, { queryParams: { redirect_to: redirectRouteURL } }),
      'calls transitionTo with the expected args'
    );

    spy.restore();
  });

  test('#transitionToTargetRoute with auth as a target', function (assert) {
    const subject = createClusterRoute({ needsInit: false, sealed: false });
    const spy = sinon.spy(subject, 'transitionTo');
    // in this case it's already transitioning to the AUTH route so we don't need to call transitionTo again
    subject.transitionToTargetRoute({ targetName: AUTH });
    assert.ok(spy.notCalled, 'transitionTo is not called');
    spy.restore();
  });

  test('#transitionToTargetRoute with auth target, coming from cluster route', function (assert) {
    const subject = createClusterRoute({ needsInit: false, sealed: false });
    const spy = sinon.spy(subject.router, 'transitionTo');
    subject.transitionToTargetRoute({ targetName: CLUSTER_INDEX });
    assert.ok(spy.calledWithExactly(AUTH), 'calls transitionTo without redirect_to');
    spy.restore();
  });
});
