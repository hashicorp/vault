/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { renderTerraformBlocks, ref } from 'vault/utils/terraform-code-generator/terraform-registry';
import type { TerraformBlock } from 'vault/utils/terraform-code-generator/terraform-registry';

module('Unit | Utility | terraform-registry', function () {
  module('ref', function () {
    test('with attribute returns three-part reference', function (assert) {
      assert.strictEqual(ref('vault_mount', 'kv', 'path'), 'vault_mount.kv.path');
    });

    test('without attribute returns two-part reference', function (assert) {
      assert.strictEqual(ref('vault_mount', 'kv'), 'vault_mount.kv');
    });
  });

  module('renderTerraformBlocks', function () {
    test('joins blocks with one blank line', function (assert) {
      const blocks: TerraformBlock[] = [
        { type: 'variable', content: 'variable "example" {}' },
        { type: 'resource', content: 'resource "vault_policy" "example" {}' },
      ];
      assert.strictEqual(
        renderTerraformBlocks(blocks),
        'variable "example" {}\n\nresource "vault_policy" "example" {}'
      );
    });

    test('preserves caller-supplied block order', function (assert) {
      const blocks: TerraformBlock[] = [
        { type: 'resource', content: 'resource "a" {}' },
        { type: 'variable', content: 'variable "b" {}' },
        { type: 'resource', content: 'resource "c" {}' },
      ];
      assert.strictEqual(
        renderTerraformBlocks(blocks),
        'resource "a" {}\n\nvariable "b" {}\n\nresource "c" {}'
      );
    });

    test('returns empty string for empty input', function (assert) {
      assert.strictEqual(renderTerraformBlocks([]), '');
    });
  });
});
