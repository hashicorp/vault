/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { JSONAPISerializer } from 'miragejs';

export default JSONAPISerializer.extend({
  typeKeyForModel(model) {
    return model.modelName;
  },
});
