/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import {
  aclTemplate,
  formatCapabilities,
  isAclCapability,
  ACL_CAPABILITIES,
  PolicyStanza,
} from 'core/utils/code-generators/policy';

module('Integration | Util | code-generators/policy', function (hooks) {
  setupTest(hooks);

  test('aclTemplate: it formats a policy', async function (assert) {
    const formatted = aclTemplate('my-path/*', ['list', 'read', 'delete']);
    const expected = `path "my-path/*" {
    capabilities = ["read", "delete", "list"]
}`;
    assert.strictEqual(formatted, expected, 'it formats an ACL policy');
  });

  test('aclTemplate: it handles empty path and capabilities', async function (assert) {
    const formatted = aclTemplate('', []);
    const expected = `path "" {
    capabilities = []
}`;
    assert.strictEqual(formatted, expected, 'it formats empty policy');
  });

  test('aclTemplate: it handles single capability', async function (assert) {
    const formatted = aclTemplate('auth/token/lookup-self', ['read']);
    const expected = `path "auth/token/lookup-self" {
    capabilities = ["read"]
}`;
    assert.strictEqual(formatted, expected, 'it formats policy with single capability');
  });

  test('formatCapabilities: it formats capabilities in consistent order', async function (assert) {
    const formatted = formatCapabilities(['list', 'read', 'delete']);
    const expected = '"read", "delete", "list"';
    assert.strictEqual(formatted, expected, 'it formats capabilities in ACL_CAPABILITIES order');
  });

  test('formatCapabilities: it filters out invalid capabilities', async function (assert) {
    const formatted = formatCapabilities(['read', 'invalid', 'list']);
    const expected = '"read", "list"';
    assert.strictEqual(formatted, expected, 'it filters out invalid capabilities');
  });

  test('formatCapabilities: it returns empty string for empty array', async function (assert) {
    const formatted = formatCapabilities([]);
    assert.strictEqual(formatted, '', 'it returns empty string for no capabilities');
  });

  test('formatCapabilities: it handles single capability', async function (assert) {
    const formatted = formatCapabilities(['read']);
    const expected = '"read"';
    assert.strictEqual(formatted, expected, 'it formats single capability');
  });

  test('formatCapabilities: it handles all capabilities', async function (assert) {
    const sorted = [...ACL_CAPABILITIES].sort(); // alphabetize so input order is different than expected output
    const formatted = formatCapabilities(sorted);
    const expected = '"create", "read", "update", "delete", "list", "patch", "sudo"';
    assert.strictEqual(formatted, expected, 'it formats all capabilities in order');
  });

  test('isAclCapability: it returns true for valid capabilities', async function (assert) {
    ACL_CAPABILITIES.forEach((cap) => {
      assert.true(isAclCapability(cap), `${cap} is a valid capability`);
    });
  });

  test('isAclCapability: it returns false for invalid capabilities', async function (assert) {
    assert.false(isAclCapability('invalid'), 'invalid is not a valid capability');
    assert.false(isAclCapability('write'), 'write is not a valid capability');
    assert.false(isAclCapability(''), 'empty string is not a valid capability');
    assert.false(isAclCapability('READ'), 'uppercase READ is not a valid capability');
  });

  test('PolicyStanza: it initializes with empty capabilities and path', async function (assert) {
    const stanza = new PolicyStanza();
    assert.strictEqual(stanza.path, '', 'path is empty');
    assert.strictEqual(stanza.capabilities.size, 0, 'capabilities set is empty');
  });

  test('PolicyStanza: it sets path when instantiated with a path value', async function (assert) {
    const stanza = new PolicyStanza({ path: 'my-path' });
    assert.strictEqual(stanza.path, 'my-path', 'path is sets');
  });

  test('PolicyStanza: it generates preview for single capability', async function (assert) {
    const stanza = new PolicyStanza();
    stanza.path = 'secret/data/*';
    stanza.capabilities.add('read');

    const expected = `path "secret/data/*" {
    capabilities = ["read"]
}`;
    assert.strictEqual(stanza.preview, expected, 'it generates correct preview');
  });

  test('PolicyStanza: it generates preview for multiple capabilities', async function (assert) {
    const stanza = new PolicyStanza();
    stanza.path = 'auth/*';
    stanza.capabilities.add('list');
    stanza.capabilities.add('read');
    stanza.capabilities.add('create');

    const expected = `path "auth/*" {
    capabilities = ["create", "read", "list"]
}`;
    assert.strictEqual(stanza.preview, expected, 'it generates preview with multiple capabilities');
  });

  test('PolicyStanza: it generates preview without path and capabilities', async function (assert) {
    const stanza = new PolicyStanza();
    const expected = `path "" {
    capabilities = []
}`;
    assert.strictEqual(stanza.preview, expected, 'it generates preview with empty capabilities');
  });

  test('PolicyStanza: it updates preview when capabilities change', async function (assert) {
    const stanza = new PolicyStanza();
    stanza.path = 'secret/*';
    stanza.capabilities.add('read');

    const firstPreview = stanza.preview;

    stanza.capabilities.add('list');
    const secondPreview = stanza.preview;
    const expected = `path "secret/*" {
    capabilities = ["read", "list"]
}`;
    assert.notStrictEqual(firstPreview, secondPreview, 'preview updates when capabilities change');
    assert.true(secondPreview.includes('"read", "list"'), 'new preview includes both capabilities');
    assert.strictEqual(secondPreview, expected, 'new preview reflects updates');
  });
});
