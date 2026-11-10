/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import RoleAwsForm from 'vault/forms/aws-role';

module('Unit | Form | aws role', function () {
  const fieldsByType = {
    iam_user: ['name', 'credential_type', 'policy_arns', 'policy_document'],
    assumed_role: ['name', 'credential_type', 'role_arns', 'policy_document'],
    federation_token: ['name', 'credential_type', 'policy_document'],
    session_token: ['name', 'credential_type'],
  };

  for (const [credential_type, fields] of Object.entries(fieldsByType)) {
    test(`fields and payload for ${credential_type}`, function (assert) {
      const form = new RoleAwsForm(
        {
          backend: 'aws',
          name: 'example',
          credential_type,
          role_arns: ['role-arn'],
          policy_arns: ['policy-arn'],
          policy_document: '{}',
        },
        { isNew: true }
      );
      assert.deepEqual(
        form.formFields.map((field) => field.name),
        fields
      );
      const { isValid, data } = form.toJSON();
      assert.true(isValid, 'A named role with a credential type is valid');
      assert.deepEqual(data, {
        credential_type,
        role_arns: credential_type === 'assumed_role' ? ['role-arn'] : [],
        policy_arns: credential_type === 'iam_user' ? ['policy-arn'] : [],
        policy_document: '',
      });
      assert.deepEqual(form.data.role_arns, ['role-arn'], 'Serialization does not mutate form inputs');
    });
  }

  test('requires a role name and credential type', function (assert) {
    const { isValid, state } = new RoleAwsForm().toJSON();
    assert.false(isValid);
    assert.false(state.name.isValid);
    assert.false(state.credential_type.isValid);
  });

  test('switching type updates fields and preserves policy text', function (assert) {
    const form = new RoleAwsForm({
      name: 'example',
      credential_type: 'iam_user',
      policy_document: '{"Version":"2012-10-17"}',
    });
    form.data.credential_type = 'session_token';
    assert.deepEqual(
      form.formFields.map((field) => field.name),
      ['name', 'credential_type']
    );
    assert.strictEqual(form.toJSON().data.policy_document, '{"Version":"2012-10-17"}');
  });
});
