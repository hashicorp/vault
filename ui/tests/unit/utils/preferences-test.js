/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { PREFERENCES, getPreference, setPreference } from 'vault/utils/preferences';

module('Unit | Util | preferences', function (hooks) {
  hooks.beforeEach(function () {
    window.localStorage.clear();
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
});
