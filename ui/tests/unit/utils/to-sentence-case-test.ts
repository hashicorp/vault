/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { toSentenceCase } from 'vault/utils/to-sentence-case';

module('Unit | Utility | to-sentence-case', function () {
  module('toSentenceCase', function () {
    test('capitalises the first character and lowercases the rest', function (assert) {
      assert.strictEqual(toSentenceCase('kubernetes'), 'Kubernetes');
      assert.strictEqual(toSentenceCase('userpass'), 'Userpass');
      assert.strictEqual(toSentenceCase('cubbyhole'), 'Cubbyhole');
    });

    test('replaces underscores with spaces and applies sentence case', function (assert) {
      assert.strictEqual(toSentenceCase('some_example_text'), 'Some example text');
      assert.strictEqual(toSentenceCase('secret_engine'), 'Secret engine');
      assert.strictEqual(toSentenceCase('secret_aws_engine'), 'Secret AWS engine');
      assert.strictEqual(toSentenceCase('auth_method'), 'Auth method');
    });

    test('replaces hyphens with spaces and applies sentence case', function (assert) {
      assert.strictEqual(toSentenceCase('client-id'), 'Client id');
      assert.strictEqual(toSentenceCase('lease-count-quota'), 'Lease count quota');
    });

    test('uppercases usage-reporting acronyms', function (assert) {
      assert.strictEqual(toSentenceCase('aws'), 'AWS');
      assert.strictEqual(toSentenceCase('pki'), 'PKI');
      assert.strictEqual(toSentenceCase('kv'), 'KV');
      assert.strictEqual(toSentenceCase('kmip'), 'KMIP');
      assert.strictEqual(toSentenceCase('saml'), 'SAML');
      assert.strictEqual(toSentenceCase('totp'), 'TOTP');
      assert.strictEqual(toSentenceCase('oidc'), 'OIDC');
      assert.strictEqual(toSentenceCase('mfa'), 'MFA');
    });

    test('uppercases acronyms wherever they appear', function (assert) {
      assert.strictEqual(toSentenceCase('aws_auth_method'), 'AWS auth method');
      assert.strictEqual(toSentenceCase('kv_pki_status'), 'KV PKI status');
      assert.strictEqual(toSentenceCase('oidc-jwt-auth'), 'OIDC JWT auth');
    });

    test('keeps normal words sentence-cased', function (assert) {
      assert.strictEqual(toSentenceCase('nomad'), 'Nomad');
    });

    test('preserves branded Vault engine names', function (assert) {
      assert.strictEqual(toSentenceCase('approle'), 'AppRole');
      assert.strictEqual(toSentenceCase('alicloud_dynamic'), 'AliCloud dynamic');
      assert.strictEqual(toSentenceCase('github_auth'), 'GitHub auth');
      assert.strictEqual(toSentenceCase('rabbitmq'), 'RabbitMQ');
    });

    test('handles an already sentence-cased string without double-capitalising', function (assert) {
      assert.strictEqual(toSentenceCase('Authentication methods'), 'Authentication methods');
    });

    test('handles an empty string', function (assert) {
      assert.strictEqual(toSentenceCase(''), '');
    });
  });

  module('acronymsOnly option', function () {
    test('preserves word casing while uppercasing acronyms', function (assert) {
      assert.strictEqual(toSentenceCase('userpass', { acronymsOnly: true }), 'userpass');
      assert.strictEqual(toSentenceCase('aws', { acronymsOnly: true }), 'AWS');
      assert.strictEqual(toSentenceCase('aws_auth', { acronymsOnly: true }), 'AWS auth');
    });

    test('applies branded word overrides', function (assert) {
      assert.strictEqual(toSentenceCase('rabbitmq', { acronymsOnly: true }), 'RabbitMQ');
      assert.strictEqual(toSentenceCase('approle', { acronymsOnly: true }), 'AppRole');
      assert.strictEqual(toSentenceCase('github', { acronymsOnly: true }), 'GitHub');
    });

    test('handles an empty string', function (assert) {
      assert.strictEqual(toSentenceCase('', { acronymsOnly: true }), '');
    });
  });
});
