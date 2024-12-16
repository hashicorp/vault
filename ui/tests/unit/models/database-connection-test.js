/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

import { AVAILABLE_PLUGIN_TYPES } from 'vault/utils/model-helpers/database-helpers';

module('Unit | Model | database/connection', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.version = this.owner.lookup('service:version');
    this.version.type = 'community';

    this.createModel = (plugin) =>
      this.store.createRecord('database/connection', {
        plugin_name: plugin,
      });
  });

  for (const plugin of AVAILABLE_PLUGIN_TYPES.map((p) => p.value)) {
    module(`it computes fields for plugin_type: ${plugin}`, function (hooks) {
      // If we refactor the database/connection model we can use this const in the updated form component
      // Ideal scenario be the API provides this schema so we don't have to manually write :)
      const EXPECTED_FIELDS = {
        'elasticsearch-database-plugin': {
          default: [
            'url',
            'username',
            'password',
            'ca_cert',
            'ca_path',
            'client_cert',
            'client_key',
            'tls_server_name',
            'insecure',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'url',
            'ca_cert',
            'ca_path',
            'client_cert',
            'client_key',
            'tls_server_name',
            'insecure',
            'username_template',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          pluginFieldGroups: undefined,
          statementFields: [],
        },
        'mongodb-database-plugin': {
          default: ['username', 'password', 'write_concern', 'username_template'],
          showAttrs: [
            'plugin_name',
            'name',
            'connection_url',
            'password_policy',
            'write_concern',
            'username_template',
            'tls',
            'tls_ca',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy'],
          pluginFieldGroups: ['username', 'password', 'write_concern', 'username_template'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
          'TLS options': ['tls', 'tls_ca'],
        },
        'mssql-database-plugin': {
          default: [
            'username',
            'password',
            'username_template',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'connection_url',
            'password_policy',
            'username_template',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
        },
        'mysql-aurora-database-plugin': {
          default: [
            'connection_url',
            'username',
            'password',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'connection_url',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
            'tls',
            'tls_ca',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
          'TLS options': ['tls', 'tls_ca'],
        },
        'mysql-legacy-database-plugin': {
          default: [
            'connection_url',
            'username',
            'password',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'connection_url',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
            'tls',
            'tls_ca',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
          'TLS options': ['tls', 'tls_ca'],
        },
        'mysql-database-plugin': {
          default: [
            'connection_url',
            'username',
            'password',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'connection_url',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
            'tls',
            'tls_ca',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
          'TLS options': ['tls', 'tls_ca'],
        },
        'mysql-rds-database-plugin': {
          default: [
            'connection_url',
            'username',
            'password',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'connection_url',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
            'tls',
            'tls_ca',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
          'TLS options': ['tls', 'tls_ca'],
        },
        'vault-plugin-database-oracle': {
          default: [
            'connection_url',
            'username',
            'password',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'connection_url',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
        },
        'postgresql-database-plugin': {
          default: [
            'connection_url',
            'username',
            'password',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
          ],
          showAttrs: [
            'plugin_name',
            'name',
            'password_policy',
            'connection_url',
            'max_open_connections',
            'max_idle_connections',
            'max_connection_lifetime',
            'username_template',
            'root_rotation_statements',
            'allowed_roles',
          ],
          fieldAttrs: ['plugin_name', 'name', 'verify_connection', 'password_policy'],
          statementFields: [
            {
              name: 'root_rotation_statements',
              options: {
                defaultShown: 'Default',
                editType: 'stringArray',
                subText:
                  "The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.",
              },
              type: undefined,
            },
          ],
        },
      };

      hooks.beforeEach(function () {
        this.model = this.createModel(plugin);
        this.getActual = (group) => this.model[group];
        this.getExpected = (group) => EXPECTED_FIELDS[plugin][group];

        this.pluginHasTlsOptions = AVAILABLE_PLUGIN_TYPES.find((o) => o.value === plugin).fields.some(
          (f) => f.subgroup === 'TLS options'
        );
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

      test('it computes statementFields', function (assert) {
        const actual = this.getActual('statementFields');
        const expected = this.getExpected('statementFields');
        assert.propEqual(actual, expected, 'actual computed attrs match expected');
      });
    });
  }
});
