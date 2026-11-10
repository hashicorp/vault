/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import AwsCredentialForm from 'vault/forms/aws-credential';

module('Unit | Form | aws credential', function () {
  test('defaults to IAM and strips generation-only metadata from the payload', function (assert) {
    const form = new AwsCredentialForm();
    assert.strictEqual(form.data.credential_type, 'iam_user', 'Defaults to IAM when role read is denied');
    assert.strictEqual(form.data.ttl, '3600s', 'STS forms start with a one-hour TTL');
    assert.deepEqual(form.toJSON().data, {}, 'IAM never sends TTL or credential type');
    form.set('credential_type', 'assumed_role');
    form.set('role_arn', 'arn:aws:iam::123456789012:role/example');
    assert.deepEqual(
      form.toJSON().data,
      { ttl: '3600s', role_arn: 'arn:aws:iam::123456789012:role/example' },
      'FormField updates are serialized from tracked form data'
    );
    assert.strictEqual(form.data.credential_type, 'assumed_role', 'Serialization does not mutate inputs');
  });

  test('switching credential type omits stale ARN and disabled TTL', function (assert) {
    const form = new AwsCredentialForm({ credential_type: 'assumed_role', ttl: '', role_arn: 'arn' });
    assert.deepEqual(form.toJSON().data, { role_arn: 'arn' }, 'Disabling TTL omits it');
    form.set('credential_type', 'federation_token');
    assert.deepEqual(form.toJSON().data, {}, 'Federation tokens omit a previously entered ARN');
    form.set('ttl', '3h');
    form.set('credential_type', 'session_token');
    assert.deepEqual(form.toJSON().data, { ttl: '3h' }, 'Session tokens accept only TTL');
    form.set('credential_type', 'iam_user');
    assert.deepEqual(form.toJSON().data, {}, 'IAM ignores all previous STS inputs');
  });
});
