/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { prepFormConfig, generateConfigContent } from 'vault/utils/form-config-generator';
import { OAS_STUB as SPEC } from 'vault/tests/helpers/stubs';

module('Unit | Utility | form-config-generator', function ({ beforeEach }) {
  beforeEach(function () {
    // This is essentially what happens in the `main()` function
    // within the CLI script generate-form-configs.ts
    this.config = prepFormConfig(SPEC, 'mountsEnableSecretsEngine');
    this.result = generateConfigContent(this.config);
  });

  module('#prepFormConfig', function () {
    test('returns null for unknown operation', function (assert) {
      const result = prepFormConfig(SPEC, 'unknownOperation');
      assert.strictEqual(result, null, 'returns null for unknown operation');
    });

    test('returns `name`, `apiClass`, and `requestType` properties', function (assert) {
      assert.propContains(
        this.config,
        {
          name: 'mountsEnableSecretsEngine',
          apiClass: 'sys',
          requestType: 'SystemApiMountsEnableSecretsEngineOperationRequest',
        },
        'config contains expected name, apiClass, and requestType'
      );
    });

    module('`payload` object', function () {
      test('includes expected default values', function (assert) {
        assert.deepEqual(
          this.config.payload,
          {
            path: '',
            MountsEnableSecretsEngineRequest: {
              config: {},
              description: '',
              local: false,
              options: {},
              plugin_name: '',
              plugin_version: '',
              seal_wrap: false,
              type: '',
              allowed_managed_keys: [],
            },
          },
          'payload matches expected structure and defaults'
        );
      });

      test('does not include deprecated properties', function (assert) {
        const { properties } = SPEC.components.schemas.MountsEnableSecretsEngineRequest;
        const deprecatedProperties = Object.keys(properties).filter((p) => properties[p].deprecated);
        const payloadKeys = Object.keys(this.config.payload.MountsEnableSecretsEngineRequest);

        // Double-check that the fixture actually includes deprecated properties
        assert.true(deprecatedProperties.length > 0, 'fixture includes deprecated properties');

        assert.notPropContains(
          payloadKeys,
          deprecatedProperties,
          'payload does not include deprecated properties'
        );
      });
    });

    module('`sections` array', function ({ beforeEach }) {
      beforeEach(function () {
        this.expectedSections = [
          {
            name: 'params',
            fields: [
              {
                name: 'path',
                type: 'TextInput',
                label: 'Path',
                helperText: 'The path to mount to. Example: "aws/east"',
              },
            ],
          },
          {
            name: 'default',
            fields: [
              {
                name: 'MountsEnableSecretsEngineRequest.config',
                type: 'TextInput',
                label: 'Config',
                helperText: 'Configuration for this mount, such as default_lease_ttl and max_lease_ttl.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.description',
                type: 'TextInput',
                label: 'Description',
                helperText: 'User-friendly description for this mount.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.local',
                type: 'TextInput',
                label: 'Local',
                helperText:
                  'Mark the mount as a local mount, which is not replicated and is unaffected by replication.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.options',
                type: 'TextInput',
                label: 'Options',
                helperText:
                  'The options to pass into the backend. Should be a json object with string keys and values.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.plugin_name',
                type: 'TextInput',
                label: 'Plugin name',
                helperText:
                  'Name of the plugin to mount based from the name registered in the plugin catalog.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.plugin_version',
                type: 'TextInput',
                label: 'Plugin version',
                helperText:
                  'The semantic version of the plugin to use, or image tag if oci_image is provided.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.type',
                type: 'TextInput',
                label: 'Type',
                helperText: 'The type of the backend. Example: "passthrough"',
              },
            ],
          },
          {
            name: 'Advanced',
            fields: [
              {
                name: 'MountsEnableSecretsEngineRequest.seal_wrap',
                type: 'TextInput',
                label: 'Seal Wrap',
                helperText: 'Whether to turn on seal wrapping for the mount.',
              },
              {
                name: 'MountsEnableSecretsEngineRequest.allowed_managed_keys',
                type: 'TextInput',
                label: 'Allowed Managed Keys',
                helperText: 'List of managed key names allowed for this mount.',
              },
            ],
          },
        ];
      });

      test('returns with expected fields', function (assert) {
        assert.deepEqual(
          this.config.sections,
          this.expectedSections,
          'sections match expected structure and fields'
        );
      });

      test('excludes deprecated properties from sections', function (assert) {
        const { properties } = SPEC.components.schemas.MountsEnableSecretsEngineRequest;
        const deprecatedProperties = Object.keys(properties).filter((p) => properties[p].deprecated);
        const allFieldNames = this.config.sections.flatMap((s) => s.fields.map((f) => f.name));

        assert.notPropContains(
          allFieldNames,
          deprecatedProperties,
          'sections do not include deprecated properties'
        );
      });

      test('is grouped by x-vault-displayAttrs group', function (assert) {
        const { properties } = SPEC.components.schemas.MountsEnableSecretsEngineRequest;

        // Create a hashmap of expected groups and their field names based on the spec
        // eg: { default: ['config', 'description', ...], Advanced: ['seal_wrap', ...] }
        const expectedGroups = Object.keys(properties).reduce((acc, prop) => {
          if (properties[prop].deprecated) {
            return acc;
          }
          const displayAttrs = properties[prop]['x-vault-displayAttrs'];
          const group = displayAttrs?.group || 'default';
          if (!acc[group]) {
            acc[group] = [];
          }
          acc[group].push(prop);
          return acc;
        }, {});

        // Sections for path is less relevant...
        // eslint-disable-next-line no-unused-vars, @typescript-eslint/no-unused-vars
        const [_pathSection, ...propSections] = this.config.sections;

        propSections.forEach((section) => {
          const expectedFieldNames = expectedGroups[section.name];
          const currentFieldNames = section.fields.map((f) =>
            // Remove the prefix to match property names as field names
            // as the prefix is utilized for mapping to the payload structure
            f.name.replace('MountsEnableSecretsEngineRequest.', '')
          );

          assert.deepEqual(
            currentFieldNames,
            expectedFieldNames,
            `${section.name} section contains correct fields from spec`
          );
        });
      });

      test('uses x-vault-displayAttrs.name as label when available', function (assert) {
        const { properties } = SPEC.components.schemas.MountsEnableSecretsEngineRequest;

        // Capture a list of fields that have x-vault-displayAttrs.name defined
        const fieldsWithDisplayNames = Object.keys(properties).reduce((acc, prop) => {
          const displayAttrs = properties[prop]['x-vault-displayAttrs'];
          if (displayAttrs && displayAttrs.name && !properties[prop].deprecated) {
            acc.push({ key: prop, label: displayAttrs.name });
          }
          return acc;
        }, []);

        fieldsWithDisplayNames.forEach(({ key, label }) => {
          const field = this.config.sections
            .flatMap((s) => s.fields)
            .find((field) => field.name.endsWith(key));

          assert.strictEqual(field?.label, label, `${key} uses x-vault-displayAttrs.name as label`);
        });
      });

      test('falls back to sentence case for label when x-vault-displayAttrs.name is missing', function (assert) {
        const defaultSection = this.config.sections.find((s) => s.name === 'default');
        const pluginNameField = defaultSection.fields.find((f) => f.name.includes('plugin_name'));

        assert.strictEqual(pluginNameField.label, 'Plugin name', 'falls back to sentence case for label');
      });
    });
  });

  module('#generateConfigContent', function () {
    test('produces content with correct imports', function (assert) {
      assert.true(
        this.result.includes("import type ApiService from 'vault/services/api'"),
        'includes ApiService import'
      );
      assert.true(
        this.result.includes(`import type { ${this.config.requestType} }`),
        'includes request type import'
      );
    });

    test('produces content with correct name property', function (assert) {
      assert.true(this.result.includes(`name: '${this.config.name}'`), 'includes correct name property');
    });

    test('produces content with correct submit implementation', function (assert) {
      const expectedSubmit = `submit: async (api: ApiService, payload: ${this.config.requestType}) => {
        return await api.${this.config.apiClass}.${this.config.name}Raw(payload);
      },`;

      assert.true(
        this.result.includes(expectedSubmit),
        'submit property has correct async function implementation'
      );
    });

    test('produces content with correct payload', function (assert) {
      assert.true(
        this.result.includes(JSON.stringify(this.config.payload, null, 2)),
        'payload matches config payload'
      );
    });

    test('produces content with correct sections', function (assert) {
      assert.true(
        this.result.includes(JSON.stringify(this.config.sections, null, 2)),
        'sections match config sections'
      );
    });

    test('produces content with correct export statement', function (assert) {
      assert.true(
        this.result.includes(`export default ${this.config.name}Config`),
        'exports config with correct name'
      );
    });
  });
});
