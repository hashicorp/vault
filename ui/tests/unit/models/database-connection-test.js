/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { EXPECTED_FIELDS } from 'vault/tests/helpers/secret-engine/database-helpers';

import { AVAILABLE_PLUGIN_TYPES } from 'vault/utils/model-helpers/database-helpers';

module('Unit | Model | database/connection', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    // setting version here so tests can be run locally on CE or ENT
    this.version.type = 'community';

    this.createModel = (plugin) =>
      this.store.createRecord('database/connection', {
        plugin_name: plugin,
      });
  });

  for (const plugin of AVAILABLE_PLUGIN_TYPES.map((p) => p.value)) {
    module(`it computes fields for plugin_type: ${plugin}`, function (hooks) {
      hooks.beforeEach(function () {
        this.model = this.createModel(plugin);
        this.getActual = (group) => this.model[group];
        this.getExpected = (group) => EXPECTED_FIELDS[plugin][group];

        const pluginFields = AVAILABLE_PLUGIN_TYPES.find((o) => o.value === plugin).fields;
        this.pluginHasTlsOptions = pluginFields.some((f) => f.subgroup === 'TLS options');
        this.pluginHasEnterpriseAttrs = pluginFields.some((f) => f.isEnterprise);
      });

      test('it computes showAttrs', function (assert) {
        const actual = this.getActual('showAttrs').map((a) => a.name);
        const expected = this.getExpected('showAttrs');
        assert.propEqual(actual, expected, 'actual computed attrs match expected');
      });

      test('it computes fieldAttrs', function (assert) {
        const actual = this.getActual('fieldAttrs').map((a) => a.name);
        const expected = this.getExpected('fieldAttrs');
        assert.propEqual(actual, expected, 'actual computed attrs match expected');
      });

      test('it computes default group', function (assert) {
        // pluginFieldGroups is an array of group objects
        const [actualDefault] = this.getActual('pluginFieldGroups');

        assert.propEqual(
          actualDefault.default.map((a) => a.name),
          this.getExpected('default'),
          'it has expected default group attributes'
        );
      });

      test('it computes statementFields', function (assert) {
        const actual = this.getActual('statementFields');
        const expected = this.getExpected('statementFields');
        assert.propEqual(actual, expected, 'actual computed attrs match expected');
      });

      if (this.pluginHasTlsOptions) {
        test('it computes TLS options group', function (assert) {
          // pluginFieldGroups is an array of group objects
          const [, actualTlsOptions] = this.getActual('pluginFieldGroups');

          assert.propEqual(
            actualTlsOptions['TLS options'].map((a) => a.name),
            this.getExpected('TLS options'),
            'it has expected TLS options'
          );
        });
      }

      if (this.pluginHasEnterpriseAttrs) {
        test('it includes enterprise fields', function (assert) {
          this.version.type = 'enterprise';
          const [actualDefault] = this.getActual('pluginFieldGroups');
          const expected = this.getExpected('default').push(this.getExpected('enterpriseOnly'));
          assert.propEqual(
            actualDefault.default.map((a) => a.name),
            expected,
            'it includes enterprise attributes'
          );
        });
      }
    });
  }
});
