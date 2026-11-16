/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import DatabaseConnectionForm, {
  getShowAttrs,
  isAvailablePlugin,
  normalizeDatabaseConnection,
} from 'vault/forms/database/connection';
import { EXPECTED_FIELDS } from 'vault/tests/helpers/secret-engine/database-helpers';
import { AVAILABLE_PLUGIN_TYPES } from 'vault/utils/model-helpers/database-helpers';

module('Unit | Form | database/connection', function () {
  for (const plugin of AVAILABLE_PLUGIN_TYPES.map((p) => p.value)) {
    const pluginHasTlsOptions = AVAILABLE_PLUGIN_TYPES.find((o) => o.value === plugin).fields.some(
      (f) => f.subgroup === 'TLS options'
    );
    module(`it computes fields for plugin_type: ${plugin}`, function (hooks) {
      hooks.beforeEach(function () {
        this.form = new DatabaseConnectionForm({ plugin_name: plugin });
        this.expected = EXPECTED_FIELDS[plugin];

        const pluginFields = AVAILABLE_PLUGIN_TYPES.find((o) => o.value === plugin).fields;
        this.enterpriseAttrs = pluginFields.filter((f) => f.isEnterprise).map((f) => f.attr);
      });

      /**
       * fieldAttrs are the fields rendered before the plugin config section and must match
       * the fields defined for the selected plugin.
       */
      test('it computes fieldAttrs', function (assert) {
        assert.propEqual(
          this.form.fieldAttrs.map((a) => a.name),
          this.expected.fieldAttrs,
          'actual computed attrs match expected'
        );
      });

      /**
       * The default group of pluginFieldGroups holds the plugin config fields that have no subgroup.
       */
      test('it computes default group', function (assert) {
        const [actualDefault] = this.form.pluginFieldGroups;
        assert.propEqual(
          actualDefault.default.map((a) => a.name),
          this.expected.default,
          'it has expected default group attributes'
        );
      });

      /**
       * statementFields are the fields for database statements, such as root rotation statements.
       */
      test('it computes statementFields', function (assert) {
        assert.propEqual(
          this.form.statementFields.map((a) => a.name),
          this.expected.statementFields.map((a) => a.name),
          'actual computed attrs match expected'
        );
      });

      if (pluginHasTlsOptions) {
        /**
         * Plugins with a TLS options subgroup return it as the second group of pluginFieldGroups.
         */
        test('it computes TLS options group', function (assert) {
          const [, actualTlsOptions] = this.form.pluginFieldGroups;
          assert.propEqual(
            actualTlsOptions['TLS options'].map((a) => a.name),
            this.expected['TLS options'],
            'it has expected TLS options'
          );
        });
      }

      /**
       * getShowAttrs drives the show view and must follow the selected plugin, ending with allowed_roles.
       */
      test('it computes show attrs', function (assert) {
        assert.propEqual(
          getShowAttrs(plugin, false).map((a) => a.name),
          this.expected.showAttrs,
          'actual computed attrs match expected'
        );
      });

      /**
       * Enterprise-only plugin fields must only be returned when the form is created for an enterprise cluster.
       */
      test('it gates enterprise fields', function (assert) {
        const names = (form) =>
          form.pluginFieldGroups.flatMap((group) => Object.values(group).flat()).map((a) => a.name);
        const community = names(this.form);
        const enterprise = names(new DatabaseConnectionForm({ plugin_name: plugin }, { isEnterprise: true }));
        for (const attr of this.enterpriseAttrs) {
          assert.false(community.includes(attr), `${attr} is hidden on community`);
          assert.true(enterprise.includes(attr), `${attr} is shown on enterprise`);
        }
        assert.strictEqual(
          enterprise.length - community.length,
          this.enterpriseAttrs.length,
          'only enterprise attrs differ between editions'
        );
      });
    });
  }

  /**
   * Before a plugin is selected only the generic fields are shown and there is no plugin config section.
   */
  test('it returns generic fields when no plugin is selected', function (assert) {
    const form = new DatabaseConnectionForm();
    assert.false(isAvailablePlugin(form.data.plugin_name), 'plugin is not available');
    assert.propEqual(
      form.fieldAttrs.map((a) => a.name),
      ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy'],
      'it returns the generic field attributes'
    );
    assert.strictEqual(form.pluginFieldGroups, null, 'there are no plugin field groups');
  });

  /**
   * The form proxy exposes data keys for every connection field so templates and
   * the plugin_name selection can read and write them directly.
   */
  test('it exposes connection data via the form proxy', function (assert) {
    const form = new DatabaseConnectionForm({ name: 'my-db', backend: 'database' });
    form.plugin_name = 'mysql-database-plugin';
    assert.strictEqual(form.name, 'my-db', 'it reads data');
    assert.strictEqual(form.data.plugin_name, 'mysql-database-plugin', 'it writes data');
    assert.true(isAvailablePlugin(form.data.plugin_name), 'plugin is available once selected');
  });

  /**
   * Only keys valid for the selected plugin are sent, undefined values are dropped and
   * name/backend (which are part of the url) are not part of the body.
   */
  test('toJSON only includes keys valid for the database type (elasticsearch)', function (assert) {
    const form = new DatabaseConnectionForm({
      plugin_name: 'elasticsearch-database-plugin',
      backend: 'db',
      name: 'elastic-test',
      verify_connection: true,
      allowed_roles: ['readonly'],
      connection_url: 'http://localhost:9200', // not valid for elasticsearch
      url: 'http://localhost:9200',
      username: 'elastic',
      password: 'changeme',
      skip_static_role_import_rotation: false,
      tls_ca: 'some-value', // not valid for elasticsearch
      ca_cert: undefined, // undefined values are not sent
      insecure: false,
    });
    assert.deepEqual(form.toJSON().data, {
      plugin_name: 'elasticsearch-database-plugin',
      verify_connection: true,
      allowed_roles: ['readonly'],
      url: 'http://localhost:9200',
      username: 'elastic',
      password: 'changeme',
      skip_static_role_import_rotation: false,
      insecure: false,
    });
    assert.strictEqual(form.data.name, 'elastic-test', 'toJSON does not mutate the form data');
  });

  /**
   * connection_details are flattened into the connection and the url is decoded so
   * the templated username/password characters are displayed and resubmitted correctly.
   */
  test('normalizeDatabaseConnection flattens the response (elasticsearch)', function (assert) {
    const connection = normalizeDatabaseConnection(
      {
        allowed_roles: ['readonly'],
        connection_details: {
          backend: 'database',
          insecure: false,
          url: 'https://localhost:9200',
          username: 'root',
        },
        password_policy: '',
        plugin_name: 'elasticsearch-database-plugin',
        plugin_version: '',
        root_credentials_rotate_statements: [],
      },
      'elastic-test',
      'database'
    );
    assert.propContains(connection, {
      allowed_roles: ['readonly'],
      backend: 'database',
      insecure: false,
      name: 'elastic-test',
      password_policy: '',
      plugin_name: 'elasticsearch-database-plugin',
      plugin_version: '',
      root_rotation_statements: [],
      url: 'https://localhost:9200',
      username: 'root',
    });
  });

  test('normalizeDatabaseConnection decodes connection_url (oracle)', function (assert) {
    const connection = normalizeDatabaseConnection(
      {
        allowed_roles: ['*'],
        connection_details: {
          connection_url: '%7B%7Busername%7D%7D/%7B%7Bpassword%7D%7D@//localhost:1521/ORCLPDB1',
          max_open_connections: 3,
          username: 'VAULTADMIN',
        },
        plugin_name: 'vault-plugin-database-oracle',
      },
      'oracle-test',
      'database'
    );
    assert.propContains(connection, {
      connection_url: '{{username}}/{{password}}@//localhost:1521/ORCLPDB1',
      max_open_connections: 3,
      name: 'oracle-test',
      username: 'VAULTADMIN',
    });
  });

  test('normalizeDatabaseConnection handles missing params and connection_details', function (assert) {
    const withoutUrl = normalizeDatabaseConnection(
      { allowed_roles: ['*'], connection_details: {}, plugin_name: 'vault-postgres-db' },
      'db-test',
      'database'
    );
    assert.strictEqual(withoutUrl.connection_url, undefined, 'connection_url is not set');
    const withoutData = normalizeDatabaseConnection(undefined, 'db-test', 'database');
    assert.propContains(withoutData, { name: 'db-test', backend: 'database' });
  });
});
