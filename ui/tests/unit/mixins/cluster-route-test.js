import { assign } from '@ember/polyfills';
import EmberObject from '@ember/object';
import ClusterRouteMixin from 'vault/mixins/cluster-route';
import { INIT, UNSEAL, AUTH, CLUSTER, DR_REPLICATION_SECONDARY } from 'vault/mixins/cluster-route';
import { module, test } from 'qunit';

module('Unit | Mixin | cluster route', function() {
  function createClusterRoute(
    clusterModel = {},
    methods = { hasKeyData: () => false, authToken: () => null }
  ) {
    let ClusterRouteObject = EmberObject.extend(
      ClusterRouteMixin,
      assign(methods, { clusterModel: () => clusterModel })
    );
    return ClusterRouteObject.create();
  }

  test('#targetRouteName init', function(assert) {
    let subject = createClusterRoute({ needsInit: true });
    subject.routeName = CLUSTER;
    assert.equal(subject.targetRouteName(), INIT, 'forwards to INIT when cluster needs init');

    subject = createClusterRoute({ needsInit: false, sealed: true });
    subject.routeName = CLUSTER;
    assert.equal(subject.targetRouteName(), UNSEAL, 'forwards to UNSEAL if sealed and initialized');

    subject = createClusterRoute({ needsInit: false });
    subject.routeName = CLUSTER;
    assert.equal(subject.targetRouteName(), AUTH, 'forwards to AUTH if unsealed and initialized');

    subject = createClusterRoute({ dr: { isSecondary: true } });
    subject.routeName = CLUSTER;
    assert.equal(
      subject.targetRouteName(),
      DR_REPLICATION_SECONDARY,
      'forwards to DR_REPLICATION_SECONDARY if is a dr secondary'
    );
  });

  test('#targetRouteName when #hasDataKey is true', function(assert) {
    let subject = createClusterRoute(
      { needsInit: false, sealed: true },
      { hasKeyData: () => true, authToken: () => null }
    );

    subject.routeName = CLUSTER;
    assert.equal(subject.targetRouteName(), INIT, 'still land on INIT if there are keys on the controller');

    subject.routeName = UNSEAL;
    assert.equal(subject.targetRouteName(), UNSEAL, 'allowed to proceed to unseal');

    subject = createClusterRoute(
      { needsInit: false, sealed: false },
      { hasKeyData: () => true, authToken: () => null }
    );

    subject.routeName = AUTH;
    assert.equal(subject.targetRouteName(), AUTH, 'allowed to proceed to auth');
  });

  test('#targetRouteName happy path forwards to CLUSTER route', function(assert) {
    let subject = createClusterRoute(
      { needsInit: false, sealed: false, dr: { isSecondary: false } },
      { hasKeyData: () => false, authToken: () => 'a token' }
    );
    subject.routeName = INIT;
    assert.equal(subject.targetRouteName(), CLUSTER, 'forwards when inited and navigating to INIT');

    subject.routeName = UNSEAL;
    assert.equal(subject.targetRouteName(), CLUSTER, 'forwards when unsealed and navigating to UNSEAL');

    subject.routeName = AUTH;
    assert.equal(subject.targetRouteName(), CLUSTER, 'forwards when authenticated and navigating to AUTH');

    subject.routeName = DR_REPLICATION_SECONDARY;
    assert.equal(
      subject.targetRouteName(),
      CLUSTER,
      'forwards when not a DR secondary and navigating to DR_REPLICATION_SECONDARY'
    );
  });
});
