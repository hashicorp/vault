/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'config/control-group';
  },

  urlForDeleteRecord(id, modelName) {
    return this.buildURL(modelName);
  },

  urlForFindRecord(id, modelName) {
    return this.buildURL(modelName);
  },

  urlForUpdateRecord(id, modelName) {
    return this.buildURL(modelName);
  },
});
