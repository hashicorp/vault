/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { next } from '@ember/runloop';
import { getOwner } from '@ember/application';

export default class PkiIssuerIndexController extends Controller {
  queryParams = ['page'];

  get mountPoint() {
    return getOwner(this).mountPoint;
  }
  // To prevent production build bug of passing D.actions to on "click": https://github.com/hashicorp/vault/pull/16983
  @action onLinkClick(D) {
    next(() => D.actions.close());
  }
}
