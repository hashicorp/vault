/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import type { FieldOptions } from 'vault/utils/forms/field';

interface AwsCredentialData {
  credential_type?: string;
  ttl?: string;
  role_arn?: string;
}

const TTL_OPTIONS: FieldOptions & { setDefault: boolean; ttlOffValue: string } = {
  editType: 'ttl',
  defaultValue: '3600s',
  setDefault: true,
  ttlOffValue: '',
  label: 'TTL',
  helpText:
    'Specifies the TTL for the use of the STS token. Valid only when credential_type is assumed_role, federation_token, or session_token.',
};

/** Keeps generation inputs separate from the returned, one-time credentials. */
export default class AwsCredentialForm extends Form<AwsCredentialData> {
  formFields = [
    new FormField('credential_type', 'string', {
      label: 'Credential type',
      possibleValues: [
        { value: 'iam_user', displayName: 'IAM User' },
        { value: 'assumed_role', displayName: 'Assumed Role' },
        { value: 'federation_token', displayName: 'Federation Token' },
        { value: 'session_token', displayName: 'Session Token' },
      ],
    }),
    new FormField('ttl', undefined, TTL_OPTIONS),
    new FormField('role_arn', 'string', {
      label: 'Role ARN',
      helpText:
        'The ARN of the role to assume if credential_type on the Vault role is assumed_role. Optional if the role has a single role ARN; required otherwise.',
    }),
  ];

  constructor(data: AwsCredentialData = {}) {
    super({ credential_type: 'iam_user', ttl: '3600s', ...data });
  }

  toJSON() {
    const { credential_type, ttl, role_arn } = this.data;
    const data: Pick<AwsCredentialData, 'ttl' | 'role_arn'> = {};
    if (credential_type !== 'iam_user') {
      if (ttl) data.ttl = ttl;
      if (credential_type === 'assumed_role' && role_arn) data.role_arn = role_arn;
    }
    return super.toJSON(data);
  }
}
