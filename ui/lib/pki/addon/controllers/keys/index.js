/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { getOwner } from '@ember/application';

export default class PkiKeysIndexController extends Controller {
  get mountPoint() {
    return getOwner(this).mountPoint;
  }
}
