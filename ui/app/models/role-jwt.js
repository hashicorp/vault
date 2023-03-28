/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import parseURL from 'core/utils/parse-url';

const DOMAIN_STRINGS = {
  github: 'GitHub',
  gitlab: 'GitLab',
  google: 'Google',
  ping: 'Ping',
  okta: 'Okta',
  auth0: 'Auth0',
};

const PROVIDER_WITH_LOGO = ['GitLab', 'Google', 'Auth0'];

export { DOMAIN_STRINGS, PROVIDER_WITH_LOGO };

export default Model.extend({
  authUrl: attr('string'),

  providerName: computed('authUrl', function () {
    const { hostname } = parseURL(this.authUrl);
    const firstMatch = Object.keys(DOMAIN_STRINGS).find((name) => hostname.includes(name));
    return DOMAIN_STRINGS[firstMatch] || null;
  }),

  providerButtonComponent: computed('providerName', function () {
    const { providerName } = this;
    return PROVIDER_WITH_LOGO.includes(providerName) ? `auth-button-${providerName.toLowerCase()}` : null;
  }),
});
