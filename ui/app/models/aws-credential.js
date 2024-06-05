/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
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

const DISPLAY_FIELDS = ['accessKey', 'secretKey', 'securityToken', 'leaseId', 'renewable', 'leaseDuration'];
export default Model.extend({
  helpText:
    'For Vault roles of credential type iam_user, there are no inputs, just submit the form. Choose a type to change the input options.',
  role: attr('object', {
    readOnly: true,
  }),

  credentialType: attr('string', {
    defaultValue: 'iam_user',
    possibleValues: CREDENTIAL_TYPES,
    readOnly: true,
  }),

  roleArn: attr('string', {
    label: 'Role ARN',
    helpText:
      'The ARN of the role to assume if credential_type on the Vault role is assumed_role. Optional if the role has a single role ARN; required otherwise.',
  }),

  ttl: attr({
    editType: 'ttl',
    defaultValue: '3600s',
    setDefault: true,
    ttlOffValue: '',
    label: 'TTL',
    helpText:
      'Specifies the TTL for the use of the STS token. Valid only when credential_type is assumed_role, federation_token, or session_token.',
  }),
  leaseId: attr('string'),
  renewable: attr('boolean'),
  leaseDuration: attr('number'),
  accessKey: attr('string'),
  secretKey: attr('string'),
  securityToken: attr('string'),

  attrs: computed('credentialType', 'accessKey', 'securityToken', function () {
    const type = this.credentialType;
    const fieldsForType = {
      iam_user: ['credentialType'],
      assumed_role: ['credentialType', 'ttl', 'roleArn'],
      federation_token: ['credentialType', 'ttl'],
      session_token: ['credentialType', 'ttl'],
    };
    if (this.accessKey || this.securityToken) {
      return expandAttributeMeta(this, DISPLAY_FIELDS.slice(0));
    }
    return expandAttributeMeta(this, fieldsForType[type].slice(0));
  }),

  toCreds: computed('accessKey', 'secretKey', 'securityToken', 'leaseId', function () {
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
  }),
});
