/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType(modelName) {
    return modelName;
  },
});
