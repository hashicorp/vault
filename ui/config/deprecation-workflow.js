/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* global self */
self.deprecationWorkflow = self.deprecationWorkflow || {};
//self.deprecationWorkflow.config = {
//throwOnUnhandled: true
//}
self.deprecationWorkflow.config = {
  // current output from deprecationWorkflow.flushDeprecations();
  // deprecations that will not be removed until 5.0.0 are filtered by deprecation-filter initializer rather than silencing below
  workflow: [
    { handler: 'silence', matchId: 'ember-data:model-save-promise' },
    { handler: 'silence', matchId: 'ember-engines.deprecation-camelized-engine-names' },
    { handler: 'silence', matchId: 'ember-engines.deprecation-router-service-from-host' },
    { handler: 'silence', matchId: 'ember-modifier.use-modify' },
    { handler: 'silence', matchId: 'ember-modifier.no-element-property' },
    { handler: 'silence', matchId: 'ember-modifier.no-args-property' },
    { handler: 'silence', matchId: 'ember-cli-mirage-config-routes-only-export' },
    { handler: 'silence', matchId: 'setting-on-hash' },
  ],
};
