/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { terraformTemplate, formatTerraformResourceArgs } from 'core/utils/code-generators/terraform';

module('Integration | Util | code-generators/terraform', function (hooks) {
  setupTest(hooks);

  test('formatTerraformResourceArgs: it formats single option', async function (assert) {
    const resourceArgs = { type: '"github"' };
    const formatted = formatTerraformResourceArgs(resourceArgs);
    const expected = ['  type = "github"'];
    assert.propEqual(formatted, expected, 'it formats single option correctly');
  });

  test('formatTerraformResourceArgs: it formats multiple resourceArgs', async function (assert) {
    const resourceArgs = {
      type: '"userpass"',
      path: '"userpass"',
      description: '"User password auth"',
    };
    const formatted = formatTerraformResourceArgs(resourceArgs);
    assert.strictEqual(formatted.length, 3, 'it returns array with 3 items');
    assert.true(formatted.includes('  type = "userpass"'), 'it includes type option');
    assert.true(formatted.includes('  path = "userpass"'), 'it includes path option');
    assert.true(formatted.includes('  description = "User password auth"'), 'it includes description option');
  });

  test('formatTerraformResourceArgs: it formats resourceArgs with numbers', async function (assert) {
    const resourceArgs = { max_ttl: 3600 };
    const formatted = formatTerraformResourceArgs(resourceArgs);
    const expected = ['  max_ttl = 3600'];
    assert.propEqual(formatted, expected, 'it formats numeric values');
  });

  test('formatTerraformResourceArgs: it handles empty resourceArgs', async function (assert) {
    const resourceArgs = {};
    const formatted = formatTerraformResourceArgs(resourceArgs);
    assert.propEqual(formatted, [], 'it returns empty array for empty resourceArgs');
  });

  test('formatTerraformResourceArgs: it handles undefined resourceArgs', async function (assert) {
    const formatted = formatTerraformResourceArgs();
    assert.propEqual(formatted, [], 'it returns empty array for empty resourceArgs');
  });

  test('terraformTemplate: it formats terraform', async function (assert) {
    const formatted = terraformTemplate({
      resource: 'vault_auth_backend',
      localId: 'example',
      resourceArgs: { type: '"github"' },
    });
    const expected = `resource "vault_auth_backend" "example" {
  type = "github"
}`;
    assert.strictEqual(formatted, expected, 'it returns formatted terraform snippet');
  });

  test('terraformTemplate: it formats terraform with multiple resourceArgs', async function (assert) {
    const formatted = terraformTemplate({
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

  test('terraformTemplate: it formats terraform with nested object', async function (assert) {
    const formatted = terraformTemplate({
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

  test('terraformTemplate: it uses default localId when not provided', async function (assert) {
    const formatted = terraformTemplate({
      resource: 'vault_auth_backend',
      resourceArgs: { type: '"github"' },
    });
    assert.true(formatted.includes('"<local identifier>"'), 'it uses default local identifier');
  });

  test('terraformTemplate: it handles empty resource name', async function (assert) {
    const formatted = terraformTemplate({
      resource: '',
      localId: 'test',
      resourceArgs: { key: '"value"' },
    });
    assert.true(formatted.includes('resource "" "test"'), 'it handles empty resource name');
  });

  test('terraformTemplate: it formats with linebreaks', async function (assert) {
    const formatted = terraformTemplate({
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

  test('terraformTemplate: it handles undefined args', async function (assert) {
    const formatted = terraformTemplate();
    const expected = `resource "<resource name>" "<local identifier>" {

}`;
    assert.strictEqual(formatted, expected, 'it returns default template when args are undefined');
  });
});
