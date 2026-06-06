/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

module.exports = function (environment) {
  const policy = {
    'default-src': ["'none'"],
    'script-src': ["'self'"],
    'font-src': ["'self'"],
    'connect-src': ["'self'"],
    'img-src': ["'self'", 'data:'],
    'style-src': ["'self'"],
    'media-src': ["'self'"],
    'form-action': ["'none'"],
  };

  policy['connect-src'].push('https://eu.i.posthog.com');

  if (environment === 'test') {
    policy['style-src'].push("'unsafe-inline'");
  }

  return {
    delivery: ['header', 'meta'],
    enabled: environment !== 'production',
    failTests: true,
    policy,
    reportOnly: false,
  };
};
