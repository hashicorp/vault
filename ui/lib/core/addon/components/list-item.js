/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import { computed } from '@ember/object';
import layout from '../templates/components/list-item';

export default Component.extend({
  layout,
  flashMessages: service(),
  tagName: '',
  linkParams: null,
  queryParams: null,
  componentName: null,
  hasMenu: true,

  callMethod: task(function* (method, model, successMessage, failureMessage, successCallback = () => {}) {
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
  }),
  link: computed('linkParams.[]', function () {
    if (!Array.isArray(this.linkParams) || !this.linkParams.length) return {};
    const [route, ...models] = this.linkParams;
    return { route, models };
  }),
});
