/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  // we have to have two list routes because ember won't let a query param (e.g. *secret_prefix) be blank.
  // and we only have a query param when we're on a nested secret route like beep/boop/
  this.route('list');
  this.route('list-nested-secret', { path: '/list/*secret_prefix' });
  this.route('create');
  this.route('secret', { path: '/:name' }, function () {
    this.route('details');
    this.route('edit');
    this.route('metadata', function () {
      this.route('edit');
      this.route('versions');
      this.route('diff');
    });
  });
  this.route('configuration');
});
