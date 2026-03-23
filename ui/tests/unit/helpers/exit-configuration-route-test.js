/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import exitConfigurationRoute from 'vault/helpers/exit-configuration-route';

module('Unit | Helper | exit-configuration-route', function () {
  test('alicloud returns list-root', function (assert) {
    const result = exitConfigurationRoute('alicloud');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('aws returns list-root', function (assert) {
    const result = exitConfigurationRoute('aws');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('azure returns backends route (isOnlyMountable)', function (assert) {
    const result = exitConfigurationRoute('azure');
    assert.strictEqual(result, 'vault.cluster.secrets.backends');
  });

  test('consul returns list-root', function (assert) {
    const result = exitConfigurationRoute('consul');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('cubbyhole returns list-root', function (assert) {
    const result = exitConfigurationRoute('cubbyhole');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('database returns list-root', function (assert) {
    const result = exitConfigurationRoute('database');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('gcp returns backends route (isOnlyMountable)', function (assert) {
    const result = exitConfigurationRoute('gcp');
    assert.strictEqual(result, 'vault.cluster.secrets.backends');
  });

  test('gcpkms returns list-root', function (assert) {
    const result = exitConfigurationRoute('gcpkms');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('kv with no version (defaults to v1) returns list-root', function (assert) {
    const result = exitConfigurationRoute('kv');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('kv v1 returns list-root', function (assert) {
    const result = exitConfigurationRoute('kv', 1);
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('kv v2 returns kv.list', function (assert) {
    const result = exitConfigurationRoute('kv', 2);
    assert.strictEqual(result, 'vault.cluster.secrets.backend.kv.list');
  });

  test('kmip returns kmip.scopes.index', function (assert) {
    const result = exitConfigurationRoute('kmip');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.kmip.scopes.index');
  });

  test('transform returns list-root', function (assert) {
    const result = exitConfigurationRoute('transform');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('keymgmt returns list-root', function (assert) {
    const result = exitConfigurationRoute('keymgmt');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('kubernetes returns kubernetes.overview', function (assert) {
    const result = exitConfigurationRoute('kubernetes');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.kubernetes.overview');
  });

  test('ldap returns ldap.overview', function (assert) {
    const result = exitConfigurationRoute('ldap');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.ldap.overview');
  });

  test('nomad returns list-root', function (assert) {
    const result = exitConfigurationRoute('nomad');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('pki returns pki.overview', function (assert) {
    const result = exitConfigurationRoute('pki');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.pki.overview');
  });

  test('rabbitmq returns list-root', function (assert) {
    const result = exitConfigurationRoute('rabbitmq');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('ssh returns list-root', function (assert) {
    const result = exitConfigurationRoute('ssh');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('totp returns list-root', function (assert) {
    const result = exitConfigurationRoute('totp');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('transit returns list-root', function (assert) {
    const result = exitConfigurationRoute('transit');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('unknown engine type returns list-root', function (assert) {
    const result = exitConfigurationRoute('unknown-engine');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });

  test('empty engine type returns list-root', function (assert) {
    const result = exitConfigurationRoute('');
    assert.strictEqual(result, 'vault.cluster.secrets.backend.list-root');
  });
});
