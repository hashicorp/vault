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
  workflow: [],
};
