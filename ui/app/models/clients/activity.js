/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
export default class Activity extends Model {
  @attr('array') byMonth;
  @attr('array') byNamespace;
  @attr('object') total;
  @attr('string') startTime;
  @attr('string') endTime;
  @attr('string') responseTimestamp;
}
