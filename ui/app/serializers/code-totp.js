/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from './application';

export default ApplicationSerializer.extend({
  normalizeItems(payload, requestType) {
    Object.assign(payload, payload.data);
    delete payload.data;
    return payload;
  },
});
