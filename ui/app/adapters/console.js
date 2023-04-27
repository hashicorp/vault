/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  pathForType(modelName) {
    return modelName;
  },
});
