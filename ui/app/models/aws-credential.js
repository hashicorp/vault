/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

const CREDENTIAL_TYPES = [
  {
    value: 'iam_user',
    displayName: 'IAM User',
  },
  {
    value: 'assumed_role',
    displayName: 'Assumed Role',
  },
  {
    value: 'federation_token',
    displayName: 'Federation Token',
  },
  {
    value: 'session_token',
    displayName: 'Session Token',
  },
];

@withExpandedAttributes()
export default class AwsCredential extends Model {
  @attr('object', {
    readOnly: true,
  })
  role;

  @attr('string', {
    defaultValue: 'iam_user',
    possibleValues: CREDENTIAL_TYPES,
    readOnly: true,
  })
  credentialType;

  @attr('string', {
    label: 'Role ARN',
    helpText:
      'The ARN of the role to assume if credential_type on the Vault role is assumed_role. Optional if the role has a single role ARN; required otherwise.',
  })
  roleArn;

  @attr({
    editType: 'ttl',
    defaultValue: '3600s',
    setDefault: true,
    ttlOffValue: '',
    label: 'TTL',
    helpText:
      'Specifies the TTL for the use of the STS token. Valid only when credential_type is assumed_role, federation_token, or session_token.',
  })
  ttl;

  @attr('string') leaseId;
  @attr('boolean') renewable;
  @attr('number') leaseDuration;
  @attr('string') accessKey;
  @attr('string', { masked: true }) secretKey;
  @attr('string', { masked: true }) securityToken;

  get toCreds() {
    const props = {
      accessKey: this.accessKey,
      secretKey: this.secretKey,
      securityToken: this.securityToken,
      leaseId: this.leaseId,
    };
    const propsWithVals = Object.keys(props).reduce((ret, prop) => {
      if (props[prop]) {
        ret[prop] = props[prop];
        return ret;
      }
      return ret;
    }, {});
    return JSON.stringify(propsWithVals, null, 2);
  }
}
