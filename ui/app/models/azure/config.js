/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class AzureConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', {
    label: 'Subscription ID',
  })
  subscriptionId;
  @attr('string', {
    label: 'Tenant ID',
  })
  tenantId;
  @attr('string', {
    label: 'Client ID',
  })
  clientId;
  @attr('string') clientSecret;
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
    subText: '',
    editType: 'ttl',
  })
  identityTokenTtl;
  @attr('string') environment;
  @attr({
    label: 'Root password TTL',
    editType: 'ttl',
  })
  rootPasswordTtl;

  get attrs() {
    const keys = [
      'subscriptionId',
      'tenantId',
      'clientId',
      'clientSecret',
      'identityTokenAudience',
      'identityTokenTtl',
      'environment',
      'rootPasswordTtl',
    ];
    return expandAttributeMeta(this, keys);
  }

  // "filedGroupsWif" and "fieldGroupsAzure" are passed to the FormFieldGroups component to determine which group to show in the form (ex: @groupName="fieldGroupsWif")
  get fieldGroupsWif() {
    return fieldToAttrs(this, this.formFieldGroups('wif'));
  }

  get fieldGroupsAzure() {
    return fieldToAttrs(this, this.formFieldGroups('azure'));
  }

  formFieldGroups(accessType = 'azure') {
    const formFieldGroups = [];
    formFieldGroups.push({
      default: ['subscriptionId', 'tenantId', 'clientId', 'environment'],
    });
    if (accessType === 'wif') {
      formFieldGroups.push({
        default: ['identityTokenAudience', 'identityTokenTtl'],
      });
    }
    if (accessType === 'azure') {
      formFieldGroups.push({
        default: ['clientSecret', 'rootPasswordTtl'],
      });
    }
    return formFieldGroups;
  }

  // return client and secret key for edit/create view
  get formFields() {
    const keys = [
      'subscriptionId',
      'tenantId',
      'clientId',
      'clientSecret',
      'identityTokenAudience',
      'identityTokenTtl',
      'rootPasswordTtl',
      'environment',
    ];
    return expandAttributeMeta(this, keys);
  }
}
