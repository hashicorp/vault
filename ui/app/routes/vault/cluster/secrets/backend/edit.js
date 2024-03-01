/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EditBase from './secret-edit';

export default EditBase.extend({
  queryParams: {
    version: {
      refreshModel: true,
    },
  },
});
