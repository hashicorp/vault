/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

module.exports = function (environment) {
  const policy = {
    'default-src': ["'none'"],
    'script-src': ["'self'"],
    'script-src-elem': ["'*.vercel-ui-preview-hashicorp.vercel.app*'", "'*hashicorp.vercel.app*'"], //TODO check this
    'font-src': ["'self'"],
    'connect-src': ["'self'"],
    'img-src': ["'self'", 'data:'],
    'style-src': ["'unsafe-inline'", "'self'"],
    'media-src': ["'self'"],
    'form-action': ["'none'"],
  };

  policy['connect-src'].push('https://eu.i.posthog.com');

  return {
    delivery: ['header', 'meta'],
    enabled: environment !== 'production',
    failTests: true,
    policy,
    reportOnly: false,
  };
};
