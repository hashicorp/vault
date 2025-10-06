/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// this model is just used for integration tests
//

import Model, { belongsTo, attr } from '@ember-data/model';

export default class TestFormModel extends Model {
  @belongsTo('mount-config', { async: false, inverse: null }) config;
  @belongsTo('mount-config', { async: false, inverse: null }) otherConfig;

  @attr('string') path;
  @attr('string', { editType: 'textarea' }) description;
}
