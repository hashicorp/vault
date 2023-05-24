/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { dropTask, task } from 'ember-concurrency';
import { inject as service } from '@ember/service';

export default class VaultClusterAccessMethodsController extends Controller {
  @service flashMessages;

  queryParams = ['page, pageFilter'];

  page = 1;
  pageFilter = null;
  filter = null;

  @task
  @dropTask
  *disableMethod(method) {
    const { type, path } = method;
    try {
      yield method.destroyRecord();
      this.flashMessages.success(`The ${type} Auth Method at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling Auth Method at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }
}
