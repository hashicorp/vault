/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { terraformTemplate, formatTerraformOptions } from 'core/utils/code-generators/terraform';

module('Integration | Util | code-generators/terraform', function (hooks) {
  setupTest(hooks);

  test('formatTerraformOptions: it formats single option', async function (assert) {
    const options = { type: '"github"' };
    const formatted = formatTerraformOptions(options);
    const expected = ['  type = "github"'];
    assert.propEqual(formatted, expected, 'it formats single option correctly');
  });

  test('formatTerraformOptions: it formats multiple options', async function (assert) {
    const options = {
      type: '"userpass"',
      path: '"userpass"',
      description: '"User password auth"',
    };
    const formatted = formatTerraformOptions(options);
    assert.strictEqual(formatted.length, 3, 'it returns array with 3 items');
    assert.true(formatted.includes('  type = "userpass"'), 'it includes type option');
    assert.true(formatted.includes('  path = "userpass"'), 'it includes path option');
    assert.true(formatted.includes('  description = "User password auth"'), 'it includes description option');
  });

  test('formatTerraformOptions: it formats options with numbers', async function (assert) {
    const options = { max_ttl: 3600 };
    const formatted = formatTerraformOptions(options);
    const expected = ['  max_ttl = 3600'];
    assert.propEqual(formatted, expected, 'it formats numeric values');
  });

  test('formatTerraformOptions: it handles empty options', async function (assert) {
    const options = {};
    const formatted = formatTerraformOptions(options);
    assert.propEqual(formatted, [], 'it returns empty array for empty options');
  });

  test('terraformTemplate: it generates basic terraform resource', async function (assert) {
    const formatted = terraformTemplate({
      resource: 'vault_auth_backend',
      localId: 'example',
      options: { type: '"github"' },
    });
    const expected = `resource "vault_auth_backend" "example" {
  type = "github"
}`;
    assert.strictEqual(formatted, expected, 'it generates terraform resource');
  });

  test('terraformTemplate: it generates terraform resource with multiple options', async function (assert) {
    const formatted = terraformTemplate({
      resource: 'vault_mount',
      localId: 'kv',
      options: {
        path: '"secret"',
        type: '"kv-v2"',
        description: '"KV Version 2 secret engine"',
      },
    });

    assert.true(formatted.includes('resource "vault_mount" "kv"'), 'it includes resource declaration');
    assert.true(formatted.includes('path = "secret"'), 'it includes path option');
    assert.true(formatted.includes('type = "kv-v2"'), 'it includes type option');
    assert.true(formatted.includes('description = "KV Version 2 secret engine"'), 'it includes description');
  });

  test('terraformTemplate: it uses default localId when not provided', async function (assert) {
    const formatted = terraformTemplate({
      resource: 'vault_auth_backend',
      options: { type: '"github"' },
    });
    assert.true(formatted.includes('"<local identifier>"'), 'it uses default local identifier');
  });

  test('terraformTemplate: it handles empty resource name', async function (assert) {
    const formatted = terraformTemplate({
      resource: '',
      localId: 'test',
      options: { key: '"value"' },
    });
    assert.true(formatted.includes('resource "" "test"'), 'it handles empty resource name');
  });

  test('terraformTemplate: it formats multiple options with proper spacing', async function (assert) {
    const formatted = terraformTemplate({
      resource: 'vault_auth_backend',
      localId: 'userpass',
      options: {
        type: '"userpass"',
        path: '"userpass"',
        tune: '{}',
      },
    });

    // Check that options are separated by double newlines
    const lines = formatted.split('\n\n');
    assert.strictEqual(lines.length, 3, 'options are separated by blank lines');
  });
});
