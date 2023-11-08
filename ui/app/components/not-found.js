/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Component from '@glimmer/component';

export default class NotFound extends Component {
  @service router;

  get path() {
    return this.router.currentURL;
  }
}
