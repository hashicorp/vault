/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';

import { SegmentProvider } from 'vault/utils/analytics-providers/segment';

module('Unit | Utils | analytics providers | segment', function (hooks) {
  hooks.afterEach(function () {
    sinon.restore();
  });

  test('identify sets instanceId from clusterId and subscriptionId from licenseId', function (assert) {
    const provider = new SegmentProvider();
    const identifyStub = sinon.stub(provider.client, 'identify');

    provider.identify('user-123', {
      licenseId: 'license-abc',
      clusterId: 'cluster-123',
      isEnterprise: true,
    });

    assert.true(identifyStub.calledOnce, 'identify is called');
    assert.propContains(
      identifyStub.firstCall.args[1],
      {
        instanceId: 'cluster-123',
        subscriptionId: 'license-abc',
      },
      'instanceId is the cluster ID and subscriptionId is the license ID'
    );
  });

  test('identify omits subscriptionId on community clusters (no license)', function (assert) {
    const provider = new SegmentProvider();
    const identifyStub = sinon.stub(provider.client, 'identify');

    provider.identify('user-123', {
      clusterId: 'cluster-123',
      isEnterprise: false,
    });

    assert.true(identifyStub.calledOnce, 'identify is called');
    const props = identifyStub.firstCall.args[1];
    assert.strictEqual(props.instanceId, 'cluster-123', 'instanceId is still the cluster ID');
    assert.notOk('subscriptionId' in props, 'subscriptionId is omitted when there is no license');
  });

  test('trackEvent reuses instanceId and omits subscriptionId for community after identify', function (assert) {
    const provider = new SegmentProvider();
    const trackStub = sinon.stub(provider.client, 'track');

    provider.identify('user-123', {
      clusterId: 'cluster-123',
      isEnterprise: false,
    });

    provider.trackEvent('UI Interaction');

    assert.true(trackStub.calledOnce, 'track is called');
    const props = trackStub.firstCall.args[1];
    assert.strictEqual(
      props.instanceId,
      'cluster-123',
      'event payload includes the cluster ID as instanceId'
    );
    assert.notOk('subscriptionId' in props, 'event payload omits subscriptionId for community');
  });
});
