/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  terraformResourceTemplate,
  terraformVariableTemplate,
  formatTerraformArgs,
} from 'core/utils/code-generators/terraform';

module('Integration | Util | code-generators/terraform', function (hooks) {
  setupTest(hooks);

  test('formatTerraformArgs: it formats single option', async function (assert) {
    const resourceArgs = { type: '"github"' };
    const formatted = formatTerraformArgs(resourceArgs);
    const expected = ['  type = "github"'];
    assert.propEqual(formatted, expected, 'it formats single option correctly');
  });

  test('formatTerraformArgs: it formats multiple resourceArgs', async function (assert) {
    const resourceArgs = {
      type: '"userpass"',
      path: '"userpass"',
      description: '"User password auth"',
    };
    const formatted = formatTerraformArgs(resourceArgs);
    assert.strictEqual(formatted.length, 3, 'it returns array with 3 items');
    assert.true(formatted.includes('  type = "userpass"'), 'it includes type option');
    assert.true(formatted.includes('  path = "userpass"'), 'it includes path option');
    assert.true(formatted.includes('  description = "User password auth"'), 'it includes description option');
  });

  test('formatTerraformArgs: it formats resourceArgs with numbers', async function (assert) {
    const resourceArgs = { max_ttl: 3600 };
    const formatted = formatTerraformArgs(resourceArgs);
    const expected = ['  max_ttl = 3600'];
    assert.propEqual(formatted, expected, 'it formats numeric values');
  });

  test('formatTerraformArgs: it handles empty resourceArgs', async function (assert) {
    const resourceArgs = {};
    const formatted = formatTerraformArgs(resourceArgs);
    assert.propEqual(formatted, [], 'it returns empty array for empty resourceArgs');
  });

  test('formatTerraformArgs: it handles undefined resourceArgs', async function (assert) {
    const formatted = formatTerraformArgs();
    assert.propEqual(formatted, [], 'it returns empty array for empty resourceArgs');
  });

  test('terraformResourceTemplate: it formats terraform', async function (assert) {
    const formatted = terraformResourceTemplate({
      resource: 'vault_auth_backend',
      localId: 'example',
      resourceArgs: { type: '"github"' },
    });
    const expected = `resource "vault_auth_backend" "example" {
  type = "github"
}`;
    assert.strictEqual(formatted, expected, 'it returns formatted terraform snippet');
  });

  test('terraformResourceTemplate: it formats terraform with multiple resourceArgs', async function (assert) {
    const formatted = terraformResourceTemplate({
      resource: 'vault_mount',
      localId: 'secrets',
      resourceArgs: {
        path: '"secret"',
        type: '"generic"',
        description: '"A generic secret engine"',
      },
    });

    assert.true(formatted.includes('resource "vault_mount" "secrets"'), 'it includes resource declaration');
    assert.true(formatted.includes('path = "secret"'), 'it includes path option');
    assert.true(formatted.includes('type = "generic"'), 'it includes type option');
    assert.true(formatted.includes('description = "A generic secret engine"'), 'it includes description');
  });

  test('terraformResourceTemplate: it formats terraform with nested object', async function (assert) {
    const formatted = terraformResourceTemplate({
      resource: 'vault_mount',
      localId: 'kv',
      resourceArgs: {
        description: '"This is an example KV Version 2 secret engine mount"',
        path: '"my-kv-path"',
        type: '"kv-v2"',
        options: { version: '"2"', type: '"kv-v2"' },
      },
    });
    const expected = `resource "vault_mount" "kv" {
  description = "This is an example KV Version 2 secret engine mount"
  path = "my-kv-path"
  type = "kv-v2"
  options = {
    version = "2"
    type = "kv-v2"
  }
}`;
    assert.strictEqual(formatted, expected, 'it returns formatted terraform snippet');
  });

  test('terraformResourceTemplate: it uses default localId when not provided', async function (assert) {
    const formatted = terraformResourceTemplate({
      resource: 'vault_auth_backend',
      resourceArgs: { type: '"github"' },
    });
    assert.true(formatted.includes('"<local identifier>"'), 'it uses default local identifier');
  });

  test('terraformResourceTemplate: it handles empty resource name', async function (assert) {
    const formatted = terraformResourceTemplate({
      resource: '',
      localId: 'test',
      resourceArgs: { key: '"value"' },
    });
    assert.true(formatted.includes('resource "" "test"'), 'it handles empty resource name');
  });

  test('terraformResourceTemplate: it formats with linebreaks', async function (assert) {
    const formatted = terraformResourceTemplate({
      resource: 'vault_auth_backend',
      localId: 'userpass',
      resourceArgs: {
        type: '"userpass"',
        path: '"userpass"',
      },
    });
    // Check that snippet is separated by newlines
    const lines = formatted.split('\n');
    assert.strictEqual(lines.length, 4, 'the formatted snippet renders on 4 separate lines');
  });

  test('terraformResourceTemplate: it handles undefined args', async function (assert) {
    const formatted = terraformResourceTemplate();
    const expected = `resource "<resource name>" "<local identifier>" {

}`;
    assert.strictEqual(formatted, expected, 'it returns default template when args are undefined');
  });

  test('terraformVariableTemplate: it formats basic variable', async function (assert) {
    const formatted = terraformVariableTemplate({
      variable: 'namespace_path',
      variableArgs: {
        description: '"The namespace path"',
        type: 'string',
      },
    });
    const expected = `variable "namespace_path" {
  description = "The namespace path"
  type = string
}`;
    assert.strictEqual(formatted, expected, 'it returns formatted variable snippet');
  });

  test('terraformVariableTemplate: it uses defaults when args not provided', async function (assert) {
    const formatted = terraformVariableTemplate();
    const expected = `variable "<variable name>" {

}`;
    assert.strictEqual(formatted, expected, 'it returns default variable template');
  });

  test('terraformVariableTemplate: it handles empty variable name', async function (assert) {
    const formatted = terraformVariableTemplate({
      variable: '',
      variableArgs: { type: 'string' },
    });
    assert.true(formatted.includes('variable ""'), 'it handles empty variable name');
  });
});
