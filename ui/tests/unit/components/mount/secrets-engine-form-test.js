/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import SecretsEngineForm from 'vault/forms/secrets/engine';
import { getExternalPluginNameFromBuiltin } from 'vault/utils/external-plugin-helpers';

module('Unit | Component | mount/secrets-engine-form', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    // Setup default model
    const defaults = {
      config: { listing_visibility: false },
      options: { version: 2 },
    };
    this.model = new SecretsEngineForm(defaults, { isNew: true });
    this.model.type = 'keymgmt';
    this.model.data.path = 'keymgmt';

    this.availableVersions = [
      { version: '1.0.0', isBuiltin: false },
      { version: '1.1.0', isBuiltin: false },
      { version: '2.0.0', isBuiltin: false },
      { version: '', isBuiltin: true },
    ];
  });

  test('getExternalPluginNameFromBuiltin returns correct name for keymgmt', function (assert) {
    const externalName = getExternalPluginNameFromBuiltin('keymgmt');
    assert.strictEqual(
      externalName,
      'vault-plugin-secrets-keymgmt',
      'generates correct external plugin name for keymgmt'
    );
  });

  test('model normalizedType returns correct value', function (assert) {
    assert.strictEqual(this.model.normalizedType, 'keymgmt', 'returns correct normalized type');
  });
});
