/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import MfaMethodController from 'vault/controllers/vault/cluster/access/mfa/methods/method';

// ---------------------------------------------------------------------------
// Fixture helpers
// All secret values use the TEST_ONLY_ prefix — they are clearly
// non-production strings and are never committed as real vendor keys.
// ---------------------------------------------------------------------------

function makeDuoRecord(store) {
  return store.createRecord('mfa-method', {
    type: 'duo',
    secret_key: 'TEST_ONLY_DUO_SECRET_KEY',
    integration_key: 'TEST_ONLY_DUO_INTEGRATION_KEY',
    api_hostname: 'api.duosecurity.example.com',
    push_info: 'test push info',
    use_passcode: false,
    username_format: '%s@example.com',
  });
}

function makeOktaRecord(store) {
  return store.createRecord('mfa-method', {
    type: 'okta',
    api_token: 'TEST_ONLY_OKTA_API_TOKEN',
    org_name: 'example-org',
    username_format: '%s',
    base_url: 'okta.com',
    primary_email: true,
  });
}

function makePingIdRecord(store) {
  return store.createRecord('mfa-method', {
    type: 'pingid',
    settings_file_base64: 'TEST_ONLY_PINGID_SETTINGS',
    use_signature: true,
    idp_url: 'https://idp.example.pingidentity.com',
    admin_url: 'https://admin.example.pingidentity.com',
    authenticator_url: 'https://authenticator.pingidentity.example.com',
    org_alias: 'test-org',
    username_format: '%s',
  });
}

function makeTotpRecord(store) {
  return store.createRecord('mfa-method', {
    type: 'totp',
    issuer: 'Vault',
    period: 30,
    key_size: 20,
    qr_size: 200,
    algorithm: 'SHA1',
    digits: 6,
    skew: 1,
    max_validation_attempts: 5,
    enable_self_enrollment: false,
  });
}

function makeController(owner, method) {
  // Register with singleton:false so each call creates a fresh instance and
  // tests do not share state. Re-registering is idempotent for the test run.
  owner.register('controller:vault/cluster/access/mfa/methods/method', MfaMethodController, {
    singleton: false,
  });
  const controller = owner.lookup('controller:vault/cluster/access/mfa/methods/method');
  controller.model = { method };
  return controller;
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

module('Unit | Controller | mfa method detail read model', function (hooks) {
  setupTest(hooks);

  // -------------------------------------------------------------------------
  // AC1 — displayFields for Duo must not include secret_key, integration_key,
  //        or api_token. Same constraint for Okta.
  // -------------------------------------------------------------------------

  test('Duo displayFields does not include secret_key, integration_key, or api_token', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeDuoRecord(store));

    const fields = controller.displayFields;
    assert.false(fields.includes('secret_key'), 'secret_key is absent from Duo displayFields');
    assert.false(fields.includes('integration_key'), 'integration_key is absent from Duo displayFields');
    assert.false(fields.includes('api_token'), 'api_token is absent from Duo displayFields');
  });

  test('Okta displayFields does not include secret_key, integration_key, or api_token', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeOktaRecord(store));

    const fields = controller.displayFields;
    assert.false(fields.includes('secret_key'), 'secret_key is absent from Okta displayFields');
    assert.false(fields.includes('integration_key'), 'integration_key is absent from Okta displayFields');
    assert.false(fields.includes('api_token'), 'api_token is absent from Okta displayFields');
  });

  // -------------------------------------------------------------------------
  // AC2 — PingID displayFields must not include settings_file_base64, and
  //        JSON.stringify of the sanitized model must not contain any of
  //        the fixture secret strings.
  // -------------------------------------------------------------------------

  test('PingID displayFields does not include settings_file_base64', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makePingIdRecord(store));

    assert.false(
      controller.displayFields.includes('settings_file_base64'),
      'settings_file_base64 is absent from PingID displayFields'
    );
  });

  test('sanitized model JSON does not contain Duo secret strings', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeDuoRecord(store));

    const json = JSON.stringify(controller.sanitizedModel);
    assert.false(
      json.includes('TEST_ONLY_DUO_SECRET_KEY'),
      'Duo secret_key value is absent from sanitized model JSON'
    );
    assert.false(
      json.includes('TEST_ONLY_DUO_INTEGRATION_KEY'),
      'Duo integration_key value is absent from sanitized model JSON'
    );
  });

  test('sanitized model JSON does not contain Okta api_token string', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeOktaRecord(store));

    const json = JSON.stringify(controller.sanitizedModel);
    assert.false(
      json.includes('TEST_ONLY_OKTA_API_TOKEN'),
      'Okta api_token value is absent from sanitized model JSON'
    );
  });

  test('sanitized model JSON does not contain PingID settings string', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makePingIdRecord(store));

    const json = JSON.stringify(controller.sanitizedModel);
    assert.false(
      json.includes('TEST_ONLY_PINGID_SETTINGS'),
      'PingID settings_file_base64 value is absent from sanitized model JSON'
    );
  });

  // -------------------------------------------------------------------------
  // AC3 — TOTP displayFields retains all required non-secret fields;
  //        Duo retains api_hostname.
  // -------------------------------------------------------------------------

  test('TOTP displayFields retains issuer, period, key_size, qr_size, algorithm, digits, skew, max_validation_attempts', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeTotpRecord(store));

    const fields = controller.displayFields;
    [
      'issuer',
      'period',
      'key_size',
      'qr_size',
      'algorithm',
      'digits',
      'skew',
      'max_validation_attempts',
    ].forEach((field) => {
      assert.true(fields.includes(field), `TOTP displayFields includes ${field}`);
    });
  });

  test('Duo displayFields retains api_hostname', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeDuoRecord(store));

    assert.true(controller.displayFields.includes('api_hostname'), 'Duo displayFields retains api_hostname');
  });

  // -------------------------------------------------------------------------
  // AC4 — Source ember-data record is not mutated; secret_key is still
  //        present on the write model after sanitizedModel is computed.
  // -------------------------------------------------------------------------

  test('source ember-data record still has secret_key after sanitizedModel is accessed', function (assert) {
    const store = this.owner.lookup('service:store');
    const method = makeDuoRecord(store);
    const controller = makeController(this.owner, method);

    // Trigger the sanitized model computation
    const sanitized = controller.sanitizedModel;

    assert.strictEqual(
      method.secret_key,
      'TEST_ONLY_DUO_SECRET_KEY',
      'source record secret_key is not cleared by sanitizedModel'
    );
    assert.ok(sanitized, 'sanitized model was computed');
  });

  // -------------------------------------------------------------------------
  // AC5 (indicator rows) — writeOnlyIndicators rows expose only label and
  //      configured boolean; no secret property and no secret string value.
  // -------------------------------------------------------------------------

  test('writeOnlyIndicators for Duo contain only label and configured, no secret property', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeDuoRecord(store));

    const indicators = controller.writeOnlyIndicators;
    assert.true(indicators.length > 0, 'indicator rows are present for Duo');

    indicators.forEach((row, i) => {
      assert.strictEqual(typeof row.label, 'string', `row[${i}] has a string label`);
      assert.strictEqual(typeof row.configured, 'boolean', `row[${i}] has a boolean configured`);
      assert.false('secret' in row, `row[${i}] has no secret property`);
      const rowValues = Object.values(row);
      assert.false(
        rowValues.some((v) => typeof v === 'string' && v.includes('TEST_ONLY_')),
        `row[${i}] values do not contain any TEST_ONLY_ secret string`
      );
    });
  });

  test('writeOnlyIndicators for Okta marks api_token as configured when present', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeOktaRecord(store));

    const indicators = controller.writeOnlyIndicators;
    assert.strictEqual(indicators.length, 1, 'one indicator row for Okta');
    assert.true(indicators[0].configured, 'api_token is marked configured');
    assert.strictEqual(typeof indicators[0].label, 'string', 'indicator has a string label');
  });

  test('writeOnlyIndicators for TOTP is empty', function (assert) {
    const store = this.owner.lookup('service:store');
    const controller = makeController(this.owner, makeTotpRecord(store));

    assert.deepEqual(controller.writeOnlyIndicators, [], 'TOTP has no write-only indicator rows');
  });

  // -------------------------------------------------------------------------
  // AC6 — Mock fixtures use TEST_ONLY_ prefix; no real vendor keys committed.
  // -------------------------------------------------------------------------

  test('fixture secret values use TEST_ONLY_ prefix', function (assert) {
    const store = this.owner.lookup('service:store');
    const duo = makeDuoRecord(store);
    const okta = makeOktaRecord(store);
    const pingid = makePingIdRecord(store);

    assert.true(duo.secret_key.startsWith('TEST_ONLY_'), 'Duo secret_key uses TEST_ONLY_ prefix');
    assert.true(duo.integration_key.startsWith('TEST_ONLY_'), 'Duo integration_key uses TEST_ONLY_ prefix');
    assert.true(okta.api_token.startsWith('TEST_ONLY_'), 'Okta api_token uses TEST_ONLY_ prefix');
    assert.true(
      pingid.settings_file_base64.startsWith('TEST_ONLY_'),
      'PingID settings_file_base64 uses TEST_ONLY_ prefix'
    );
  });

  // -------------------------------------------------------------------------
  // Edge cases
  // -------------------------------------------------------------------------

  test('TOTP sanitized model does not expose secret_key even when source unexpectedly carries it', function (assert) {
    const store = this.owner.lookup('service:store');
    // Simulate an unexpected API payload that echoes secret_key on a TOTP record
    const method = store.createRecord('mfa-method', {
      type: 'totp',
      issuer: 'Vault',
      secret_key: 'TEST_ONLY_DUO_SECRET_KEY',
      settings_file_base64: 'TEST_ONLY_PINGID_SETTINGS',
    });
    const controller = makeController(this.owner, method);

    assert.false(
      controller.displayFields.includes('secret_key'),
      'secret_key absent from TOTP displayFields'
    );
    assert.false(
      controller.displayFields.includes('settings_file_base64'),
      'settings_file_base64 absent from TOTP displayFields'
    );

    const json = JSON.stringify(controller.sanitizedModel);
    assert.false(
      json.includes('TEST_ONLY_DUO_SECRET_KEY'),
      'secret_key string is absent from TOTP sanitized model'
    );
    assert.false(
      json.includes('TEST_ONLY_PINGID_SETTINGS'),
      'settings_file_base64 string is absent from TOTP sanitized model'
    );

    assert.deepEqual(
      controller.writeOnlyIndicators,
      [],
      'TOTP has no write-only indicator rows even with unexpected secret fields'
    );
  });

  test('Duo with only some write-only fields present marks each correctly', function (assert) {
    const store = this.owner.lookup('service:store');
    const method = store.createRecord('mfa-method', {
      type: 'duo',
      secret_key: 'TEST_ONLY_DUO_SECRET_KEY',
      // integration_key is intentionally absent (null/undefined)
      api_hostname: 'api.duosecurity.example.com',
    });
    const controller = makeController(this.owner, method);

    const indicators = controller.writeOnlyIndicators;
    const secretKeyRow = indicators.find((r) => r.label === controller.label('secret_key'));
    const integrationKeyRow = indicators.find((r) => r.label === controller.label('integration_key'));

    assert.ok(secretKeyRow, 'indicator row for secret_key exists');
    assert.ok(integrationKeyRow, 'indicator row for integration_key exists');
    assert.true(secretKeyRow.configured, 'secret_key is marked configured (value present)');
    assert.false(integrationKeyRow.configured, 'integration_key is marked not configured (value absent)');
  });

  test('whitespace-only credential is treated as not configured', function (assert) {
    const store = this.owner.lookup('service:store');
    const method = store.createRecord('mfa-method', {
      type: 'duo',
      secret_key: '   ', // whitespace-only — must not be configured
      integration_key: 'TEST_ONLY_DUO_INTEGRATION_KEY',
      api_hostname: 'api.duosecurity.example.com',
    });
    const controller = makeController(this.owner, method);

    const indicators = controller.writeOnlyIndicators;
    const secretKeyRow = indicators.find((r) => r.label === controller.label('secret_key'));
    assert.false(secretKeyRow.configured, 'whitespace-only secret_key is treated as not configured');
  });

  test('large PingID base64 blob is fully absent from sanitized model — not truncated', function (assert) {
    const largePingIdSettings = 'TEST_ONLY_PINGID_SETTINGS_' + 'A'.repeat(5000);
    const store = this.owner.lookup('service:store');
    const method = store.createRecord('mfa-method', {
      type: 'pingid',
      settings_file_base64: largePingIdSettings,
      idp_url: 'https://idp.example.pingidentity.com',
      admin_url: 'https://admin.example.pingidentity.com',
      authenticator_url: 'https://auth.example.pingidentity.com',
      org_alias: 'test-org',
    });
    const controller = makeController(this.owner, method);

    const json = JSON.stringify(controller.sanitizedModel);
    assert.false(
      json.includes('TEST_ONLY_PINGID_SETTINGS'),
      'large settings_file_base64 blob is fully absent from PingID sanitized model'
    );
    const pingIdRow = controller.writeOnlyIndicators.find(
      (r) => r.label === controller.label('settings_file_base64')
    );
    assert.true(pingIdRow.configured, 'large PingID blob is marked configured');
  });

  test('null model method returns safe defaults for sanitizedModel and writeOnlyIndicators', function (assert) {
    this.owner.register('controller:vault/cluster/access/mfa/methods/method', MfaMethodController, {
      singleton: false,
    });
    const controller = this.owner.lookup('controller:vault/cluster/access/mfa/methods/method');
    controller.model = { method: null };

    assert.strictEqual(controller.sanitizedModel, null, 'sanitizedModel is null when method is null');
    assert.deepEqual(
      controller.writeOnlyIndicators,
      [],
      'writeOnlyIndicators is empty array when method is null'
    );
  });
});
