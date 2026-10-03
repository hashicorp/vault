/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
export default class VersionHistoryModel extends Model {
  @attr('string') version;
  @attr('string') previousVersion;
  @attr('string') timestampInstalled;
}
