/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

// Note: while the API docs indicate subscriptionId and tenantId are required, the UI does not enforce this because the user may pass these values in as environment variables.
// https://developer.hashicorp.com/vault/api-docs/secret/azure#configure-access
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
  @attr('string', { sensitive: true }) clientSecret; // obfuscated, never returned by API
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
  @attr('string') environment;
  @attr({
    label: 'Root password TTL',
    editType: 'ttl',
    helperTextDisabled:
      'Specifies how long the root password is valid for in Azure when rotate-root generates a new client secret. Defaults to 182 days or 6 months, 1 day and 13 hours.',
  })
  rootPasswordTtl;

  // for configuration details view
  // do not include clientSecret because it is never returned by the API
  get attrs() {
    const keys = [
      'subscriptionId',
      'tenantId',
      'clientId',
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

  // return client secret for edit/create view
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
