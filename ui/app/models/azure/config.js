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
  @attr('string', { label: 'Subscription ID' }) subscriptionId;
  @attr('string', { label: 'Tenant ID' }) tenantId;
  @attr('string', { label: 'Client ID' }) clientId;
  @attr('string', { sensitive: true }) clientSecret; // obfuscated, never returned by API

  @attr('string', {
    subText:
      'This value can also be provided with the AZURE_ENVIRONMENT environment variable. If not specified, Vault will use Azure Public Cloud.',
  })
  environment;

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

  @attr({
    label: 'Root password TTL',
    editType: 'ttl',
    // default is 15768000 sec. The api docs say 182 days, but this should be updated to 182.5 days.
    helperTextDisabled: 'Vault will use the default of 182 days.',
    helperTextEnabled:
      'Specifies how long the root password is valid for in Azure when rotate-root generates a new client secret.',
  })
  rootPasswordTtl;

  configurableParams = [
    'subscriptionId',
    'tenantId',
    'clientId',
    'clientSecret',
    'identityTokenAudience',
    'identityTokenTtl',
    'rootPasswordTtl',
    'environment',
  ];

  /* GETTERS used by configure-azure component 
  these getters help:
  1. determine if the model is new or existing
  2. if wif or azure attributes have been configured
  */
  get isConfigured() {
    // if every value is falsy, this engine has not been configured yet
    return !this.configurableParams.every((param) => !this[param]);
  }

  get isWifPluginConfigured() {
    return !!this.identityTokenAudience || !!this.identityTokenTtl;
  }

  // the "clientSecret" param is not checked because it's never return by the API.
  // thus we can never say for sure if the account accessType has been configured so we always return false
  isAccountPluginConfigured = false;

  /* GETTERS used to generate array of fields to be displayed in: 
  1. details view
  2. edit/create view
*/
  get displayAttrs() {
    const formFields = expandAttributeMeta(this, this.configurableParams);
    return formFields.filter((attr) => attr.name !== 'clientSecret');
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
    formFieldGroups.push({
      default: ['subscriptionId', 'tenantId', 'clientId'],
    });
    if (accessType === 'account') {
      formFieldGroups.push({
        default: ['clientSecret'],
      });
    }
    if (accessType === 'wif') {
      formFieldGroups.push({
        default: ['identityTokenAudience', 'identityTokenTtl'],
      });
    }
    formFieldGroups.push({
      'More options': ['environment', 'rootPasswordTtl'],
    });

    return formFieldGroups;
  }
}
