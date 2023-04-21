/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Serializer | pki/crl', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.serializer = this.owner.lookup('serializer:pki/crl');
    this.store = this.owner.lookup('service:store');
  });

  test('it exists', function (assert) {
    assert.ok(this.serializer);
  });

  test('it serializes falsy ttl form values', async function (assert) {
    this.store.pushPayload('pki/crl', {
      modelName: 'pki/crl',
      id: 'some-pki-mount',
    });
    const record = this.store.peekRecord('pki/crl', 'some-pki-mount');
    // update record with values from ttl fields
    record.autoRebuildData = {
      enabled: false,
      duration: '24h',
    };
    record.deltaCrlBuildingData = {
      enabled: false,
      duration: '3d',
    };
    record.crlExpiryData = {
      enabled: false,
      duration: '24h',
    };
    record.ocspExpiryData = {
      enabled: false,
      duration: '15m',
    };
    const serializedRecord = record.serialize();
    assert.propEqual(
      serializedRecord,
      {
        auto_rebuild: false,
        auto_rebuild_grace_period: '24h',
        delta_rebuild_interval: '3d',
        disable: true,
        enable_delta: false,
        expiry: '24h',
        ocsp_disable: true,
        ocsp_expiry: '15m',
      },
      'it correctly transforms falsy values'
    );
  });

  test('it serializes api values to falsy (no form changes)', function (assert) {
    this.store.pushPayload('pki/crl', {
      modelName: 'pki/crl',
      id: 'some-pki-mount',
      auto_rebuild: false,
      auto_rebuild_grace_period: '12h',
      delta_rebuild_interval: '15m',
      disable: false,
      enable_delta: false,
      expiry: '72h',
      ocsp_disable: false,
      ocsp_expiry: '12h',
    });
    const record = this.store.peekRecord('pki/crl', 'some-pki-mount');
    const serializedRecord = record.serialize();
    assert.propEqual(
      serializedRecord,
      {
        auto_rebuild: false,
        auto_rebuild_grace_period: '12h',
        delta_rebuild_interval: '15m',
        disable: false,
        enable_delta: false,
        expiry: '72h',
        ocsp_disable: false,
        ocsp_expiry: '12h',
      },
      'it sends falsy values'
    );
  });

  test('it serializes truthy ttl form values', async function (assert) {
    this.store.pushPayload('pki/crl', {
      modelName: 'pki/crl',
      id: 'some-pki-mount',
    });
    const record = this.store.peekRecord('pki/crl', 'some-pki-mount');
    record.autoRebuildData = {
      enabled: true,
      duration: '12h',
    };
    record.deltaCrlBuildingData = {
      enabled: true,
      duration: '12h',
    };
    record.crlExpiryData = {
      enabled: true,
      duration: '12h',
    };
    record.ocspExpiryData = {
      enabled: true,
      duration: '12h',
    };

    const serializedRecord = record.serialize();
    assert.propEqual(
      serializedRecord,
      {
        auto_rebuild: true,
        auto_rebuild_grace_period: '12h',
        delta_rebuild_interval: '12h',
        disable: false,
        enable_delta: true,
        expiry: '12h',
        ocsp_disable: false,
        ocsp_expiry: '12h',
      },
      'it correctly transforms truthy values'
    );
  });

  test('it serializes api values to truthy (no form changes)', function (assert) {
    this.store.pushPayload('pki/crl', {
      modelName: 'pki/crl',
      id: 'some-pki-mount',
      auto_rebuild: true,
      auto_rebuild_grace_period: '3d',
      delta_rebuild_interval: '24m',
      disable: true,
      enable_delta: true,
      expiry: '34h',
      ocsp_disable: true,
      ocsp_expiry: '24h',
    });
    const record = this.store.peekRecord('pki/crl', 'some-pki-mount');
    const serializedRecord = record.serialize();

    assert.propEqual(
      serializedRecord,
      {
        auto_rebuild: true,
        auto_rebuild_grace_period: '3d',
        delta_rebuild_interval: '24m',
        disable: true,
        enable_delta: true,
        expiry: '34h',
        ocsp_disable: true,
        ocsp_expiry: '24h',
      },
      'it sends truthy values'
    );
  });

  test('it normalizes response when payload booleans are all false', function (assert) {
    const payload = {
      auto_rebuild: false,
      auto_rebuild_grace_period: '12h',
      cross_cluster_revocation: false,
      delta_rebuild_interval: '15m',
      disable: false,
      enable_delta: false,
      expiry: '72h',
      ocsp_disable: false,
      ocsp_expiry: '12h',
      unified_crl: false,
      unified_crl_on_existing_paths: false,
    };
    const normalizedRecord = this.serializer.normalizeResponse(
      this.store,
      this.store.modelFor('pki/crl'),
      payload,
      'some-pki-mount',
      'findRecord'
    );
    assert.propEqual(
      normalizedRecord,
      {
        data: {
          attributes: {
            autoRebuildGracePeriod: '12h',
            autoRebuild: false,
            autoRebuildData: {
              duration: '12h',
              enabled: false,
            },
            expiry: '72h',
            disable: false,
            crlExpiryData: {
              duration: '72h',
              enabled: true,
            },
            deltaRebuildInterval: '15m',
            enableDelta: false,
            deltaCrlBuildingData: {
              duration: '15m',
              enabled: false,
            },
            ocspExpiry: '12h',
            ocspDisable: false,
            ocspExpiryData: {
              duration: '12h',
              enabled: true,
            },
          },
          id: 'some-pki-mount',
          relationships: {},
          type: 'pki/crl',
        },
        included: [],
      },
      'returned payload has correct keys and values (from falsy booleans)'
    );
  });

  test('it normalizes response when payload booleans are all true', function (assert) {
    const payload = {
      auto_rebuild: true,
      auto_rebuild_grace_period: '12h',
      cross_cluster_revocation: true,
      delta_rebuild_interval: '15m',
      disable: true,
      enable_delta: true,
      expiry: '72h',
      ocsp_disable: true,
      ocsp_expiry: '12h',
      unified_crl: true,
      unified_crl_on_existing_paths: true,
    };
    const normalizedRecord = this.serializer.normalizeResponse(
      this.store,
      this.store.modelFor('pki/crl'),
      payload,
      'some-pki-mount',
      'findRecord'
    );
    assert.propEqual(
      normalizedRecord,
      {
        data: {
          attributes: {
            autoRebuild: true,
            autoRebuildData: {
              duration: '12h',
              enabled: true,
            },
            autoRebuildGracePeriod: '12h',
            crlExpiryData: {
              duration: '72h',
              enabled: false,
            },
            deltaCrlBuildingData: {
              duration: '15m',
              enabled: true,
            },
            ocspExpiryData: {
              duration: '12h',
              enabled: false,
            },
            deltaRebuildInterval: '15m',
            disable: true,
            enableDelta: true,
            expiry: '72h',
            ocspDisable: true,
            ocspExpiry: '12h',
          },
          id: 'some-pki-mount',
          relationships: {},
          type: 'pki/crl',
        },
        included: [],
      },
      'returned payload has correct keys and values (from truthy booleans)'
    );
  });
});
