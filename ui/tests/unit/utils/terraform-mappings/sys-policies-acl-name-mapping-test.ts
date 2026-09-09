/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { sysPoliciesAclNameMapping } from 'vault/utils/terraform-mappings/sys-policies-acl-name-mapping';

module('Unit | Utility | terraform-mappings/sys-policies-acl-name-mapping', function () {
  test('it renders a vault_policy resource block', function (assert) {
    const secretReaderPayload = {
      name: 'my-policy',
      policy: 'path "secret/*" { capabilities = ["read"] }',
    };
    const tfResourceString = sysPoliciesAclNameMapping(secretReaderPayload);
    const expectedString = `resource "vault_policy" "<local identifier>" {\n  name = "my-policy"\n  policy = <<EOT\npath "secret/*" { capabilities = ["read"] }\nEOT\n}`;

    assert.strictEqual(tfResourceString, expectedString);
  });

  test('it uses the policy name fallback when name is empty', function (assert) {
    const noNamePayload = {
      name: '',
      policy: 'path "secret/*" { capabilities = ["read"] }',
    };
    const tfResourceString = sysPoliciesAclNameMapping(noNamePayload);

    assert.true(
      tfResourceString.includes('name = "<policy name>"'),
      'should use <policy name> when name is empty'
    );
  });

  test('it emits policy as a heredoc', function (assert) {
    const emptyPolicyPayload = {
      name: 'policy_in_name_only',
      policy: '',
    };
    const tfResourceString = sysPoliciesAclNameMapping(emptyPolicyPayload);

    assert.true(
      tfResourceString.includes('policy = <<EOT\n\nEOT'),
      'should emit an empty heredoc for empty policy'
    );
  });
});
