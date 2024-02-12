/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { task } from 'ember-concurrency';

export default class ListItemComponent extends Component {
  @service flashMessages;

  @task *callMethod(method, model, successMessage, failureMessage, successCallback = () => {}) {
    const flash = this.flashMessages;
    try {
      yield model[method]();
      flash.success(successMessage);
      successCallback();
    } catch (e) {
      const errString = e.errors.join(' ');
      flash.danger(failureMessage + ' ' + errString);
      model.rollbackAttributes();
    }
  }

  get link() {
    if (!Array.isArray(this.linkParams) || !this.linkParams.length) return {};
    const [route, ...models] = this.linkParams;
    return { route, models };
  }
}
