/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import {
  PREFERENCES,
  getPreference,
  hasPreference,
  setPreference,
  getOrCreateAnalyticsUserId,
} from 'vault/utils/preferences';

module('Unit | Util | preferences', function (hooks) {
  hooks.beforeEach(function () {
    window.localStorage.clear();
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('registry keys follow the vault:prefs:<name> convention', function (assert) {
    Object.entries(PREFERENCES).forEach(([name, def]) => {
      assert.strictEqual(def.key, `vault:prefs:${name}`, `${name} uses the namespaced key convention`);
    });
  });

  test('telemetryConsent is registered and defaults off (opt-in)', function (assert) {
    assert.false(PREFERENCES.telemetryConsent.default, 'default is false');
    assert.strictEqual(PREFERENCES.telemetryConsent.type, 'boolean', 'type is boolean');
    assert.strictEqual(PREFERENCES.telemetryConsent.key, 'vault:prefs:telemetryConsent', 'key is namespaced');
  });

  test('getPreference returns the registry default when the key is absent', function (assert) {
    assert.false(getPreference('telemetryConsent'), 'returns documented default off');
  });

  test('getPreference throws for an unknown preference', function (assert) {
    assert.throws(() => getPreference('nope'), /Unknown preference "nope"/);
  });

  test('setPreference throws for an unknown preference', function (assert) {
    assert.throws(() => setPreference('nope', true), /Unknown preference "nope"/);
  });

  test('write/read round-trips a value through localStorage', function (assert) {
    setPreference('telemetryConsent', true);
    assert.strictEqual(
      window.localStorage.getItem('vault:prefs:telemetryConsent'),
      'true',
      'value is persisted under the namespaced key (JSON-serialized)'
    );
    assert.true(getPreference('telemetryConsent'), 'reads back the stored value');

    setPreference('telemetryConsent', false);
    assert.false(getPreference('telemetryConsent'), 'reads back an updated value');
  });

  module('getOrCreateAnalyticsUserId', function () {
    test('generates and persists a raw (unprefixed) id on first use', function (assert) {
      const id = getOrCreateAnalyticsUserId();

      assert.ok(id, 'returns an id');
      assert.notOk(id.startsWith('vault-'), 'the stored id is raw, without the vault- realm prefix');
      assert.strictEqual(
        window.localStorage.getItem('vault:prefs:analyticsUserId'),
        JSON.stringify(id),
        'persists the id under the namespaced key'
      );
    });

    test('returns the same id across calls (stable per browser)', function (assert) {
      const first = getOrCreateAnalyticsUserId();
      const second = getOrCreateAnalyticsUserId();

      assert.strictEqual(second, first, 'reuses the persisted id instead of generating a new one');
    });

    test('falls back to an in-memory id when localStorage is unavailable', function (assert) {
      sinon.stub(window.localStorage, 'getItem').throws(new Error('localStorage unavailable'));
      sinon.stub(window.localStorage, 'setItem').throws(new Error('localStorage unavailable'));

      const id = getOrCreateAnalyticsUserId();
      assert.ok(id, 'returns a uuid instead of throwing');
      assert.notOk(id.startsWith('vault-'), 'fallback id is also raw');
    });
  });

  module('fails safe when localStorage is unavailable', function () {
    test('getPreference returns the registered default when reads throw', function (assert) {
      sinon.stub(window.localStorage, 'getItem').throws(new Error('localStorage unavailable'));

      assert.false(getPreference('telemetryConsent'), 'falls back to the default instead of throwing');
    });

    test('hasPreference returns false when reads throw', function (assert) {
      sinon.stub(window.localStorage, 'getItem').throws(new Error('localStorage unavailable'));

      assert.false(hasPreference('telemetryConsent'), 'treats unreadable storage as "not stored"');
    });

    test('setPreference is a no-op when writes throw', function (assert) {
      sinon.stub(window.localStorage, 'setItem').throws(new Error('quota exceeded'));

      let error = null;
      try {
        setPreference('telemetryConsent', true);
      } catch (e) {
        error = e;
      }

      assert.strictEqual(error, null, 'does not throw when the write fails');
    });
  });
});
