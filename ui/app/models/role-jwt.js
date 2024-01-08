/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import parseURL from 'core/utils/parse-url';

const DOMAIN_STRINGS = {
  'github.com': 'GitHub',
  'gitlab.com': 'GitLab',
  'google.com': 'Google',
  'ping.com': 'Ping',
  'okta.com': 'Okta',
  'auth0.com': 'Auth0',
};

const PROVIDER_WITH_LOGO = ['GitHub', 'GitLab', 'Google', 'Okta', 'Auth0'];

export { DOMAIN_STRINGS, PROVIDER_WITH_LOGO };

export default class RoleJwtModel extends Model {
  @attr('string') authUrl;

  get providerName() {
    const { hostname } = parseURL(this.authUrl);
    const firstMatch = Object.keys(DOMAIN_STRINGS).find((name) => hostname.includes(name));
    return DOMAIN_STRINGS[firstMatch] || null;
  }

  get providerIcon() {
    const { providerName } = this;
    return PROVIDER_WITH_LOGO.includes(providerName) ? providerName.toLowerCase() : null;
  }
}
