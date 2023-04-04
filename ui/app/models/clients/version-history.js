/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
export default class VersionHistoryModel extends Model {
  @attr('string') version;
  @attr('string') previousVersion;
  @attr('string') timestampInstalled;
}
