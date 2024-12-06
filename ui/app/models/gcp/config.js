/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const credentialsExample = `# The example below is treated as a comment and will not be submitted
# some kind of example??
`;

export default class GcpConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord

  @attr('string', {
    editType: 'json',
    label: 'JSON credentials',
    helpText: 'Specifies the maximum Config TTL for long-lived credentials (i.e. service account keys).',
    example: credentialsExample,
    mode: 'ruby',
    sectionHeading: 'Hello section heading', // render section heading before form field
  })
  credentials; // Mutually exclusive with identityTokenAudience, meaning it's cannot be set if identityTokenAudience is set

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
    label: 'Config TTL',
    editType: 'ttl',
    helperTextDisabled: 'The TTL of generated tokens. Defaults to 1 hour.',
  })
  ttl;

  @attr({
    label: 'Max TTL',
    editType: 'ttl',
    helperTextDisabled:
      'Specifies the maximum Config TTL for long-lived credentials (i.e. service account keys).',
  })
  maxTtl;

  // for configuration details view
  get displayAttrs() {
    return this.formFields;
  }

  // formFields are iterated through to generate the edit/create view
  get formFields() {
    const keys = [
      'credentials',
      'serviceAccountEmail',
      'identityTokenAudience',
      'identityTokenTtl',
      'ttl',
      'maxTtl',
    ];
    return expandAttributeMeta(this, keys);
  }
}
