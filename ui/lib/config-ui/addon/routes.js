/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  this.route('messages', function () {
    this.route('create');
    this.route('message', { path: '/:id' }, function () {
      this.route('details');
      this.route('edit');
    });
  });

  this.route('login-settings', function () {
    this.route('rule', { path: '/:name' }, function () {
      this.route('details');
    });
  });
});
