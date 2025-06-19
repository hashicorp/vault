/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { tracked } from '@glimmer/tracking';

export default class WifConfigForm<T extends object> extends Form<T> {
  // for community users they will not be able to change this. for enterprise users, they will have the option to select "wif".
  @tracked accessType: 'account' | 'wif' = 'account';

  declare type: 'aws' | 'azure' | 'gcp';

  commonWifFields = {
    issuer: new FormField('issuer', 'string', {
      label: 'Issuer',
      subText:
        "The Issuer URL to be used in configuring Vault as an identity provider. If not set, Vault's default issuer will be used.",
      docLink: '/vault/api-docs/secret/identity/tokens#configure-the-identity-tokens-backend',
      placeholder: 'https://vault-test.com',
    }),

    identityTokenAudience: new FormField('identityTokenAudience', 'string', {
      subText:
        'The audience claim value for plugin identity tokens. Must match an allowed audience configured for the target IAM OIDC identity provider.',
    }),

    identityTokenTtl: new FormField('identityTokenTtl', 'string', {
      label: 'Identity token TTL',
      helperTextDisabled:
        'The TTL of generated tokens. Defaults to 1 hour, turn on the toggle to specify a different value.',
      helperTextEnabled: 'The TTL of generated tokens.',
      editType: 'ttl',
    }),

    serviceAccountEmail: new FormField('serviceAccountEmail', 'string', {
      subText: 'Email ID for the Service Account to impersonate for Workload Identity Federation.',
    }),
  };
}
