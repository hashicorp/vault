/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  // we have to have two list routes because *nestedSecret can't be blank and we only have *secret param when we're in a nested view like beep/bop
  this.route('list-root', { path: '/list/' });
  this.route('list', { path: '/list/*nestedSecret' });
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
