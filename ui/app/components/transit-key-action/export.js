/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class ExportComponent extends Component {
  @tracked
  wrapTTL = null;
  @tracked
  exportVersion = false;
}
