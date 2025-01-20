/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'control-group';
  },

  async findRecord(store, type, id) {
    const baseUrl = this.buildURL(type.modelName);
    return this.ajax(`${baseUrl}/request`, 'POST', {
      data: {
        accessor: id,
      },
    }).then((response) => {
      response.id = id;
      return response;
    });
  },

  urlForUpdateRecord(id, modelName) {
    const base = this.buildURL(modelName);
    return `${base}/authorize`;
  },
});
