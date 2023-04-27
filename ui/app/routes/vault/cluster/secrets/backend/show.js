/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import EditBase from './secret-edit';

export default EditBase.extend({
  queryParams: {
    version: {
      refreshModel: true,
    },
  },
});
