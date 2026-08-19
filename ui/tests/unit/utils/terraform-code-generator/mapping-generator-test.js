/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  resourceArgLine,
  crossReference,
  featureKeyHint,
  generateInteractiveScaffold,
  generateOpenApiScaffold,
} from 'vault/utils/terraform-code-generator/mapping-generator';

module('Unit | Utility | terraform-code-generator/mapping-generator', function () {
  // ---------------------------------------------------------------------------
  // resourceArgLine
  // ---------------------------------------------------------------------------

  module('#resourceArgLine', function () {
    test('wraps string fields in template-literal quotes', function (assert) {
      assert.strictEqual(
        resourceArgLine({ name: 'path', type: 'string' }),
        '      path: `"${payload.path}"`,'
      );
    });

    test('uses formatEot for heredoc fields', function (assert) {
      assert.strictEqual(
        resourceArgLine({ name: 'policy', type: 'heredoc' }),
        '      policy: formatEot(payload.policy),'
      );
    });

    test('passes boolean fields through without quoting', function (assert) {
      assert.strictEqual(
        resourceArgLine({ name: 'disable_remount', type: 'boolean' }),
        '      disable_remount: payload.disable_remount,'
      );
    });

    test('passes number fields through without quoting', function (assert) {
      assert.strictEqual(
        resourceArgLine({ name: 'max_ttl', type: 'number' }),
        '      max_ttl: payload.max_ttl,'
      );
    });

    test('emits a TODO comment for object fields', function (assert) {
      assert.true(resourceArgLine({ name: 'options', type: 'object' }).includes('// TODO: options'));
    });

    test('emits a TODO comment for array fields', function (assert) {
      assert.true(
        resourceArgLine({ name: 'token_policies', type: 'array' }).includes('// TODO: token_policies')
      );
    });

    test('prefers fieldType over type when both are present (OpenAPI mode)', function (assert) {
      assert.strictEqual(
        resourceArgLine({ name: 'policy', type: 'string', fieldType: 'heredoc' }),
        '      policy: formatEot(payload.policy),'
      );
    });
  });

  // ---------------------------------------------------------------------------
  // crossReference
  // ---------------------------------------------------------------------------

  module('#crossReference', function () {
    test('returns only fields present in both OpenAPI and Terraform', function (assert) {
      const openApiFields = [
        { name: 'name', openApiType: 'string' },
        { name: 'policy', openApiType: 'string' },
        { name: 'only_in_api', openApiType: 'string' },
      ];
      const tfAttributes = { name: {}, policy: {}, only_in_tf: {} };
      const { matched, inOpenApiOnly, inTfOnly } = crossReference(openApiFields, tfAttributes);
      assert.deepEqual(
        matched.map((f) => f.name),
        ['name', 'policy']
      );
      assert.deepEqual(
        inOpenApiOnly.map((f) => f.name),
        ['only_in_api']
      );
      assert.deepEqual(inTfOnly, ['only_in_tf']);
    });

    test('excludes id, accessor, and namespace from Terraform fields', function (assert) {
      const openApiFields = [{ name: 'name', openApiType: 'string' }];
      const tfAttributes = { name: {}, id: {}, accessor: {}, namespace: {} };
      const { matched, inTfOnly } = crossReference(openApiFields, tfAttributes);
      assert.deepEqual(
        matched.map((f) => f.name),
        ['name']
      );
      assert.deepEqual(inTfOnly, []);
    });

    test('assigns fieldType: heredoc for known heredoc field names', function (assert) {
      const openApiFields = [{ name: 'policy', openApiType: 'string' }];
      const { matched } = crossReference(openApiFields, { policy: {} });
      assert.strictEqual(matched[0].fieldType, 'heredoc');
    });

    test('assigns fieldType: boolean for boolean openApiType', function (assert) {
      const openApiFields = [{ name: 'disable_remount', openApiType: 'boolean' }];
      const { matched } = crossReference(openApiFields, { disable_remount: {} });
      assert.strictEqual(matched[0].fieldType, 'boolean');
    });

    test('assigns fieldType: object for object openApiType', function (assert) {
      const openApiFields = [{ name: 'options', openApiType: 'object' }];
      const { matched } = crossReference(openApiFields, { options: {} });
      assert.strictEqual(matched[0].fieldType, 'object');
    });
  });

  // ---------------------------------------------------------------------------
  // featureKeyHint
  // ---------------------------------------------------------------------------

  module('#featureKeyHint', function () {
    test('strips /sys/ prefix and path params', function (assert) {
      assert.strictEqual(featureKeyHint('/sys/policies/acl/{name}'), 'policies/acl');
    });

    test('strips multiple path params', function (assert) {
      assert.strictEqual(featureKeyHint('/sys/auth/{path}/tune'), 'auth/tune');
    });

    test('leaves non-sys paths intact apart from path params', function (assert) {
      assert.strictEqual(featureKeyHint('/auth/{path}/login'), '/auth/login');
    });
  });

  // ---------------------------------------------------------------------------
  // generateInteractiveScaffold
  // ---------------------------------------------------------------------------

  module('#generateInteractiveScaffold', function () {
    const base = {
      method: 'mountsEnableSecretsEngine',
      tfResource: 'vault_mount',
      featureKey: 'secrets/kv',
      fields: [],
    };

    test('produces correct interface and mapping names from method', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [{ name: 'path', type: 'string', required: true }],
      });
      assert.true(result.includes('export interface MountsEnableSecretsEnginePayload'));
      assert.true(result.includes('export const mountsEnableSecretsEngineMapping'));
    });

    test('marks optional fields with ?', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [
          { name: 'path', type: 'string', required: true },
          { name: 'description', type: 'string', required: false },
        ],
      });
      assert.true(result.includes('path: string;'));
      assert.true(result.includes('description?: string;'));
    });

    test('includes formatEot import for heredoc fields', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [{ name: 'policy', type: 'heredoc', required: true }],
      });
      assert.true(result.includes("import { formatEot } from 'core/utils/code-generators/formatters'"));
      assert.true(result.includes('formatEot(payload.policy)'));
    });

    test('omits formatEot import when no heredoc fields', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [{ name: 'path', type: 'string', required: true }],
      });
      assert.false(result.includes('formatEot'));
    });

    test('emits TODO comment for object fields', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [{ name: 'options', type: 'object', required: false }],
      });
      assert.true(result.includes('// TODO: options'));
    });

    test('includes the feature key in the registry comment', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [{ name: 'path', type: 'string', required: true }],
      });
      assert.true(result.includes("registry['secrets/kv']"));
    });

    test('includes the Terraform docs URL', function (assert) {
      const result = generateInteractiveScaffold({
        ...base,
        fields: [{ name: 'path', type: 'string', required: true }],
      });
      assert.true(result.includes('docs/resources/mount'));
    });
  });

  // ---------------------------------------------------------------------------
  // generateOpenApiScaffold
  // ---------------------------------------------------------------------------

  module('#generateOpenApiScaffold', function () {
    const base = {
      apiPath: '/sys/policies/acl/{name}',
      tfResource: 'vault_policy',
      matched: [],
      inOpenApiOnly: [],
      inTfOnly: [],
    };

    test('includes the API path in the header comment', function (assert) {
      const result = generateOpenApiScaffold({ ...base });
      assert.true(result.includes('API path:           /sys/policies/acl/{name}'));
    });

    test('includes the Terraform resource in the header comment', function (assert) {
      const result = generateOpenApiScaffold({ ...base });
      assert.true(result.includes('Terraform resource: vault_policy'));
    });

    test('derives the registry key hint from the API path', function (assert) {
      const result = generateOpenApiScaffold({ ...base });
      assert.true(result.includes("registry['policies/acl']"));
    });

    test('includes Terraform docs URL with vault_ stripped', function (assert) {
      const result = generateOpenApiScaffold({ ...base });
      assert.true(result.includes('docs/resources/policy'));
    });

    test('adds an omission note for inOpenApiOnly fields', function (assert) {
      const result = generateOpenApiScaffold({
        ...base,
        inOpenApiOnly: [{ name: 'token_ttl', description: 'TTL for token', openApiType: 'string' }],
      });
      assert.true(result.includes('In Vault API but not in vault_policy (omitted)'));
      assert.true(result.includes('token_ttl'));
    });

    test('adds an omission note for inTfOnly fields', function (assert) {
      const result = generateOpenApiScaffold({
        ...base,
        inTfOnly: ['disable_remount'],
      });
      assert.true(result.includes('In vault_policy but not in Vault API request body (omitted)'));
      assert.true(result.includes('disable_remount'));
    });

    test('omits both notes when there are no discrepancies', function (assert) {
      const result = generateOpenApiScaffold({ ...base });
      assert.false(result.includes('omitted'));
    });

    test('includes formatEot import when a matched field is heredoc', function (assert) {
      const result = generateOpenApiScaffold({
        ...base,
        matched: [{ name: 'policy', fieldType: 'heredoc', required: true }],
      });
      assert.true(result.includes("import { formatEot } from 'core/utils/code-generators/formatters'"));
    });

    test('omits formatEot import when no heredoc fields', function (assert) {
      const result = generateOpenApiScaffold({
        ...base,
        matched: [{ name: 'name', fieldType: 'string', required: true }],
      });
      assert.false(result.includes('formatEot'));
    });
  });
});
