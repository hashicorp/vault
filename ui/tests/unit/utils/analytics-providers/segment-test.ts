/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';

import localStorage from 'vault/lib/local-storage';
import { SegmentProvider } from 'vault/utils/analytics-providers/segment';
import { setStringPreference, STRING_PREFERENCES } from 'vault/utils/preferences';

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

  module('trackEvent roles from persona preference', function (hooks) {
    const rolesFor = (stored?: string) => {
      if (stored !== undefined) setStringPreference('persona', stored);
      const provider = new SegmentProvider();
      const trackStub = sinon.stub(provider.client, 'track');
      provider.trackEvent('UI Interaction');
      return trackStub.firstCall.args[1]?.['roles'];
    };

    hooks.beforeEach(function () {
      localStorage.removeItem(STRING_PREFERENCES['persona']?.key ?? '');
    });

    test('sends "unknown" when no persona is stored', function (assert) {
      assert.deepEqual(rolesFor(undefined), ['unknown']);
    });

    test('sends the selected predefined persona', function (assert) {
      assert.deepEqual(rolesFor(JSON.stringify({ value: 'developer' })), ['developer']);
    });

    test('sends "other" when other is selected without a custom role', function (assert) {
      assert.deepEqual(rolesFor(JSON.stringify({ value: 'other' })), ['other']);
    });

    test('sends "other-<kebab-role>" when a custom role is provided', function (assert) {
      assert.deepEqual(rolesFor(JSON.stringify({ value: 'other', customRole: 'DevOps SRE' })), [
        'other-devops-sre',
      ]);
    });

    test('sends "other" when the custom role is only whitespace', function (assert) {
      assert.deepEqual(rolesFor(JSON.stringify({ value: 'other', customRole: '   ' })), ['other']);
    });
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

  test('reset clears the Segment identity', function (assert) {
    const provider = new SegmentProvider();
    const resetStub = sinon.stub(provider.client, 'reset');

    provider.reset();

    assert.true(resetStub.calledOnce, 'client.reset() is called to clear ajs_user_id and anonymousId');
  });

  module('purgeStorage', function (hooks) {
    const SEGMENT_KEYS = [
      'ajs_user_id',
      'ajs_anonymous_id',
      'persisted-queue:v1:testkey:dest-Segment.io:items',
      'persisted-queue:v1:testkey:dest-Segment.io:seen',
    ];

    hooks.afterEach(function () {
      [...SEGMENT_KEYS, 'vault:prefs:telemetryConsent', 'vault:prefs:persona'].forEach((k) =>
        window.localStorage.removeItem(k)
      );
    });

    test('removes Segment identity and queue keys but preserves vault:prefs keys', function (assert) {
      SEGMENT_KEYS.forEach((k) => window.localStorage.setItem(k, 'x'));
      // Vault's own prefs (including the consent flag itself) must survive.
      window.localStorage.setItem('vault:prefs:telemetryConsent', 'false');
      window.localStorage.setItem('vault:prefs:persona', 'developer');

      SegmentProvider.purgeStorage();

      SEGMENT_KEYS.forEach((k) => assert.strictEqual(window.localStorage.getItem(k), null, `${k} is purged`));
      assert.strictEqual(
        window.localStorage.getItem('vault:prefs:telemetryConsent'),
        'false',
        'the consent flag is preserved'
      );
      assert.strictEqual(
        window.localStorage.getItem('vault:prefs:persona'),
        'developer',
        'unrelated vault prefs are preserved'
      );
    });

    test('reset() purges storage in addition to client.reset()', function (assert) {
      const provider = new SegmentProvider();
      const clientReset = sinon.stub(provider.client, 'reset');
      window.localStorage.setItem('ajs_user_id', 'IBMid-123');

      provider.reset();

      assert.true(clientReset.calledOnce, 'delegates to client.reset() for in-memory identity');
      assert.strictEqual(
        window.localStorage.getItem('ajs_user_id'),
        null,
        'also purges the persisted keys client.reset() leaves behind'
      );
    });
  });
});
