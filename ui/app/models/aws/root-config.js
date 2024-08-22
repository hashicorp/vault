/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { regions } from 'vault/helpers/aws-regions';

export default class AwsRootConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord

  // IAM only fields
  @attr('string') accessKey;
  @attr('string', { sensitive: true }) secretKey; // obfuscated, never returned by API

  // WIF only fields
  @attr('string', {
    subText: 'Role ARN to assume for plugin workload identity federation.',
  })
  roleArn;
  @attr('string', {
    subText:
      'The audience claim value for plugin identity tokens. Must match an allowed audience configured for the target IAM OIDC identity provider..',
  })
  identityTokenAudience;
  @attr({
    defaultValue: '1h',
    label: 'Identity token TTL',
    helperTextDisabled:
      'The TTL of generated tokens. Defaults to 1 hour, turn on the toggle to specify a different value.',
    helperTextEnabled:
      'The TTL of generated tokens. Defaults to 1 hour, turn on the toggle to specify a different value.',
    subText: '',
    editType: 'ttl',
  })
  identityTokenTll;

  // WIF + IAM fields
  @attr('string', {
    possibleValues: regions(),
    subText:
      'Specifies the AWS region. If not set it will use the AWS_REGION env var, AWS_DEFAULT_REGION env var, or us-east-1 in that order.',
  })
  region;
  @attr('string', { label: 'IAM endpoint' })
  iamEndpoint;
  @attr('string', { label: 'STS endpoint' }) stsEndpoint;
  @attr('number', {
    defaultValue: -1,
    label: 'Maximum retries',
    subText: 'Number of max retries the client should use for recoverable errors. Default is -1.',
  })
  maxRetries;
  // there are more options available on the API, but the UI does not support them yet.
  get attrs() {
    const keys = ['accessKey', 'region', 'iamEndpoint', 'stsEndpoint', 'maxRetries'];
    return expandAttributeMeta(this, keys);
  }

  get formFieldGroups() {
    return [
      { default: ['accessKey', 'secretKey', 'roleArn', 'identityTokenAudience', 'identityTokenTll'] },
      {
        'Root config options': ['region', 'iamEndpoint', 'stsEndpoint', 'maxRetries'],
      },
    ];
  }

  get fieldGroups() {
    return fieldToAttrs(this, this.formFieldGroups);
  }

  get fieldGroupsWif() {
    return fieldToAttrs(this, this.formFieldGroupsWif);
  }

  get fieldGroupsIam() {
    return fieldToAttrs(this, this.formFieldGroupsIam);
  }

  get formFieldGroupsWif() {
    return [
      { default: ['roleArn', 'identityTokenAudience', 'identityTokenTll'] },
      {
        'Root config options': ['region', 'iamEndpoint', 'stsEndpoint', 'maxRetries'],
      },
    ];
  }

  get formFieldGroupsIam() {
    return [
      { default: ['accessKey', 'secretKey'] },
      {
        'Root config options': ['region', 'iamEndpoint', 'stsEndpoint', 'maxRetries'],
      },
    ];
  }
}
