/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Helper from '@ember/component/helper';
import { Promise } from 'rsvp';

export default class AwaitHelper extends Helper {
  compute([promise]) {
    if (!promise || typeof promise.then !== 'function') {
      return promise;
    }
    if (promise !== this.lastPromise) {
      this.lastPromise = promise;
      this.value = null;
      this.resolve(promise);
    }
    return this.value;
  }
  async resolve(promise) {
    let value;
    try {
      value = await Promise.resolve(promise);
    } catch (error) {
      value = error;
    } finally {
      // ensure this promise is still the newest promise
      // otherwise avoid firing recompute since a newer promise is in flight
      if (promise === this.lastPromise) {
        this.value = value;
        this.recompute();
      }
    }
  }
}
