/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ListController from 'core/mixins/list-controller';
import Controller from '@ember/controller';
import { computed } from '@ember/object';
import { getOwner } from '@ember/application';
import { task } from 'ember-concurrency';
import { inject as service } from '@ember/service';

export default Controller.extend(ListController, {
  flashMessages: service(),
  roleToDelete: null,

  mountPoint: computed(function () {
    return getOwner(this).mountPoint;
  }),

  // template originally called callMethod from the list-item.js contextual component
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
    } finally {
      this.set('roleToDelete', null);
    }
  }),
});
