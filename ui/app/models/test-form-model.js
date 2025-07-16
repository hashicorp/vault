/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// this model is just used for integration tests
//

import AuthMethodModel from './auth-method';
import { belongsTo } from '@ember-data/model';

export default AuthMethodModel.extend({
  otherConfig: belongsTo('mount-config', { async: false, inverse: null }),
});
