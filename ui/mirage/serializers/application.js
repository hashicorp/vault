/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { JSONAPISerializer } from 'ember-cli-mirage';

export default JSONAPISerializer.extend({
  typeKeyForModel(model) {
    return model.modelName;
  },
});
