/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

export default AuthConfig.extend({
  orgName: attr('string', {
    helpText: 'Name of the organization to be used in the Okta API',
  }),
  apiToken: attr('string', {
    helpText:
      'Okta API token. This is required to query Okta for user group membership. If this is not supplied only locally configured groups will be enabled.',
  }),
  baseUrl: attr('string', {
    helpText:
      'If set, will be used as the base domain for API requests. Examples are okta.com, oktapreview.com, and okta-emea.com',
  }),
  bypassOktaMfa: attr('boolean', {
    defaultValue: false,
    helpText:
      "Useful if using Vault's built-in MFA mechanisms. Will also cause certain other statuses to be ignored, such as PASSWORD_EXPIRED",
  }),

  fieldGroups: computed('newFields', function () {
    let groups = [
      {
        default: ['orgName'],
      },
      {
        Options: ['apiToken', 'baseUrl', 'bypassOktaMfa'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }

    return fieldToAttrs(this, groups);
  }),
});
