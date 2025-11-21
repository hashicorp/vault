/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
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
