/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
self.deprecationWorkflow.config = {
  throwOnUnhandled: false,
};

self.deprecationWorkflow.config = {
  // current output from deprecationWorkflow.flushDeprecations();
  workflow: [
    { handler: 'silence', matchId: 'ember-engines.deprecation-router-service-from-host' },
    // ember-data
    { handler: 'silence', matchId: 'ember-data:deprecate-early-static' }, // decorator tests
    { handler: 'silence', matchId: 'ember-data:deprecate-promise-proxies' }, // Transform secrets
    { handler: 'silence', matchId: 'ember-data:no-a-with-array-like' }, // MFA
    { handler: 'silence', matchId: 'ember-data:deprecate-promise-many-array-behaviors' }, // MFA
  ],
};
