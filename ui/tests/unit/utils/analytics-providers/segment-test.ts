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

  test('identify marks HVD clusters as dedicated and omits subscriptionId', function (assert) {
    const provider = new SegmentProvider();
    const identifyStub = sinon.stub(provider.client, 'identify');

    provider.identify('user-123', {
      // HVD may still surface an internal license id; it must NOT be sent as subscriptionId.
      licenseId: 'license-abc',
      clusterId: 'cluster-123',
      isEnterprise: true,
      isHvdManaged: true,
    });

    const props = identifyStub.firstCall.args[1];
    assert.strictEqual(
      props.productPlanType,
      'Vault dedicated',
      'HVD reports productPlanType "Vault dedicated"'
    );
    assert.strictEqual(props.instanceId, 'cluster-123', 'instanceId is still the cluster ID for HVD');
    assert.notOk('subscriptionId' in props, 'HVD omits subscriptionId (no HCP org id surfaced yet)');
    assert.notOk(
      'licenseId' in props,
      'HVD strips the raw licenseId trait so the internal license id is never sent'
    );
  });

  test('trackEvent carries productPlanType for HVD after identify', function (assert) {
    const provider = new SegmentProvider();
    const trackStub = sinon.stub(provider.client, 'track');

    provider.identify('user-123', {
      clusterId: 'cluster-123',
      isEnterprise: true,
      isHvdManaged: true,
    });

    provider.trackEvent('UI Interaction');

    const props = trackStub.firstCall.args[1];
    assert.strictEqual(
      props.productPlanType,
      'Vault dedicated',
      'event payload carries productPlanType for HVD'
    );
    assert.notOk('subscriptionId' in props, 'event payload omits subscriptionId for HVD');
  });

  test('start classifies HVD before identify runs', function (assert) {
    const provider = new SegmentProvider();
    provider.start({ enabled: false, write_key: '', isHvdManaged: true });

    assert.strictEqual(
      provider.productPlanType,
      'Vault dedicated',
      'productPlanType is set to dedicated at start, before identify runs'
    );
  });

  test('start defaults productPlanType to self-managed', function (assert) {
    const provider = new SegmentProvider();
    provider.start({ enabled: false, write_key: '' });

    assert.strictEqual(
      provider.productPlanType,
      'Vault self-managed',
      'productPlanType defaults to self-managed when no HVD hint is provided'
    );
  });
});
