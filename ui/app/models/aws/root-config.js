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
    label: 'Role ARN',
    subText: 'Role ARN to assume for plugin workload identity federation.',
  })
  roleArn;
  @attr('string', {
    subText:
      'The audience claim value for plugin identity tokens. Must match an allowed audience configured for the target IAM OIDC identity provider.',
  })
  identityTokenAudience;
  @attr({
    label: 'Identity token TTL',
    helperTextDisabled:
      'The TTL of generated tokens. Defaults to 1 hour, turn on the toggle to specify a different value.',
    helperTextEnabled: 'The TTL of generated tokens.',
    editType: 'ttl',
  })
  identityTokenTtl;

  // Fields that show regardless of access type
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
    subText: 'Number of max retries the client should use for recoverable errors. Default is -1.',
  })
  maxRetries;

  configurableParams = [
    'roleArn',
    'identityTokenAudience',
    'identityTokenTtl',
    'accessKey',
    'secretKey',
    'region',
    'iamEndpoint',
    'stsEndpoint',
    'maxRetries',
  ];

  get isWifPluginConfigured() {
    return !!this.identityTokenAudience || !!this.identityTokenTtl || !!this.roleArn;
  }

  get isAccountPluginConfigured() {
    return !!this.accessKey;
  }

  get displayAttrs() {
    const formFields = expandAttributeMeta(this, this.configurableParams);
    return formFields.filter((attr) => attr.name !== 'secretKey');
  }

  // "filedGroupsWif" and "fieldGroupsAccount" are passed to the FormFieldGroups component to determine which group to show in the form (ex: @groupName="fieldGroupsWif")
  get fieldGroupsWif() {
    return fieldToAttrs(this, this.formFieldGroups('wif'));
  }

  get fieldGroupsAccount() {
    return fieldToAttrs(this, this.formFieldGroups('account'));
  }

  formFieldGroups(accessType = 'account') {
    const formFieldGroups = [];
    if (accessType === 'wif') {
      formFieldGroups.push({ default: ['roleArn', 'identityTokenAudience', 'identityTokenTtl'] });
    }
    if (accessType === 'account') {
      formFieldGroups.push({ default: ['accessKey', 'secretKey'] });
    }
    formFieldGroups.push({
      'Root config options': ['region', 'iamEndpoint', 'stsEndpoint', 'maxRetries'],
    });
    return formFieldGroups;
  }
}
