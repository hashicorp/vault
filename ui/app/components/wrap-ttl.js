/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { assert } from '@ember/debug';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class WrapTtlComponent extends Component {
  @tracked
  wrapResponse = true;

  constructor() {
    super(...arguments);
    assert('`onChange` handler is a required attr in `' + this.toString() + '`.', this.args.onChange);
  }

  get wrapTTL() {
    const { wrapResponse, ttl } = this;
    return wrapResponse ? ttl : null;
  }

  @action
  changedValue(ttlObj) {
    this.wrapResponse = ttlObj.enabled;
    this.ttl = ttlObj.goSafeTimeString;
    this.args.onChange(this.wrapTTL);
  }
}
